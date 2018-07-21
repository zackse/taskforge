package task

import (
	"encoding/json"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type Note struct {
	CreatedDate time.Time
	Body        string
}

func (n Note) String() string {
	jsn, err := json.MarshalIndent(n, "", "    ")
	if err != nil {
		return n.Body
	}

	return string(jsn)
}

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

func New(title string) Task {
	t := Task{
		Title:       title,
		CreatedDate: time.Now(),
	}

	t.ID = objectid.New().Hex()

	return t
}
