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
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/ql/token"
)

// ErrNotFound is returned by a List when a Task matching the given ID does not
// exist
var ErrNotFound = errors.New("no task with that id exists")

// List is implemented by any struct that can maintain tasks
type List interface {
	// Evaluate the AST and return a List of matching results
	Search(ast ast.AST) ([]Task, error)
	// Add a task to the List
	Add(Task) error
	// Add multiple tasks to the List, should be more efficient resource
	// utilization.
	AddMultiple(task []Task) error
	// Return a slice of Tasks in this List
	Slice() []Task
	// Find a task by ID
	FindByID(id string) (Task, error)
	// Return the current task, meaning the oldest uncompleted task in the List
	Current() (Task, error)

	// Complete a task by id
	Complete(id string) error
	// Update a task in the list, finding the original by the ID of the given task
	Update(Task) error
	// Add note to a task by ID
	AddNote(id string, note Note) error
}

// MemoryList implements List for a slice of Tasks
type MemoryList []Task

// TODO: Would counting matching tasks first be faster since we
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

// FindByID will return a pointer to the task indicated by ID, returns
// ErrNotFound if no task with that ID exists
func (ml *MemoryList) FindByID(id string) (Task, error) {
	for _, t := range *ml {
		if t.ID == id {
			return t, nil
		}
	}

	return Task{}, ErrNotFound
}

// Current will return the first task which is not completed
func (ml *MemoryList) Current() (Task, error) {
	ml.sort()

	for _, t := range *ml {
		if t.CompletedDate.IsZero() {
			return t, nil
		}
	}

	return Task{}, ErrNotFound
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

// Update will update the task indicated by the ID of the provided Task
func (ml *MemoryList) Update(other Task) error {
	list := *ml

	for i := range list {
		if list[i].ID == other.ID {
			list[i].Context = other.Context
			list[i].Priority = other.Priority
			list[i].Body = other.Body
			list[i].Title = other.Title
			list[i].CompletedDate = other.CompletedDate

			if len(other.Notes) > len(list[i].Notes) {
				list[i].Notes = other.Notes
			}

			ml = &list
			return nil
		}
	}

	return ErrNotFound
}

// AddNote will add the note to the task indicated by the ID
func (ml *MemoryList) AddNote(id string, note Note) error {
	t, err := ml.FindByID(id)
	if err != nil {
		return err
	}

	t.Notes = append(t.Notes, note)
	return ml.Update(t)
}

// Search will construct a function to find matching tasks based on the given
// ast.AST
func (ml *MemoryList) Search(tree ast.AST) ([]Task, error) {
	return ml.findWhere(eval(tree.Expression)), nil
}

func eval(exp ast.Expression) func(t Task) bool {
	switch exp.(type) {
	case ast.InfixExpression:
		return evalInfixExp(exp.(ast.InfixExpression))
	case ast.StringLiteral:
		return evalStrLiteralQ(exp.(ast.StringLiteral))
	default:
		return func(t Task) bool { return false }
	}
}

func evalInfixExp(exp ast.InfixExpression) func(t Task) bool {
	switch exp.Operator.Type {
	case token.AND:
		return func(t Task) bool {
			return eval(exp.Left)(t) && eval(exp.Right)(t)
		}
	case token.OR:
		return func(t Task) bool {
			return eval(exp.Left)(t) || eval(exp.Right)(t)
		}
	case token.LIKE:
		return like(exp.Left, exp.Right)
	case token.NLIKE:
		return func(t Task) bool { return !like(exp.Left, exp.Right)(t) }
	case token.EQ:
		return eq(exp.Left, exp.Right)
	case token.NE:
		return func(t Task) bool { return !eq(exp.Left, exp.Right)(t) }
	case token.GT:
		return gt(exp.Left, exp.Right)
	case token.GTE:
		return gte(exp.Left, exp.Right)
	case token.LT:
		return func(t Task) bool { return !gt(exp.Left, exp.Right)(t) }
	case token.LTE:
		return lte(exp.Left, exp.Right)
	default:
		return func(t Task) bool { return false }
	}
}

func evalStrLiteralQ(exp ast.StringLiteral) func(t Task) bool {
	return func(t Task) bool {
		return strings.Contains(t.Title, exp.Value) ||
			strings.Contains(t.Body, exp.Value)
	}
}

func gt(left, right ast.Expression) func(Task) bool {
	return func(t Task) bool {
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

func gte(left, right ast.Expression) func(Task) bool {
	return func(t Task) bool {
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

func lte(left, right ast.Expression) func(Task) bool {
	return func(t Task) bool {
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

func eq(left, right ast.Expression) func(Task) bool {
	return func(t Task) bool {
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
			boolStr := right.(ast.StringLiteral).Value
			completed, err := strconv.ParseBool(boolStr)
			if err != nil {
				return false
			}

			return t.CompletedDate.IsZero() != completed
		default:
			return false
		}
	}
}

func like(left, right ast.Expression) func(Task) bool {
	return func(t Task) bool {
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
