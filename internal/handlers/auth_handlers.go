package handlers

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal/auth"
	"github.com/anishsharma21/go-web-dev-template/internal/queries"
	"github.com/anishsharma21/go-web-dev-template/internal/types/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := template.HTMLEscapeString(r.FormValue("email"))
		firstName := template.HTMLEscapeString(r.FormValue("first_name"))
		lastName := template.HTMLEscapeString(r.FormValue("last_name"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		if email == "" || firstName == "" || lastName == "" || password == "" {
			slog.Error("Email, first name, last name or password is empty")
			http.Error(w, "Email, first name, last name or password is empty", http.StatusBadRequest)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Failed to hash password", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		passwordHashString := string(passwordHash)

		err = queries.SignUpNewUser(r.Context(), dbPool, models.User{
			Email:     email,
			FirstName: &firstName,
			LastName:  &lastName,
			Password:  &passwordHashString,
		})
		if err != nil {
			slog.Error("Failed to sign up new user", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		accessToken, err := auth.CreateAccessToken(email)
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
	})
}
