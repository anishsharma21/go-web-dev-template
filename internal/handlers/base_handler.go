package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/anishsharma21/go-web-dev-template/internal/types"
	"github.com/anishsharma21/go-web-dev-template/internal/types/selectors"
)

func RenderBaseView(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := types.BaseView{
			RenderBaseView:   true,
			RenderLoginView:  false,
			RenderSignupView: false,
		}

		err := tmpl.ExecuteTemplate(w, selectors.IndexPage.IndexHtml, data)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to execute base index template", "error", err, "template", selectors.IndexPage.IndexHtml)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
