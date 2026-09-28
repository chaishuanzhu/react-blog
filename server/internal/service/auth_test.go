package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"blog-server/internal/apperr"
	"blog-server/internal/config"
	"blog-server/internal/store"
)

type fakeUsers struct {
	store.Querier
	user store.User
}

func (f *fakeUsers) GetUserByEmail(_ context.Context, email string) (store.User, error) {
	if email != f.user.Email {
		return store.User{}, sql.ErrNoRows
	}
	return f.user, nil
}

func (f *fakeUsers) GetUserByID(_ context.Context, id uint64) (store.User, error) {
	if id != f.user.ID {
		return store.User{}, sql.ErrNoRows
	}
	return f.user, nil
}

func (f *fakeUsers) UpdateUserPassword(_ context.Context, arg store.UpdateUserPasswordParams) (int64, error) {
	f.user.PasswordHash = arg.PasswordHash
	f.user.TokenVersion++
	return 1, nil
}

func TestChangePasswordRevokesTokens(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.MinCost)
	q := &fakeUsers{user: store.User{ID: 7, Email: "a@example.com", PasswordHash: string(hash), Role: store.UsersRoleAdmin}}
	svc := NewAuth(q, &config.Config{JWTSecret: "0123456789abcdef0123456789abcdef", JWTTTL: time.Hour})
	ctx := context.Background()

	first, err := svc.Login(ctx, "a@example.com", "old-password")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if id, err := svc.VerifyToken(ctx, first.Token); err != nil || id != 7 {
		t.Fatalf("fresh token rejected: id=%d err=%v", id, err)
	}

	if _, err := svc.ChangePassword(ctx, 7, "wrong-password", "new-password"); !errors.Is(err, errInvalidCredentials) {
		t.Fatalf("wrong current password: got %v", err)
	}
	next, err := svc.ChangePassword(ctx, 7, "old-password", "new-password")
	if err != nil {
		t.Fatalf("change password: %v", err)
	}

	if _, err := svc.VerifyToken(ctx, first.Token); !errors.Is(err, apperr.ErrUnauthorized) {
		t.Fatalf("old token should be revoked, got %v", err)
	}
	if id, err := svc.VerifyToken(ctx, next.Token); err != nil || id != 7 {
		t.Fatalf("new token rejected: id=%d err=%v", id, err)
	}
	if _, err := svc.Login(ctx, "a@example.com", "old-password"); !errors.Is(err, errInvalidCredentials) {
		t.Fatalf("old password still accepted: %v", err)
	}
}

func TestVerifyTokenRejectsForeignSecret(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.MinCost)
	q := &fakeUsers{user: store.User{ID: 1, Email: "a@example.com", PasswordHash: string(hash), Role: store.UsersRoleAdmin}}
	issuer := NewAuth(q, &config.Config{JWTSecret: "0123456789abcdef0123456789abcdef", JWTTTL: time.Hour})
	verifier := NewAuth(q, &config.Config{JWTSecret: "fedcba9876543210fedcba9876543210", JWTTTL: time.Hour})

	res, err := issuer.Login(context.Background(), "a@example.com", "password1")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if _, err := verifier.VerifyToken(context.Background(), res.Token); !errors.Is(err, apperr.ErrUnauthorized) {
		t.Fatalf("token signed with another secret accepted: %v", err)
	}
}
