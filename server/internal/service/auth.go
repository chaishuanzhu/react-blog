package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"blog-server/internal/apperr"
	"blog-server/internal/config"
	"blog-server/internal/model"
	"blog-server/internal/store"
)

const (
	jwtIssuer         = "blog-server"
	minPasswordLength = 8
)

var errInvalidCredentials = &apperr.Error{Status: 401, Code: "INVALID_CREDENTIALS", Message: "email or password is incorrect"}

type Auth struct {
	q      store.Querier
	cfg    *config.Config
	secret []byte
	// dummyHash keeps login timing the same whether or not the email exists.
	dummyHash []byte
}

type claims struct {
	Role string `json:"role"`
	// Version must match users.token_version; bumping it revokes every token issued before.
	Version uint32 `json:"ver"`
	jwt.RegisteredClaims
}

func NewAuth(q store.Querier, cfg *config.Config) *Auth {
	secret := []byte(cfg.JWTSecret)
	if len(secret) == 0 {
		buf := make([]byte, 32)
		_, _ = rand.Read(buf)
		secret = []byte(hex.EncodeToString(buf))
		slog.Warn("JWT_SECRET not set, using a random secret; tokens will not survive restarts")
	}
	dummy, _ := bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)
	return &Auth{q: q, cfg: cfg, secret: secret, dummyHash: dummy}
}

// EnsureAdmin creates the admin account from ADMIN_EMAIL / ADMIN_PASSWORD if it does not exist yet.
// An existing account is left untouched so a password changed through the API is not reset on restart.
func (s *Auth) EnsureAdmin(ctx context.Context) error {
	email := strings.ToLower(s.cfg.AdminEmail)
	if email == "" || s.cfg.AdminPassword == "" {
		return nil
	}
	_, err := s.q.GetUserByEmail(ctx, email)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("look up admin: %w", err)
	}
	if len(s.cfg.AdminPassword) < minPasswordLength {
		return fmt.Errorf("ADMIN_PASSWORD must be at least %d characters", minPasswordLength)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(s.cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	nickname := s.cfg.AdminNickname
	if nickname == "" {
		nickname = "admin"
	}
	if _, err := s.q.CreateUser(ctx, store.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		Nickname:     nickname,
		Avatar:       avatarFor(email),
		Website:      s.cfg.SiteURL,
	}); err != nil {
		return fmt.Errorf("create admin: %w", err)
	}
	slog.Info("admin account created", "email", email)
	return nil
}

func (s *Auth) Login(ctx context.Context, email, password string) (*model.LoginResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.q.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		_ = bcrypt.CompareHashAndPassword(s.dummyHash, []byte(password))
		return nil, errInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("look up user: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errInvalidCredentials
	}
	return s.issueToken(user)
}

func (s *Auth) issueToken(user store.User) (*model.LoginResult, error) {
	expires := time.Now().Add(s.cfg.JWTTTL)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		Role:    string(user.Role),
		Version: user.TokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   strconv.FormatUint(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}).SignedString(s.secret)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}
	return &model.LoginResult{Token: token, ExpiresAt: expires.UTC(), User: toUserModel(user)}, nil
}

// VerifyToken returns the user id of a valid admin token that has not been revoked.
func (s *Auth) VerifyToken(ctx context.Context, raw string) (uint64, error) {
	var c claims
	_, err := jwt.ParseWithClaims(raw, &c, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(jwtIssuer), jwt.WithExpirationRequired())
	if err != nil || c.Role != string(store.UsersRoleAdmin) {
		return 0, apperr.ErrUnauthorized
	}
	id, err := strconv.ParseUint(c.Subject, 10, 64)
	if err != nil {
		return 0, apperr.ErrUnauthorized
	}
	u, err := s.q.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, apperr.ErrUnauthorized
	}
	if err != nil {
		return 0, fmt.Errorf("look up token user: %w", err)
	}
	if u.TokenVersion != c.Version {
		return 0, apperr.ErrUnauthorized
	}
	return id, nil
}

func (s *Auth) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	u, err := s.q.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ErrUnauthorized
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	m := toUserModel(u)
	return &m, nil
}

func (s *Auth) UpdateProfile(ctx context.Context, id uint64, in model.ProfileInput) (*model.User, error) {
	nickname := strings.TrimSpace(in.Nickname)
	if n := utf8.RuneCountInString(nickname); n < 1 || n > 32 {
		return nil, apperr.BadRequest("nickname must be 1-32 characters")
	}
	avatar, err := optionalURL("avatar", in.Avatar)
	if err != nil {
		return nil, err
	}
	website, err := optionalURL("website", in.Website)
	if err != nil {
		return nil, err
	}
	if _, err := s.q.UpdateUserProfile(ctx, store.UpdateUserProfileParams{
		Nickname: nickname, Avatar: avatar, Website: website, ID: id,
	}); err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	return s.GetUser(ctx, id)
}

// ChangePassword revokes every existing token of the user and returns a fresh one for the caller.
func (s *Auth) ChangePassword(ctx context.Context, id uint64, current, next string) (*model.LoginResult, error) {
	u, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, apperr.ErrUnauthorized
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(current)) != nil {
		return nil, errInvalidCredentials
	}
	if len(next) < minPasswordLength || len(next) > 72 {
		return nil, apperr.BadRequest(fmt.Sprintf("new password must be %d-72 bytes", minPasswordLength))
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(next), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if _, err := s.q.UpdateUserPassword(ctx, store.UpdateUserPasswordParams{PasswordHash: string(hash), ID: id}); err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}
	u, err = s.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("reload user: %w", err)
	}
	return s.issueToken(u)
}

func toUserModel(u store.User) model.User {
	return model.User{
		ID: u.ID, Email: u.Email, Nickname: u.Nickname, Avatar: u.Avatar,
		Website: u.Website, Role: string(u.Role), CreatedAt: u.CreatedAt,
	}
}
