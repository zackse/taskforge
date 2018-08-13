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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chasinglogic/taskforge/list"
	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/ql/lexer"
	"github.com/chasinglogic/taskforge/ql/parser"
	"github.com/chasinglogic/taskforge/task"
)

type ListAPI struct {
	list list.List
}

func (l ListAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		l.get(w, r)
	case "POST":
		l.post(w, r)
	default:
		unsupportedMethod(w, r)
	}
}

func (l ListAPI) get(w http.ResponseWriter, r *http.Request) {
	var response interface{}
	var err error

	switch r.URL.Path {
	case "/list/current":
		response, err = l.list.Current()
	default:
		q := r.FormValue("q")
		if q == "" {
			q = r.FormValue("query")
		}

		if q == "" {
			response, err = l.list.Slice()
			break
		}

		p := parser.New(lexer.New(q))
		tree := p.Parse()

		if p.Error() != nil {
			err = p.Error()
			fmt.Println(err)
			break
		}

		response, err = l.list.Search(tree)
	}

	if err != nil {
		if err == list.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		sendJSON(w, apiError{Message: err.Error()})
		return
	}

	sendJSON(w, response)
}

func (l ListAPI) post(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	if strings.HasSuffix(r.URL.Path, "/query") {
		var tree ast.AST
		err := decoder.Decode(&tree)
		if err != nil {
			badRequest(w, err)
			return
		}

		tasks, err := l.list.Search(tree)
		if err != nil {
			fiveHundred(w, err)
			return
		}

		sendJSON(w, tasks)
		return
	}

	var tasks []task.Task
	err := decoder.Decode(&tasks)
	if err != nil {
		badRequest(w, err)
		return
	}

	err = l.list.AddMultiple(tasks)
	if err != nil {
		fiveHundred(w, err)
	}
}
