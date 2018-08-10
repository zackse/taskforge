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

	"github.com/chasinglogic/taskforge/ql/lexer"
	"github.com/chasinglogic/taskforge/ql/parser"
	"github.com/chasinglogic/taskforge/task"
)

type listTest struct {
	Name     string
	expected []int
	fixture  func() []task.Task
	run      func(List) error
}

func (bt listTest) Test(b List) error {
	var fixture []task.Task
	if bt.fixture != nil {
		fixture = bt.fixture()
		err := b.AddMultiple(fixture)
		if err != nil {
			return err
		}

	}

	if bt.run != nil {
		err := bt.run(b)
		if err != nil {
			return err
		}
	}

	if bt.expected != nil {
		received, err := b.Slice()
		if err != nil {
			return err
		}

		return verify(fixture, received, bt.expected)
	}

	return nil
}

var ListTests = []listTest{
	{
		Name: "should add one and find by id",
		run: func(b List) error {
			t := task.New("task 1")
			err := b.Add(t)
			if err != nil {
				return err
			}

			res, err := b.FindByID(t.ID)
			if err != nil {
				return err
			}

			return compareTasks(t, res)
		},
	},
	{
		Name: "should add multiple",
		fixture: func() []task.Task {
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
				task.New("task 3"),
			}
		},
		expected: []int{
			0,
			1,
			2,
		},
	},
	{
		Name: "should complete a task",
		fixture: func() []task.Task {
			return []task.Task{
				task.New("task to complete"),
			}
		},
		run: func(l List) error {
			tasks, err := l.Slice()
			if err != nil {
				return err
			}

			id := tasks[0].ID
			err = l.Complete(id)
			if err != nil {
				return err
			}

			t, err := l.FindByID(id)
			if err != nil {
				return err
			}

			if !t.IsCompleted() {
				return errors.New("expected task to be completed and was not")
			}

			return nil
		},
	},
	{
		Name: "should return correct current task",
		fixture: func() []task.Task {
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
			}
		},
		run: func(l List) error {
			tasks, err := l.Slice()
			if err != nil {
				return err
			}

			if len(tasks) == 0 {
				return errors.New("expected tasks got none")
			}

			current, err := l.Current()
			if err != nil {
				return err
			}

			if err := compareTasks(current, tasks[0]); err != nil {
				return err
			}

			err = l.Complete(tasks[0].ID)
			if err != nil {
				return err
			}

			current, err = l.Current()
			if err != nil {
				return err
			}

			if err := compareTasks(current, tasks[1]); err != nil {
				return err
			}

			return nil
		},
	},
	{
		Name: "should add a note to a task",
		fixture: func() []task.Task {
			return []task.Task{
				task.New("task to be noted"),
			}
		},
		run: func(l List) error {
			tasks, err := l.Slice()
			if err != nil {
				return err
			}

			note := task.NewNote("a note")
			if err := l.AddNote(tasks[0].ID, note); err != nil {
				return err
			}

			t, err := l.FindByID(tasks[0].ID)
			if err != nil {
				return err
			}

			if len(t.Notes) != 1 {
				return fmt.Errorf("expected 1 note got %d", len(t.Notes))
			}

			if t.Notes[0].Body != note.Body {
				return fmt.Errorf("expected %s got %s", note.Body, t.Notes[0].Body)
			}

			if t.Notes[0].ID != note.ID {
				return fmt.Errorf("expected %s got %s", note.ID, t.Notes[0].Body)
			}

			return nil
		},
	},
	{
		Name: "should update a task",
		fixture: func() []task.Task {
			return []task.Task{
				task.New("task to update"),
			}
		},
		run: func(l List) error {
			tasks, err := l.Slice()
			if err != nil {
				return err
			}

			toUpdate := tasks[0]
			toUpdate.Title = "task updated"
			err = l.Update(toUpdate)
			if err != nil {
				return err
			}

			updated, err := l.FindByID(tasks[0].ID)
			if err != nil {
				return err
			}

			if updated.Title != toUpdate.Title {
				return fmt.Errorf("expected %v got %v",
					toUpdate.String(), updated.String())
			}

			return nil
		},
	},
	queryTest(qt{
		query: "title = \"task 1\"",
		fixture: func() []task.Task {
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
			}
		},
		expected: []int{0},
	}),
	queryTest(qt{
		query: "context = other",
		fixture: func() []task.Task {
			other := task.New("other task")
			other.Context = "other"
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
				other,
				task.New("task 3"),
			}
		},
		expected: []int{2},
	}),
	queryTest(qt{
		query: "title = \"task 4\" or title = \"task 1\" or title = \"other task\"",
		fixture: func() []task.Task {
			other := task.New("other task")
			other.Context = "other"
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
				other,
				task.New("task 3"),
				task.New("task 4"),
			}
		},
		expected: []int{0, 2, 4},
	}),
	queryTest(qt{
		query: "(title = \"task 1\" and context = \"default\") or (context = \"other\")",
		fixture: func() []task.Task {
			other := task.New("other task")
			other.Context = "other"
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
				other,
				task.New("task 3"),
				task.New("task 4"),
			}
		},
		expected: []int{0, 2},
	}),
	queryTest(qt{
		query: "context = \"work\"",
		fixture: func() []task.Task {
			other := task.New("other task")
			other.Context = "other"
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
				other,
				task.New("task 3"),
				task.New("task 4"),
			}
		},
		expected: []int{},
	}),
	queryTest(qt{
		query: "task",
		fixture: func() []task.Task {
			other := task.New("other task")
			other.Context = "other"
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
				other,
				task.New("task 3"),
				task.New("task 4"),
			}
		},
		expected: []int{0, 1, 2, 3, 4},
	}),
	queryTest(qt{
		query: "priority = 0",
		fixture: func() []task.Task {
			other := task.New("other task")
			other.Context = "other"
			return []task.Task{
				task.New("task 1"),
				task.New("task 2"),
				other,
				task.New("task 3"),
				task.New("task 4"),
			}
		},
		expected: []int{0, 1, 2, 3, 4},
	}),
	queryTest(qt{
		query: "priority > 1",
		fixture: func() []task.Task {
			other := task.New("other task")
			other.Context = "other"
			other.Priority = 2
			t2 := task.New("task 2")
			t2.Priority = 1
			return []task.Task{
				task.New("task 1"),
				t2,
				other,
				task.New("task 3"),
				task.New("task 4"),
			}
		},
		expected: []int{2},
	}),
}

func verify(received, fixture []task.Task, expected []int) error {
	if len(received) != len(expected) {
		return fmt.Errorf("Expected %d results got %d", len(expected), len(received))
	}

	for i := range received {
		if err := compareTasks(received[i], fixture[expected[i]]); err != nil {
			return err
		}
	}

	return nil
}

func compareTasks(task1, task2 task.Task) error {
	if !(task1.Title == task2.Title && task1.ID == task2.ID) {
		return fmt.Errorf("Expected %s Got %s", task1.String(), task2.String())
	}

	return nil
}

func compareTaskSlice(tasks, other []task.Task) error {
	for i := range tasks {
		if err := compareTasks(tasks[i], other[i]); err != nil {
			return err
		}
	}

	return nil
}

type qt struct {
	query    string
	fixture  func() []task.Task
	expected []int
}

func queryTest(q qt) listTest {
	fixture := q.fixture()
	return listTest{
		Name: fmt.Sprintf("QUERYTEST: should return %d results with query: %s",
			len(q.expected), q.query),
		fixture: func() []task.Task { return fixture },
		run: func(l List) error {
			p := parser.New(lexer.New(q.query))
			query := p.Parse()

			if p.Error() != nil {
				return p.Error()
			}

			results, err := l.Search(query)
			if err != nil {
				return err
			}

			if len(results) != len(q.expected) {
				return fmt.Errorf("expected %d results got %d",
					len(q.expected), len(results))
			}

			return verify(results, fixture, q.expected)
		},
	}
}

// 	for _, test := range queryTests {
// 		t.Run(test.query, func(t *testing.T) {
// 			p := parser.New(lexer.New(test.query))
// 			ast := p.Parse()

// 			if err := p.Error(); err != nil {
// 				t.Errorf("error parsing query: %s", err)
// 				return
// 			}

// 			l, err := l.Search(ast)
// 			if err != nil {
// 				t.Errorf("error searching l: %s", err)
// 				return
// 			}

// 			if len(l) != len(test.expected) {
// 				t.Errorf("Expected %v Got %v", test.expected, l)
// 			}

// 		expectedTitles:
// 			for x := range test.expected {
// 				for y := range l {
// 					if l[y].Title == test.expected[x] {
// 						continue expectedTitles
// 					}
// 				}

// 				t.Errorf("Expected %v Got %v", test.expected, l)
// 			}
// 		})
// 	}
// }
