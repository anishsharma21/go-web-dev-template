package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/anishsharma21/go-web-dev-template/internal/types/selectors"
)

func RenderBaseView(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, selectors.IndexPage.BaseHTML, nil)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to execute template", "error", err, "template", selectors.IndexPage.BaseHTML)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
