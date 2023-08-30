package proxy

import "net/http"

type Server struct {
	server *http.Server
	router *http.ServeMux
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}
