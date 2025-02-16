package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal/auth"
	"github.com/anishsharma21/go-web-dev-template/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func Login(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := template.HTMLEscapeString(r.FormValue("email"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		if email == "" || password == "" {
			slog.Error("Email or password is empty")
			http.Error(w, "Email or password is empty", http.StatusBadRequest)
			return
		}

		user, err := queries.GetUserByEmail(r.Context(), dbPool, email)
		if err != nil {
			slog.Error("Failed to find user when logging in", "error", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Email != email {
			slog.Error("User email does not match", "user_email", user.Email, "email", email)
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
		if err != nil {
			slog.Error("Failed to compare password hashes", "error", err)
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		accessToken, err := auth.CreateAccessToken(user.Email)
		if err != nil {
			slog.Error("Failed to create JWT token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		refreshToken, err := auth.CreateRefreshToken(email)
		if err != nil {
			slog.Error("Failed to create refresh token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now().Add(24 * 7 * time.Hour),
		})

		response := map[string]string{
			"token": accessToken,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("Failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		slog.InfoContext(r.Context(), fmt.Sprintf("User logged in: %s", email))
	})
}
