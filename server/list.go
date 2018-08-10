package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chasinglogic/taskforge/list"
	"github.com/chasinglogic/taskforge/ql/lexer"
	"github.com/chasinglogic/taskforge/ql/parser"
	"github.com/chasinglogic/taskforge/task"
)

func unsupportedMethod(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(apiError{Message: "unsupported method"}.Marshal())
}

type apiError struct {
	Message string
}

func (ae apiError) Marshal() []byte {
	jsn, err := json.Marshal(ae)
	if err != nil {
		return []byte(err.Error())
	}

	return jsn
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

		w.Write(apiError{Message: err.Error()}.Marshal())
		return
	}

	sendJSON(w, response)
}

func (l ListAPI) post(w http.ResponseWriter, r *http.Request) {
	var tasks []task.Task
	err := json.NewDecoder(r.Body).Decode(&tasks)
	if err != nil {
		badRequest(w, err)
	}

	err = l.list.AddMultiple(tasks)
	if err != nil {
		fiveHundred(w, err)
	}
}
