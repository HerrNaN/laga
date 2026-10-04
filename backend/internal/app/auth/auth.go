package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/herrnan/laga/internal/app/auth/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	registrationCookie = "laga_registration"
	loginCookie        = "laga_login"
	sessionCookie      = "laga_session"
	registrationTTL    = 5 * time.Minute
	loginTTL           = 5 * time.Minute
	sessionTTL         = 30 * 24 * time.Hour
)

type Handler struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	wa      *webauthn.WebAuthn
	secure  bool
	log     *slog.Logger
}

type user struct {
	id          []byte
	name        string
	credentials []webauthn.Credential
}

func (u user) WebAuthnID() []byte                         { return u.id }
func (u user) WebAuthnName() string                       { return u.name }
func (u user) WebAuthnDisplayName() string                { return u.name }
func (u user) WebAuthnCredentials() []webauthn.Credential { return u.credentials }

func New(pool *pgxpool.Pool, rpID, origin string, log *slog.Logger) (*Handler, error) {
	if rpID == "" || origin == "" {
		return nil, errors.New("rp id and origin are required")
	}

	wa, err := webauthn.New(&webauthn.Config{
		RPID:          rpID,
		RPOrigins:     []string{origin},
		RPDisplayName: "Laga",
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementRequired,
			UserVerification: protocol.VerificationRequired,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("configure passkeys: %w", err)
	}

	return &Handler{pool: pool, queries: db.New(pool), wa: wa, secure: strings.HasPrefix(origin, "https://"), log: log}, nil
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/register/begin", h.beginRegistration)
	mux.HandleFunc("POST /api/auth/register/finish", h.finishRegistration)
	mux.HandleFunc("POST /api/auth/login/begin", h.beginLogin)
	mux.HandleFunc("POST /api/auth/login/finish", h.finishLogin)
	mux.HandleFunc("GET /api/auth/me", h.me)
}

func randomToken() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(value), nil
}

func tokenHash(token string) ([]byte, error) {
	value, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(value) != 32 {
		return nil, errors.New("invalid token")
	}

	hash := sha256.Sum256(value)

	return hash[:], nil
}

func timestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func (h *Handler) cookie(w http.ResponseWriter, name, value, path string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(ttl.Seconds()),
	})
}

func (h *Handler) beginRegistration(w http.ResponseWriter, r *http.Request) {
	id := make([]byte, 32)
	if _, err := rand.Read(id); err != nil {
		h.fail(w, "generate user ID", err)
		return
	}

	u := user{id: id, name: "Laga " + base64.RawURLEncoding.EncodeToString(id[:6])}

	options, ceremony, err := h.wa.BeginRegistration(u)
	if err != nil {
		h.fail(w, "begin registration", err)
		return
	}

	registrationToken, err := randomToken()
	if err != nil {
		h.fail(w, "generate registration token", err)
		return
	}

	hash, _ := tokenHash(registrationToken)

	ceremonyJSON, err := json.Marshal(ceremony)
	if err != nil {
		h.fail(w, "encode registration", err)
		return
	}

	err = h.queries.InsertRegistrationAttempt(r.Context(), db.InsertRegistrationAttemptParams{
		TokenHash: hash,
		UserID:    u.id,
		Name:      u.name,
		Ceremony:  ceremonyJSON,
		ExpiresAt: timestamp(time.Now().Add(registrationTTL)),
	})
	if err != nil {
		h.fail(w, "save registration", err)
		return
	}

	h.cookie(w, registrationCookie, registrationToken, "/api/auth/register", registrationTTL)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(options.Response)
	if err != nil {
		h.log.Error("write registration options", "error", err)
	}
}

func (h *Handler) finishRegistration(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(registrationCookie)
	if err != nil {
		http.Error(w, "Registration expired; try again", http.StatusBadRequest)
		return
	}

	hash, err := tokenHash(cookie.Value)
	if err != nil {
		http.Error(w, "Registration expired; try again", http.StatusBadRequest)
		return
	}

	// Consumed in one statement so a ceremony can only be finished once.
	attempt, err := h.queries.DeleteRegistrationAttempt(r.Context(), hash)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Registration expired; try again", http.StatusBadRequest)
		return
	} else if err != nil {
		h.fail(w, "load registration", err)
		return
	}

	u := user{id: attempt.UserID, name: attempt.Name}

	var ceremony webauthn.SessionData

	err = json.Unmarshal(attempt.Ceremony, &ceremony)
	if err != nil {
		h.fail(w, "decode registration", err)
		return
	}

	credential, err := h.wa.FinishRegistration(u, ceremony, r)
	if err != nil {
		http.Error(w, "Passkey verification failed; try again", http.StatusBadRequest)
		return
	}

	credentialJSON, err := json.Marshal(credential)
	if err != nil {
		h.fail(w, "encode credential", err)
		return
	}

	sessionToken, err := randomToken()
	if err != nil {
		h.fail(w, "generate session token", err)
		return
	}

	sessionHash, _ := tokenHash(sessionToken)

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		h.fail(w, "begin account transaction", err)
		return
	}
	defer tx.Rollback(r.Context())

	queries := h.queries.WithTx(tx)
	err = queries.InsertUser(r.Context(), db.InsertUserParams{ID: u.id, Name: u.name})
	if err != nil {
		h.fail(w, "create account", err)
		return
	}

	err = queries.InsertCredential(r.Context(), db.InsertCredentialParams{
		ID:     credential.ID,
		UserID: u.id,
		Data:   credentialJSON,
	})
	if err != nil {
		h.fail(w, "create account", err)
		return
	}

	err = queries.InsertSession(r.Context(), db.InsertSessionParams{
		TokenHash: sessionHash,
		UserID:    u.id,
		ExpiresAt: timestamp(time.Now().Add(sessionTTL)),
	})
	if err != nil {
		h.fail(w, "create account", err)
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		h.fail(w, "create account", err)
		return
	}

	h.cookie(w, registrationCookie, "", "/api/auth/register", -time.Second)
	h.cookie(w, sessionCookie, sessionToken, "/", sessionTTL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"id":   base64.RawURLEncoding.EncodeToString(u.id),
		"name": u.name,
	})
}

func (h *Handler) beginLogin(w http.ResponseWriter, r *http.Request) {
	options, ceremony, err := h.wa.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationRequired))
	if err != nil {
		h.fail(w, "begin login", err)
		return
	}

	loginToken, err := randomToken()
	if err != nil {
		h.fail(w, "generate login token", err)
		return
	}

	hash, _ := tokenHash(loginToken)
	ceremonyJSON, err := json.Marshal(ceremony)
	if err != nil {
		h.fail(w, "encode login", err)
		return
	}

	err = h.queries.InsertLoginAttempt(r.Context(), db.InsertLoginAttemptParams{
		TokenHash: hash,
		Ceremony:  ceremonyJSON,
		ExpiresAt: timestamp(time.Now().Add(loginTTL)),
	})
	if err != nil {
		h.fail(w, "save login", err)
		return
	}

	h.cookie(w, loginCookie, loginToken, "/api/auth/login", loginTTL)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(options.Response); err != nil {
		h.log.Error("write login options", "error", err)
	}
}

func (h *Handler) finishLogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(loginCookie)
	if err != nil {
		http.Error(w, "Sign-in expired; try again", http.StatusBadRequest)
		return
	}

	hash, err := tokenHash(cookie.Value)
	if err != nil {
		http.Error(w, "Sign-in expired; try again", http.StatusBadRequest)
		return
	}

	// Consume the challenge atomically to prevent replaying a login ceremony.
	ceremonyJSON, err := h.queries.DeleteLoginAttempt(r.Context(), hash)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Sign-in expired; try again", http.StatusBadRequest)
		return
	} else if err != nil {
		h.fail(w, "load login", err)
		return
	}
	h.cookie(w, loginCookie, "", "/api/auth/login", -time.Second)

	var ceremony webauthn.SessionData
	if err = json.Unmarshal(ceremonyJSON, &ceremony); err != nil {
		h.fail(w, "decode login", err)
		return
	}

	lookup := func(rawID, userHandle []byte) (webauthn.User, error) {
		account, err := h.queries.GetUserByCredential(r.Context(), db.GetUserByCredentialParams{
			CredentialID: rawID,
			UserID:       userHandle,
		})
		if err != nil {
			return nil, err
		}

		var credential webauthn.Credential
		if err = json.Unmarshal(account.Data, &credential); err != nil {
			return nil, err
		}

		return user{id: account.ID, name: account.Name, credentials: []webauthn.Credential{credential}}, nil
	}

	u, credential, err := h.wa.FinishPasskeyLogin(lookup, ceremony, r)
	if err != nil {
		http.Error(w, "Passkey verification failed; try again", http.StatusBadRequest)
		return
	}

	credentialJSON, err := json.Marshal(credential)
	if err != nil {
		h.fail(w, "encode credential", err)
		return
	}

	sessionToken, err := randomToken()
	if err != nil {
		h.fail(w, "generate session token", err)
		return
	}
	sessionHash, _ := tokenHash(sessionToken)

	tx, err := h.pool.Begin(r.Context())
	if err != nil {
		h.fail(w, "begin login transaction", err)
		return
	}
	defer tx.Rollback(r.Context())

	queries := h.queries.WithTx(tx)
	err = queries.UpdateCredential(r.Context(), db.UpdateCredentialParams{
		ID: credential.ID, UserID: u.WebAuthnID(), Data: credentialJSON,
	})
	if err != nil {
		h.fail(w, "update credential", err)
		return
	}

	err = queries.InsertSession(r.Context(), db.InsertSessionParams{
		TokenHash: sessionHash,
		UserID:    u.WebAuthnID(),
		ExpiresAt: timestamp(time.Now().Add(sessionTTL)),
	})
	if err != nil {
		h.fail(w, "create session", err)
		return
	}

	if err = tx.Commit(r.Context()); err != nil {
		h.fail(w, "complete login", err)
		return
	}

	h.cookie(w, sessionCookie, sessionToken, "/", sessionTTL)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"id":   base64.RawURLEncoding.EncodeToString(u.WebAuthnID()),
		"name": u.WebAuthnName(),
	})
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		http.Error(w, "Not signed in", http.StatusUnauthorized)
		return
	}

	hash, err := tokenHash(cookie.Value)
	if err != nil {
		http.Error(w, "Not signed in", http.StatusUnauthorized)
		return
	}

	account, err := h.queries.GetUserBySession(r.Context(), hash)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Not signed in", http.StatusUnauthorized)
		return
	} else if err != nil {
		h.fail(w, "load session", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"id":   base64.RawURLEncoding.EncodeToString(account.ID),
		"name": account.Name,
	})
}

func (h *Handler) fail(w http.ResponseWriter, operation string, err error) {
	h.log.Error(operation, "error", err)
	http.Error(w, "Server error", http.StatusInternalServerError)
}
