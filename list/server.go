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


package list

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/task"
)

type ServerList struct {
	ServerURL string
	Token     string

	client *http.Client
}

func (sl *ServerList) req(method string, endpoint string, body interface{}) (*http.Request, error) {
	buf := bytes.NewBuffer([]byte{})
	if body != nil {
		jsn, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(jsn)
	}

	url := fmt.Sprintf("%s%s", sl.ServerURL, endpoint)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", sl.Token))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (sl *ServerList) Init() error {
	sl.client = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := sl.req("GET", "/status", nil)
	if err != nil {
		return err
	}

	res, err := sl.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("can not connect to task server")
	}

	return nil
}

func (sl *ServerList) do(req *http.Request, out interface{}) error {
	if sl.client == nil {
		return errors.New("server list has not been initialized")
	}

	if req == nil {
		return errors.New("received nil request")
	}

	res, err := sl.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		var apiErr struct {
			Message string
		}

		err := json.NewDecoder(res.Body).Decode(&apiErr)
		if err != nil {
			return fmt.Errorf("failed with status %d: unable to deserialize body", res.StatusCode)
		}

		return fmt.Errorf("failed with status %d: %s", res.StatusCode, apiErr.Message)
	}

	if out != nil {
		return json.NewDecoder(res.Body).Decode(out)
	}

	return nil
}

// Evaluate the AST and return a List of matching results
func (sl *ServerList) Search(tree ast.AST) ([]task.Task, error) {
	req, err := sl.req("GET", "/list", nil)
	if err != nil {
		return nil, err
	}

	form, _ := url.ParseQuery("")
	form.Add("q", tree.String())
	req.URL.RawQuery = form.Encode()

	var tasks []task.Task
	err = sl.do(req, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Add a task to the List
func (sl *ServerList) Add(t task.Task) error {
	req, err := sl.req("POST", "/task", t)
	if err != nil {
		return err
	}

	return sl.do(req, nil)
}

// Add multiple tasks to the List, should be more efficient resource
// utilization.
func (sl *ServerList) AddMultiple(tasks []task.Task) error {
	req, err := sl.req("POST", "/list", tasks)
	if err != nil {
		return err
	}

	return sl.do(req, nil)
}

// Return a slice of task.Tasks in this List
func (sl *ServerList) Slice() ([]task.Task, error) {
	req, err := sl.req("GET", "/list", nil)
	if err != nil {
		return nil, err
	}

	var tasks []task.Task
	err = sl.do(req, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Find a task by ID
func (sl *ServerList) FindByID(id string) (task.Task, error) {
	req, err := sl.req("GET", fmt.Sprintf("/task/%s", id), nil)
	if err != nil {
		return task.Task{}, err
	}

	var t task.Task
	err = sl.do(req, &t)
	return t, err
}

// Return the current task, meaning the oldest uncompleted task in the List
func (sl *ServerList) Current() (task.Task, error) {
	req, err := sl.req("GET", "/list/current", nil)
	if err != nil {
		return task.Task{}, err
	}

	var t task.Task
	err = sl.do(req, &t)
	return t, err
}

// Complete a task by id
func (sl *ServerList) Complete(id string) error {
	req, err := sl.req("PUT", fmt.Sprintf("/task/%s/complete", id), nil)
	if err != nil {
		return err
	}

	return sl.do(req, nil)
}

// Update a task in the listist, finding the original by the ID of the given task
func (sl *ServerList) Update(t task.Task) error {
	req, err := sl.req("PUT", fmt.Sprintf("/task/%s", t.ID), t)
	if err != nil {
		return err
	}

	return sl.do(req, nil)
}

// Add note to a task by ID
func (sl *ServerList) AddNote(id string, note task.Note) error {
	req, err := sl.req("PUT", fmt.Sprintf("/task/%s/addNote", id), note)
	if err != nil {
		return err
	}

	return sl.do(req, nil)
}
