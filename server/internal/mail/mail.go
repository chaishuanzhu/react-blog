package mail

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"

	"blog-server/internal/config"
)

type Message struct {
	To      string
	Subject string
	HTML    string
}

type Sender interface {
	Send(ctx context.Context, msg Message) error
}

type SMTPSender struct {
	cfg config.SMTP
}

func NewSMTPSender(cfg config.SMTP) *SMTPSender {
	return &SMTPSender{cfg: cfg}
}

func (s *SMTPSender) Send(ctx context.Context, msg Message) error {
	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tlsCfg := &tls.Config{ServerName: s.cfg.Host, MinVersion: tls.VersionTLS12}

	var conn net.Conn
	var err error
	if s.cfg.Port == 465 {
		conn, err = (&tls.Dialer{NetDialer: dialer, Config: tlsCfg}).DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("dial smtp: %w", err)
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("smtp handshake: %w", err)
	}
	defer client.Close()

	if s.cfg.Port != 465 {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return fmt.Errorf("smtp server does not support STARTTLS")
		}
		if err := client.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}
	if err := client.Auth(smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(s.cfg.From); err != nil {
		return fmt.Errorf("smtp MAIL FROM: %w", err)
	}
	if err := client.Rcpt(msg.To); err != nil {
		return fmt.Errorf("smtp RCPT TO: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA: %w", err)
	}
	if _, err := w.Write(buildMessage(s.cfg, msg, time.Now())); err != nil {
		w.Close()
		return fmt.Errorf("write message: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("finish message: %w", err)
	}
	return client.Quit()
}

func buildMessage(cfg config.SMTP, msg Message, now time.Time) []byte {
	from := (&mail.Address{Name: cfg.FromName, Address: cfg.From}).String()
	domain := cfg.From[strings.LastIndex(cfg.From, "@")+1:]
	id := make([]byte, 12)
	_, _ = rand.Read(id)

	var b bytes.Buffer
	header := func(k, v string) { fmt.Fprintf(&b, "%s: %s\r\n", k, v) }
	header("From", from)
	header("To", msg.To)
	header("Subject", mime.BEncoding.Encode("UTF-8", msg.Subject))
	header("Date", now.Format(time.RFC1123Z))
	header("Message-ID", fmt.Sprintf("<%s@%s>", hex.EncodeToString(id), domain))
	header("MIME-Version", "1.0")
	header("Content-Type", `text/html; charset="UTF-8"`)
	header("Content-Transfer-Encoding", "base64")
	b.WriteString("\r\n")

	encoded := base64.StdEncoding.EncodeToString([]byte(msg.HTML))
	for len(encoded) > 76 {
		b.WriteString(encoded[:76] + "\r\n")
		encoded = encoded[76:]
	}
	b.WriteString(encoded + "\r\n")
	return b.Bytes()
}

// Queue delivers messages on a background goroutine so request handlers never wait on SMTP.
type Queue struct {
	sender Sender
	ch     chan Message
	wg     sync.WaitGroup
}

func NewQueue(sender Sender, size int) *Queue {
	q := &Queue{sender: sender, ch: make(chan Message, size)}
	q.wg.Add(1)
	go q.worker()
	return q
}

func (q *Queue) Enqueue(msg Message) {
	select {
	case q.ch <- msg:
	default:
		slog.Warn("mail queue full, message dropped", "to", msg.To, "subject", msg.Subject)
	}
}

// Close stops accepting messages and waits for queued ones to be delivered.
func (q *Queue) Close() {
	close(q.ch)
	q.wg.Wait()
}

func (q *Queue) worker() {
	defer q.wg.Done()
	for msg := range q.ch {
		var err error
		for attempt := 1; attempt <= 3; attempt++ {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err = q.sender.Send(ctx, msg)
			cancel()
			if err == nil {
				slog.Info("mail sent", "to", msg.To, "subject", msg.Subject)
				break
			}
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}
		if err != nil {
			slog.Error("mail delivery failed", "to", msg.To, "subject", msg.Subject, "err", err)
		}
	}
}

type Notifier interface {
	Enqueue(msg Message)
}

// Noop is used when SMTP is not configured.
type Noop struct{}

func (Noop) Enqueue(msg Message) {
	slog.Debug("mail disabled, message skipped", "to", msg.To, "subject", msg.Subject)
}
