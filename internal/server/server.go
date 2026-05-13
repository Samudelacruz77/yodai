package server

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sdelacru/yodai/internal/chat"
	"github.com/sdelacru/yodai/internal/config"
	"github.com/sdelacru/yodai/web"
)

func NewRouter(cfg *config.Config, chatClient *chat.Client) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	handler := chat.NewHandler(chatClient, cfg)

	r.Route("/api", func(r chi.Router) {
		r.Post("/chat", handler.HandleChat)
		r.Post("/chat/stream", handler.HandleStreamChat)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
	})

	staticFS, _ := fs.Sub(web.StaticFS, ".")
	r.Handle("/*", http.FileServer(http.FS(staticFS)))

	return r
}
