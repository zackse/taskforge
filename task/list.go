package task

import (
	"errors"
	"sort"
	"time"
)

var ErrNotFound = errors.New("no task with that id exists")

type List interface {
	// Return a new List which has all completed task if yes_or_no is true and all
	// uncompleted tasks if yes_or_no is false.
	Completed(completed bool) []Task
	// Return a new List with only tasks in the given context
	Context(context string) []Task
	// Evaluate the AST and return a List of matching results
	// search(ast: query::ast::AST) []Task
	// Add a task to the List
	Add(Task) error
	// Add multiple tasks to the List, should be more efficient resource
	// utilization.
	AddMultiple(task []Task) error
	// Return a slice of Tasks in this List
	Slice() []Task
	// Find a task by ID
	FindById(id string) *Task
	// Return the current task, meaning the oldest uncompleted task in the List
	Current() *Task

	// Complete a task by id
	Complete(id string) error
	// Update a task in the list, finding the original by the ID of the given task
	Update(Task) error
	// Add note to a task by ID
	AddNote(id string, note Note) error
}

// MemoryList implements List for a slice of Tasks
type MemoryList []Task

// TODO: Would counting matching tasks first be faster since we could
// frontload allocation?
func (ml *MemoryList) findWhere(f func(t Task) bool) []Task {
	new := []Task{}

	for _, t := range *ml {
		if f(t) {
			new = append(new, t)
		}
	}

	return new
}

func (ml *MemoryList) sort() {
	list := *ml

	sort.Slice(list, func(i, j int) bool {
		return (list[i].Priority > list[j].Priority) ||
			(list[i].Priority == list[j].Priority && list[i].CreatedDate.Before(list[j].CreatedDate))
	})

	ml = &list
}

// Completed returns a slice of completed tasks
func (ml *MemoryList) Completed(completed bool) []Task {
	return ml.findWhere(func(t Task) bool {
		return (t.CompletedDate.IsZero() && completed) ||
			(!t.CompletedDate.IsZero() && !completed)
	})
}

// Context returns a slice of tasks in the given context
func (ml *MemoryList) Context(context string) []Task {
	return ml.findWhere(func(t Task) bool {
		return t.Context == context
	})
}

// Add adds a task to this list
func (ml *MemoryList) Add(task Task) error {
	*ml = append(*ml, task)
	return nil
}

// AddMultiple adds multiple tasks to this list
func (ml *MemoryList) AddMultiple(tasks []Task) error {
	*ml = append(*ml, tasks...)
	return nil
}

// Slice turns the list into a slice of tasks
func (ml *MemoryList) Slice() []Task {
	return []Task(*ml)
}

// FindById will return a pointer to the task indicated by ID, nil if no task
// found with that id
func (ml *MemoryList) FindById(id string) *Task {
	for _, t := range *ml {
		if t.ID == id {
			return &t
		}
	}

	return nil
}

// Current will return the first task which is not completed
func (ml *MemoryList) Current() *Task {
	ml.sort()

	for _, t := range *ml {
		if t.CompletedDate.IsZero() {
			return &t
		}
	}

	return nil
}

// Complete will complete the task indicated by id
func (ml *MemoryList) Complete(id string) error {
	t := ml.FindById(id)
	if t == nil {
		return ErrNotFound
	}

	t.CompletedDate = time.Now()
	return nil
}

// Update will update the task indicated by the ID of the provided Task
func (ml *MemoryList) Update(other Task) error {
	t := ml.FindById(other.ID)
	if t == nil {
		return ErrNotFound
	}

	t.Context = other.Context
	t.Priority = other.Priority
	t.Body = other.Body
	t.Title = other.Title

	if len(other.Notes) > len(t.Notes) {
		t.Notes = append(t.Notes, other.Notes...)
	}

	return nil
}

// AddNote will add the note to the task indicated by the ID
func (ml *MemoryList) AddNote(id string, note Note) error {
	t := ml.FindById(id)
	if t == nil {
		return ErrNotFound
	}

	t.Notes = append(t.Notes, note)
	return nil
}
