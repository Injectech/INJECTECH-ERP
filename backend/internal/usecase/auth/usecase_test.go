package auth

import (
	"context"
	"testing"
	"time"

	domainuser "backend/internal/domain/user"
)

type fakeUserRepo struct {
	byID    map[string]domainuser.User
	byEmail map[string]string
}

func newFakeRepo() *fakeUserRepo {
	return &fakeUserRepo{byID: make(map[string]domainuser.User), byEmail: make(map[string]string)}
}

func (f *fakeUserRepo) Create(_ context.Context, u domainuser.User) (domainuser.User, error) {
	f.byID[u.ID] = u
	f.byEmail[u.Email] = u.ID
	return u, nil
}

func (f *fakeUserRepo) GetByID(_ context.Context, id string) (domainuser.User, error) {
	return f.byID[id], nil
}

func (f *fakeUserRepo) GetByEmail(_ context.Context, email string) (domainuser.User, error) {
	id := f.byEmail[email]
	return f.byID[id], nil
}

func (f *fakeUserRepo) Update(_ context.Context, u domainuser.User) error {
	f.byID[u.ID] = u
	return nil
}
func (f *fakeUserRepo) SoftDelete(_ context.Context, id string) error { delete(f.byID, id); return nil }

func TestRegisterAndRefreshRotation(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, time.Minute, time.Hour, "access", "refresh")

	session, err := svc.Register(context.Background(), "a@b.com", "secret", "tester")
	if err != nil {
		t.Fatalf("register error: %v", err)
	}
	tokens := session.Tokens

	// Stale refresh should fail after rotation
	rotated, err := svc.Refresh(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("refresh error: %v", err)
	}

	if rotated.RefreshToken == tokens.RefreshToken {
		t.Fatalf("expected rotated refresh token")
	}

	if _, err := svc.Refresh(tokens.RefreshToken); err == nil {
		t.Fatalf("expected old refresh token to be invalid")
	}
}

func TestLoginUsesHashedPassword(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, time.Minute, time.Hour, "access", "refresh")

	// seed user
	if _, err := svc.Register(context.Background(), "user@test.com", "pw123", "user"); err != nil {
		t.Fatalf("seed register error: %v", err)
	}

	if _, err := svc.Login(context.Background(), "user@test.com", "wrong"); err == nil {
		t.Fatalf("expected invalid credential error")
	}

	if _, err := svc.Login(context.Background(), "user@test.com", "pw123"); err != nil {
		t.Fatalf("login should succeed: %v", err)
	}
}
