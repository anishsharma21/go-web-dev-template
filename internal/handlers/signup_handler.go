package handlers

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal/auth"
	"github.com/anishsharma21/go-web-dev-template/internal/queries"
	"github.com/anishsharma21/go-web-dev-template/internal/types/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var isProduction = os.Getenv("ENV") == "production"

func SignUp(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := template.HTMLEscapeString(r.FormValue("email"))
		firstName := template.HTMLEscapeString(r.FormValue("first_name"))
		lastName := template.HTMLEscapeString(r.FormValue("last_name"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		if email == "" || firstName == "" || lastName == "" || password == "" {
			slog.ErrorContext(r.Context(), "Missing required fields for signup", "email", email, "first_name", firstName, "last_name", lastName)
			http.Error(w, "Email, first name, last name or password is invalid", http.StatusBadRequest)
			return
		}

		if len(password) < 8 {
			slog.ErrorContext(r.Context(), "Password too short", "email", email)
			http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to hash password", "error", err)
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
			slog.ErrorContext(r.Context(), "Failed to sign up new user", "error", err, "email", email, "first_name", firstName, "last_name", lastName)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user, err := queries.GetUserByEmail(r.Context(), dbPool, email)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to find user after signup", "error", err, "email", email, "first_name", firstName, "last_name", lastName)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		accessToken, err := auth.CreateAccessToken(user.ID)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to create JWT access token", "error", err, "email", email, "first_name", firstName, "last_name", lastName)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		refreshToken, err := auth.CreateRefreshToken(user.ID)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to create refresh token", "error", err, "email", email, "first_name", firstName, "last_name", lastName)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   isProduction,
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
			slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
