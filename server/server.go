// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.


package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chasinglogic/taskforge/list"
)

type Server struct {
	Port   int
	Addr   string
	Tokens []string

	list       list.List
	httpServer *http.Server
	taskAPI    http.Handler
	listAPI    http.Handler
}

func New(l list.List, tokens ...string) *Server {
	return &Server{
		Port:    8080,
		Addr:    "localhost",
		Tokens:  tokens,
		list:    l,
		listAPI: ListAPI{list: l},
		taskAPI: TaskAPI{list: l},
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		unauthorized(w)
		return
	}

	// Strip the "Token " from the front of the string
	token = token[len("Token "):]
	if !s.validToken(token) {
		unauthorized(w)
		return
	}

	switch {
	case strings.HasPrefix(r.URL.Path, "/task"):
		s.taskAPI.ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/list"):
		s.listAPI.ServeHTTP(w, r)
	case r.URL.Path == "/status":
		sendJSON(w, apiError{Message: "all systems go"})
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

func (s *Server) Listen() error {
	addr := fmt.Sprintf("%s:%d", s.Addr, s.Port)
	fmt.Println("task server listening on:", addr)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.httpServer.Shutdown(context.Background())
}

func (s *Server) validToken(providedToken string) bool {
	for _, token := range s.Tokens {
		if token == providedToken {
			return true
		}
	}

	return false
}

func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	sendJSON(w, apiError{
		Message: "unauthorized",
	})
}

func unsupportedMethod(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	sendJSON(w, apiError{Message: "unsupported method"})
}

type apiError struct {
	Message string
}

func fiveHundred(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	sendJSON(w, apiError{Message: err.Error()})
}

func badRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	sendJSON(w, apiError{Message: err.Error()})
}

func sendJSON(w http.ResponseWriter, response interface{}) {
	jsn, err := json.Marshal(response)
	if err != nil {
		fiveHundred(w, err)
		return
	}

	_, err = w.Write(jsn)
	if err != nil {
		fmt.Println("ERROR writing to client:", err)
	}
}
