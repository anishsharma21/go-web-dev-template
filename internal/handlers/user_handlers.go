package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/anishsharma21/go-web-dev-template/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RenderBaseUserView(dbPool *pgxpool.Pool, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := queries.GetUsers(r.Context(), dbPool)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to get users", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "base-users-view", users)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to execute template", "error", err, "template", "base-users-view")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
