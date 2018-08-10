package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chasinglogic/taskforge/list"
)

type Server struct {
	Port int
	Addr string

	list    list.List
	taskAPI http.Handler
	listAPI http.Handler
}

func New(l list.List) *Server {
	return &Server{
		Port:    8080,
		Addr:    "localhost",
		list:    l,
		listAPI: ListAPI{list: l},
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/task"):
		s.taskAPI.ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/list"):
		s.listAPI.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

func (s *Server) Listen() {
	addr := fmt.Sprintf("%s:%d", s.Addr, s.Port)
	fmt.Println("task server listening on:", addr)
	http.ListenAndServe(addr, s)
}
