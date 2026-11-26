package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"eneba-marketplace/backend/internal/httpx"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

type Service struct{ DB *sql.DB }
type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
type session struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type userKey struct{}

func Current(r *http.Request) User { u, _ := r.Context().Value(userKey{}).(User); return u }
func hashToken(t string) string    { h := sha256.Sum256([]byte(t)); return hex.EncodeToString(h[:]) }
func validEmail(email string) bool {
	a, e := mail.ParseAddress(email)
	return e == nil && a.Address == email && len(email) <= 254
}
func (s Service) Register(m *http.ServeMux) {
	m.HandleFunc("POST /api/auth/register", s.register)
	m.HandleFunc("POST /api/auth/login", s.login)
	m.Handle("GET /api/auth/me", s.Require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { httpx.JSON(w, 200, Current(r)) })))
	m.Handle("POST /api/auth/logout", s.Require(http.HandlerFunc(s.logout)))
}
func (s Service) register(w http.ResponseWriter, r *http.Request) {
	var c credentials
	if !httpx.Decode(w, r, &c) {
		return
	}
	c.Email = strings.ToLower(strings.TrimSpace(c.Email))
	if !validEmail(c.Email) || len(c.Password) < 8 || len(c.Password) > 72 {
		httpx.Error(w, 400, "Enter a valid email and a password of 8–72 bytes")
		return
	}
	h, e := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if e != nil {
		httpx.Error(w, 500, "Could not create account")
		return
	}
	var u User
	u.Email = c.Email
	e = s.DB.QueryRowContext(r.Context(), `INSERT INTO users(email,password_hash) VALUES($1,$2) RETURNING id`, u.Email, string(h)).Scan(&u.ID)
	if e != nil {
		var pg *pgconn.PgError
		if errors.As(e, &pg) && pg.Code == "23505" {
			httpx.Error(w, 409, "Email already registered")
		} else {
			httpx.Error(w, 500, "Could not create account")
		}
		return
	}
	s.issue(w, r, u, 201)
}
func (s Service) login(w http.ResponseWriter, r *http.Request) {
	var c credentials
	if !httpx.Decode(w, r, &c) {
		return
	}
	var u User
	var h string
	e := s.DB.QueryRowContext(r.Context(), `SELECT id,email,password_hash FROM users WHERE email=$1`, strings.ToLower(strings.TrimSpace(c.Email))).Scan(&u.ID, &u.Email, &h)
	if e != nil && !errors.Is(e, sql.ErrNoRows) {
		httpx.Error(w, 500, "Could not sign in")
		return
	}
	if e != nil || bcrypt.CompareHashAndPassword([]byte(h), []byte(c.Password)) != nil {
		httpx.Error(w, 401, "Invalid email or password")
		return
	}
	s.issue(w, r, u, 200)
}
func (s Service) issue(w http.ResponseWriter, r *http.Request, u User, status int) {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		httpx.Error(w, 500, "Could not start session")
		return
	}
	t := hex.EncodeToString(b)
	_, e := s.DB.ExecContext(r.Context(), `INSERT INTO sessions(token_hash,user_id,expires_at) VALUES($1,$2,$3)`, hashToken(t), u.ID, time.Now().Add(24*time.Hour))
	if e != nil {
		httpx.Error(w, 500, "Could not start session")
		return
	}
	httpx.JSON(w, status, session{u, t})
}
func bearer(r *http.Request) string {
	v := r.Header.Get("Authorization")
	if !strings.HasPrefix(v, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(v, "Bearer ")
}
func (s Service) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := bearer(r)
		if len(t) != 64 {
			httpx.Error(w, 401, "Sign in required")
			return
		}
		var u User
		e := s.DB.QueryRowContext(r.Context(), `SELECT u.id,u.email FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.expires_at>now()`, hashToken(t)).Scan(&u.ID, &u.Email)
		if errors.Is(e, sql.ErrNoRows) {
			httpx.Error(w, 401, "Session expired")
			return
		}
		if e != nil {
			httpx.Error(w, 500, "Could not verify session")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey{}, u)))
	})
}
func (s Service) logout(w http.ResponseWriter, r *http.Request) {
	if _, e := s.DB.ExecContext(r.Context(), `DELETE FROM sessions WHERE token_hash=$1`, hashToken(bearer(r))); e != nil {
		httpx.Error(w, 500, "Could not sign out")
		return
	}
	httpx.JSON(w, 200, map[string]bool{"ok": true})
}
