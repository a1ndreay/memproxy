package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/a1ndreay/memproxy/pkg/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server wraps router and backend.
type Server struct {
	router     *chi.Mux
	backend    cache.Backend
	originAddr string
}

// New creates a new Server instance.
func New(backend cache.Backend, originAddr string) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	srv := &Server{router: r, backend: backend, originAddr: originAddr}
	r.Get("/{key}", srv.getHandler())
	r.Post("/{key}", srv.setHandler())
	r.Delete("/{key}", srv.deleteHandler())
	return srv
}

// ListenAndServe starts HTTP listener.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) getHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		value, err := s.backend.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if value == nil {
			originURL := fmt.Sprintf("%s/%s", s.originAddr, key)
			resp, err := http.Get(originURL)
			if err != nil || resp.StatusCode != http.StatusOK {
				http.Error(w, "cache miss and origin fetch failed "+err.Error(), http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := s.backend.Set(key, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(value)
	}
}

func (s *Server) setHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.backend.Set(key, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) deleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if err := s.backend.Delete(key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
