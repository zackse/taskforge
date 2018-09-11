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

package task

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Note is a note added to a task, analogous to a comment but only for ones
// personal use.
type Note struct {
	ID          string    `bson:"_id" json:"id"`
	CreatedDate time.Time `json:"created_date"`
	Body        string    `json:"body"`
}

func (n Note) String() string {
	jsn, err := json.MarshalIndent(n, "", "    ")
	if err != nil {
		return n.Body
	}

	return string(jsn)
}

// NewNote will properly fill out the metadata of a note with the given body
func NewNote(body string) Note {
	return Note{
		ID:          objectid.New().Hex(),
		CreatedDate: time.Now(),
		Body:        body,
	}
}

// Task is a unit of work
type Task struct {
	ID            string    `bson:"_id" json:"id"`
	Title         string    `json:"title"`
	CreatedDate   time.Time `json:"created_date"`
	Context       string    `json:"context"`
	Priority      float64   `json:"priority,omitempty"`
	Notes         []Note    `json:"notes,omitempty"`
	CompletedDate time.Time `json:"completed_date,omitempty"`
	Body          string    `json:"body,omitempty"`
}

func (t Task) String() string {
	jsn, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return t.Title
	}

	return string(jsn)
}

// New creates a task with the given title. This properly populates meta data
// fields with non-zero values.
func New(title string) Task {
	return Task{
		ID:          objectid.New().Hex(),
		Title:       title,
		CreatedDate: time.Now(),
		Context:     "default",
		Notes:       []Note{},
	}
}

// Complete completes this task.
func (t *Task) Complete() {
	t.CompletedDate = time.Now()
}

// IsCompleted returns a boolean indicating whether this task is complete or
// not.
func (t *Task) IsCompleted() bool {
	return !t.CompletedDate.IsZero()
}

// IsComplete is an alias for task.IsCompleted
func (t *Task) IsComplete() bool {
	return t.IsCompleted()
}

// Sort a slice of tasks using the expected sort order of
// Priority then CreatedDate.
func Sort(l []Task) []Task {
	sort.Slice(l, func(i, j int) bool {
		return (l[i].Priority > l[j].Priority) ||
			(l[i].Priority == l[j].Priority &&
				l[i].CreatedDate.Before(l[j].CreatedDate))
	})

	return l
}
