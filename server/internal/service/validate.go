package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"

	"blog-server/internal/apperr"
)

const (
	mysqlDuplicateEntry = 1062
	mysqlFKViolation    = 1452
)

func requiredText(field, v string, maxRunes int) (string, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", apperr.BadRequest(field + " is required")
	}
	if utf8.RuneCountInString(v) > maxRunes {
		return "", apperr.BadRequest(fmt.Sprintf("%s must be at most %d characters", field, maxRunes))
	}
	return v, nil
}

func optionalText(field, v string, maxRunes int) (string, error) {
	v = strings.TrimSpace(v)
	if utf8.RuneCountInString(v) > maxRunes {
		return "", apperr.BadRequest(fmt.Sprintf("%s must be at most %d characters", field, maxRunes))
	}
	return v, nil
}

func optionalURL(field, v string) (string, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", nil
	}
	return requiredURL(field, v)
}

func requiredURL(field, v string) (string, error) {
	v = strings.TrimSpace(v)
	u, err := url.Parse(v)
	if v == "" || err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || len(v) > 255 {
		return "", apperr.BadRequest(field + " must be an http(s) URL of at most 255 characters")
	}
	return v, nil
}

func isMySQLError(err error, code uint16) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && me.Number == code
}
