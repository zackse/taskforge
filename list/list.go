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
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/ql/token"
	"github.com/chasinglogic/taskforge/task"
)

// ErrNotFound is returned by a List when a task.Task matching the given ID does not
// exist
var ErrNotFound = errors.New("no task with that id exists")

// List is implemented by any struct that can maintain tasks
type List interface {
	// Init performs any initialization required for a List that does IO
	// such as creating files or connecting to databases.
	Init() error
	// Evaluate the AST and return a List of matching results
	Search(ast ast.AST) ([]task.Task, error)
	// Add a task to the List
	Add(task.Task) error
	// Add multiple tasks to the List, should be more efficient resource
	// utilization.
	AddMultiple(task []task.Task) error
	// Return a slice of task.Tasks in this List
	Slice() ([]task.Task, error)
	// Find a task by ID
	FindByID(id string) (task.Task, error)
	// Return the current task, meaning the oldest uncompleted task in the List
	Current() (task.Task, error)

	// Complete a task by id
	Complete(id string) error
	// Update a task in the listist, finding the original by the ID of the given task
	Update(task.Task) error
	// Add note to a task by ID
	AddNote(id string, note task.Note) error
}

// Config is a convenience type used to represent the unstructured config that
// will be stored by the application. All lists should accept being decoded
// using mapstructure from this type. Additionally, the Init() function should
// return an error if a required argument was not provided.
type Config map[string]interface{}

// GetByName returns the appropriate List implementation by
// a human redable string name. It returns a user friendly error
// string if not found.
func GetByName(name string) (List, error) {
	switch name {
	case "file":
		return &File{}, nil
	case "mongodb":
		fallthrough
	case "mongo":
		return &MongoList{}, nil
	default:
		return nil, fmt.Errorf("no l found with name %s", name)
	}
}

// MemoryList implements List for a slice of Tasks
type MemoryList []task.Task

// TODO: Would counting matching tasks first be faster since we
// frontload allocation?
func (ml *MemoryList) findWhere(f func(t task.Task) bool) []task.Task {
	new := []task.Task{}

	for _, t := range *ml {
		if f(t) {
			new = append(new, t)
		}
	}

	return new
}

func (ml *MemoryList) sort() {
	l := *ml

	sort.Slice(l, func(i, j int) bool {
		return (l[i].Priority > l[j].Priority) ||
			(l[i].Priority == l[j].Priority && l[i].CreatedDate.Before(l[j].CreatedDate))
	})

	ml = &l
}

// Add adds a task to this l
func (ml *MemoryList) Add(task task.Task) error {
	*ml = append(*ml, task)
	return nil
}

// AddMultiple adds multiple tasks to this l
func (ml *MemoryList) AddMultiple(tasks []task.Task) error {
	*ml = append(*ml, tasks...)
	return nil
}

// Slice turns the l into a slice of tasks
func (ml *MemoryList) Slice() ([]task.Task, error) {
	return []task.Task(*ml), nil
}

// FindByID will return a pointer to the task indicated by ID, returns
// ErrNotFound if no task with that ID exists
func (ml *MemoryList) FindByID(id string) (task.Task, error) {
	for _, t := range *ml {
		if t.ID == id {
			return t, nil
		}
	}

	return task.Task{}, ErrNotFound
}

// Current will return the first task which is not completed
func (ml *MemoryList) Current() (task.Task, error) {
	ml.sort()

	for _, t := range *ml {
		if t.CompletedDate.IsZero() {
			return t, nil
		}
	}

	return task.Task{}, ErrNotFound
}

// Complete will complete the task indicated by id
func (ml *MemoryList) Complete(id string) error {
	t, err := ml.FindByID(id)
	if err != nil {
		return err
	}

	t.CompletedDate = time.Now()
	return ml.Update(t)
}

// Update will update the task indicated by the ID of the provided task.Task
func (ml *MemoryList) Update(other task.Task) error {
	l := *ml

	for i := range l {
		if l[i].ID == other.ID {
			l[i].Context = other.Context
			l[i].Priority = other.Priority
			l[i].Body = other.Body
			l[i].Title = other.Title
			l[i].CompletedDate = other.CompletedDate

			if len(other.Notes) > len(l[i].Notes) {
				l[i].Notes = other.Notes
			}

			ml = &l
			return nil
		}
	}

	return ErrNotFound
}

// AddNote will add the note to the task indicated by the ID
func (ml *MemoryList) AddNote(id string, note task.Note) error {
	t, err := ml.FindByID(id)
	if err != nil {
		return err
	}

	t.Notes = append(t.Notes, note)
	return ml.Update(t)
}

// Search will construct a function to find matching tasks based on the given
// ast.AST
func (ml *MemoryList) Search(tree ast.AST) ([]task.Task, error) {
	return ml.findWhere(eval(tree.Expression)), nil
}

func eval(exp ast.Expression) func(t task.Task) bool {
	switch exp.(type) {
	case ast.InfixExpression:
		return evalInfixExp(exp.(ast.InfixExpression))
	case ast.StringLiteral:
		return evalStrLiteralQ(exp.(ast.StringLiteral))
	default:
		return func(t task.Task) bool { return false }
	}
}

func evalInfixExp(exp ast.InfixExpression) func(t task.Task) bool {
	switch exp.Operator.Type {
	case token.AND:
		return func(t task.Task) bool {
			return eval(exp.Left)(t) && eval(exp.Right)(t)
		}
	case token.OR:
		return func(t task.Task) bool {
			return eval(exp.Left)(t) || eval(exp.Right)(t)
		}
	case token.LIKE:
		return like(exp.Left, exp.Right)
	case token.NLIKE:
		return func(t task.Task) bool { return !like(exp.Left, exp.Right)(t) }
	case token.EQ:
		return eq(exp.Left, exp.Right)
	case token.NE:
		return func(t task.Task) bool { return !eq(exp.Left, exp.Right)(t) }
	case token.GT:
		return gt(exp.Left, exp.Right)
	case token.GTE:
		return gte(exp.Left, exp.Right)
	case token.LT:
		return func(t task.Task) bool { return !gt(exp.Left, exp.Right)(t) }
	case token.LTE:
		return lte(exp.Left, exp.Right)
	default:
		return func(t task.Task) bool { return false }
	}
}

func evalStrLiteralQ(exp ast.StringLiteral) func(t task.Task) bool {
	return func(t task.Task) bool {
		return strings.Contains(t.Title, exp.Value) ||
			strings.Contains(t.Body, exp.Value)
	}
}

func gt(left, right ast.Expression) func(task.Task) bool {
	return func(t task.Task) bool {
		switch left.(ast.StringLiteral).Value {
		case "priority":
			return t.Priority > right.(ast.NumberLiteral).Value
		case "created_date":
			fallthrough
		case "createdDate":
			return t.CreatedDate.After(right.(ast.DateLiteral).Value)
		case "completed_date":
			fallthrough
		case "completedDate":
			return t.CompletedDate.After(right.(ast.DateLiteral).Value)
		default:
			return false
		}
	}
}

func gte(left, right ast.Expression) func(task.Task) bool {
	return func(t task.Task) bool {
		switch left.(ast.StringLiteral).Value {
		case "priority":
			return t.Priority >= right.(ast.NumberLiteral).Value
		case "created_date":
			fallthrough
		case "createdDate":
			return t.CreatedDate.After(right.(ast.DateLiteral).Value) ||
				t.CreatedDate == right.(ast.DateLiteral).Value
		case "completed_date":
			fallthrough
		case "completedDate":
			return t.CompletedDate.After(right.(ast.DateLiteral).Value) ||
				t.CompletedDate == right.(ast.DateLiteral).Value
		default:
			return false
		}
	}
}

func lte(left, right ast.Expression) func(task.Task) bool {
	return func(t task.Task) bool {
		switch left.(ast.StringLiteral).Value {
		case "priority":
			return t.Priority <= right.(ast.NumberLiteral).Value
		case "created_date":
			fallthrough
		case "createdDate":
			return t.CreatedDate.Before(right.(ast.DateLiteral).Value) ||
				t.CreatedDate == right.(ast.DateLiteral).Value
		case "completed_date":
			fallthrough
		case "completedDate":
			return t.CompletedDate.Before(right.(ast.DateLiteral).Value) ||
				t.CompletedDate == right.(ast.DateLiteral).Value
		default:
			return false
		}
	}
}

func eq(left, right ast.Expression) func(task.Task) bool {
	return func(t task.Task) bool {
		switch left.(ast.StringLiteral).Value {
		case "title":
			return t.Title == right.(ast.StringLiteral).Value
		case "body":
			return t.Body == right.(ast.StringLiteral).Value
		case "context":
			return t.Context == right.(ast.StringLiteral).Value
		case "priority":
			return t.Priority == right.(ast.NumberLiteral).Value
		case "created_date":
			fallthrough
		case "createdDate":
			return t.CreatedDate == right.(ast.DateLiteral).Value
		case "completed_date":
			fallthrough
		case "completedDate":
			return t.CompletedDate == right.(ast.DateLiteral).Value
		case "completed":
			completed := right.(ast.BooleanLiteral).Value
			return t.IsCompleted() == completed
		default:
			return false
		}
	}
}

func like(left, right ast.Expression) func(task.Task) bool {
	return func(t task.Task) bool {
		switch left.(ast.StringLiteral).Value {
		case "title":
			return strings.Contains(t.Title, right.(ast.StringLiteral).Value)
		case "body":
			return strings.Contains(t.Body, right.(ast.StringLiteral).Value)
		case "context":
			return strings.Contains(t.Context, right.(ast.StringLiteral).Value)
		default:
			return false
		}
	}
}
