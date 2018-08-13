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
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/chasinglogic/taskforge/list"
	"github.com/chasinglogic/taskforge/task"
)

type TaskAPI struct {
	list list.List
}

func (ta TaskAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ta.get(w, r)
	case "POST":
		ta.post(w, r)
	case "PUT":
		ta.put(w, r)
	default:
		unsupportedMethod(w, r)
	}
}

func (ta TaskAPI) get(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/task/"):]
	task, err := ta.list.FindByID(taskID)
	if err != nil {
		if err == list.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		sendJSON(w, apiError{Message: err.Error()})
		return
	}

	sendJSON(w, task)
}

func (ta TaskAPI) put(w http.ResponseWriter, r *http.Request) {
	var err error

	switch {
	case strings.HasSuffix(r.URL.Path, "complete"):
		taskID := r.URL.Path[len("/task/"):]
		taskID = taskID[:len(taskID)-len("/complete")]
		err = ta.list.Complete(taskID)
	case strings.HasSuffix(r.URL.Path, "addnote"):
		fallthrough
	case strings.HasSuffix(r.URL.Path, "addNote"):
		taskID := r.URL.Path[len("/task/"):]
		taskID = taskID[:len(taskID)-len("/addNote")]
		var note task.Note
		err = json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			break
		}

		err = ta.list.AddNote(taskID, note)
	default:
		taskID := r.URL.Path[len("/task/"):]
		var t task.Task
		err = json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			break
		}

		if t.ID == "" {
			t.ID = taskID
		}

		err = ta.list.Update(t)
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

	w.Write([]byte{})
}

func (ta TaskAPI) post(w http.ResponseWriter, r *http.Request) {
	var t task.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		badRequest(w, err)
		return
	}

	if t.ID == "" || t.Title == "" {
		badRequest(w, errors.New("must provide title and task ID"))
		return
	}

	if t.CreatedDate.IsZero() {
		t.CreatedDate = time.Now()
	}

	if t.Notes == nil {
		t.Notes = []task.Note{}
	}

	err = ta.list.Add(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendJSON(w, apiError{Message: err.Error()})
		return
	}
}
