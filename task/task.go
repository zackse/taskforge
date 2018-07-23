package task

import (
	"encoding/json"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Note is a note added to a task, analogous to a comment but only for ones
// personal use.
type Note struct {
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
		CreatedDate: time.Now(),
		Body:        body,
	}
}

// Task is a unit of work
type Task struct {
	ID            string    `json:"id"`
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
	t := Task{
		Title:       title,
		CreatedDate: time.Now(),
		Context:     "default",
	}

	t.ID = objectid.New().Hex()

	return t
}

// Complete completes this task.
func (t *Task) Complete() {
	t.CompletedDate = time.Now()
}
