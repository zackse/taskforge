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
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/ql/token"
	"github.com/chasinglogic/taskforge/task"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

var mongoSortOpt = findopt.Sort(bson.NewDocument(
	bson.EC.Int64("createdDate", -1),
))

// MongoList stores tasks in a MongoDB database
type MongoList struct {
	URL        string
	DB         string
	Collection string

	client *mongo.Client     `mapstructure:"-"`
	coll   *mongo.Collection `mapstructure:"-"`
}

// Init connects to MongoDB returning any errors
func (mb *MongoList) Init() error {
	if mb.URL == "" {
		mb.URL = "mongodb://localhost:27017"
	}

	if mb.DB == "" {
		mb.DB = "taskforge"
	}

	if mb.Collection == "" {
		mb.Collection = "tasks"
	}

	var err error
	mb.client, err = mongo.NewClient(mb.URL)
	if err != nil {
		return err
	}

	err = mb.client.Connect(context.Background())
	if err != nil {
		return err
	}

	mb.coll = mb.client.Database(mb.DB).Collection(mb.Collection)
	return nil
}

// Disconnect will disconnect the MongoDB client, returns an error if not connected
func (mb *MongoList) Disconnect() error {
	if mb.client == nil {
		return errors.New("mongo client not connected")
	}

	return mb.client.Disconnect(context.Background())
}

// Add a task to the List
func (mb *MongoList) Add(t task.Task) error {
	_, err := mb.coll.InsertOne(context.Background(), &t)
	return err
}

// AddMultiple tasks to the List, should be more efficient resource
// utilization.
func (mb *MongoList) AddMultiple(tasks []task.Task) error {
	docs := make([]interface{}, len(tasks))

	for i := range tasks {
		docs[i] = tasks[i]
	}

	_, err := mb.coll.InsertMany(context.Background(), docs)
	return err
}

// Slice returns a slice of Tasks in this List
func (mb *MongoList) Slice() ([]task.Task, error) {
	cur, err := mb.coll.Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	tasks := make([]task.Task, 0)

	for cur.Next(context.Background()) {
		var t task.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

// FindByID finds a task by ID
func (mb *MongoList) FindByID(id string) (task.Task, error) {
	query := bson.NewDocument(bson.EC.String("_id", id))
	var t task.Task
	err := mb.coll.FindOne(context.Background(), query).Decode(&t)
	return t, err
}

// Current returns the current task, meaning the oldest uncompleted task in the List
func (mb *MongoList) Current() (task.Task, error) {
	query := bson.NewDocument(
		bson.EC.SubDocument("completedDate", bson.NewDocument(
			bson.EC.Boolean("$exists", false),
		)),
	)

	var t task.Task
	err := mb.coll.FindOne(context.Background(), query, mongoSortOpt).Decode(&t)
	return t, err
}

// Complete a task by id
func (mb *MongoList) Complete(id string) error {
	now := time.Now()
	res, err := mb.coll.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("_id", id),
		),
		bson.NewDocument(
			bson.EC.SubDocument("$set", bson.NewDocument(
				bson.EC.Time("completedDate", now),
			)),
		),
	)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// Update a task in the listist, finding the original by the ID of the given task
func (mb *MongoList) Update(t task.Task) error {
	updateDoc := bson.NewDocument(
		bson.EC.SubDocument(
			"$set",
			bson.NewDocument(
				bson.EC.String("title", t.Title),
				bson.EC.String("body", t.Body),
				bson.EC.String("context", t.Context),
				bson.EC.Interface("priority", t.Priority),
			),
		),
	)

	res, err := mb.coll.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("_id", t.ID),
		),
		updateDoc,
	)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// AddNote to a task by ID
func (mb *MongoList) AddNote(id string, note task.Note) error {
	res, err := mb.coll.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("_id", id),
		),
		bson.NewDocument(
			bson.EC.SubDocument("$push", bson.NewDocument(
				bson.EC.Interface("notes", note),
			)),
		),
	)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// Search will evaluate the AST and return a List of matching results
func (mb *MongoList) Search(ast ast.AST) ([]task.Task, error) {
	query := mb.eval(ast.Expression)

	cur, err := mb.coll.Find(context.Background(), query, mongoSortOpt)
	if err != nil {
		return nil, err
	}

	var tasks []task.Task

	for cur.Next(context.Background()) {
		var t task.Task
		if err := cur.Decode(&t); err != nil {
			return tasks, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (mb *MongoList) eval(exp ast.Expression) *bson.Document {
	switch exp.(type) {
	case ast.InfixExpression:
		return mb.evalInfixExp(exp.(ast.InfixExpression))
	case ast.StringLiteral:
		return mb.evalStrLiteralQ(exp.(ast.StringLiteral))
	default:
		return bson.NewDocument()
	}
}

func (mb *MongoList) evalStrLiteralQ(exp ast.StringLiteral) *bson.Document {
	rgxDoc := bson.NewDocument(
		bson.EC.String("$regex", exp.Value),
		bson.EC.String("$options", "im"),
	)

	array := bson.NewArray()
	array.Append(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocument("title", rgxDoc),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocument("body", rgxDoc),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocument("notes", rgxDoc),
		),
	)

	return bson.NewDocument(bson.EC.Array("$or", array))
}

func (mb *MongoList) evalInfixExp(exp ast.InfixExpression) *bson.Document {
	switch exp.Operator.Type {
	case token.AND:
		array := bson.NewArray()
		array.Append(
			bson.VC.Document(mb.eval(exp.Left)),
			bson.VC.Document(mb.eval(exp.Right)),
		)
		return bson.NewDocument(bson.EC.Array("$and", array))
	case token.OR:
		array := bson.NewArray()
		array.Append(
			bson.VC.Document(mb.eval(exp.Left)),
			bson.VC.Document(mb.eval(exp.Right)),
		)
		return bson.NewDocument(bson.EC.Array("$or", array))
	case token.LIKE:
		return bson.NewDocument(
			bson.EC.SubDocument(
				exp.Left.(ast.StringLiteral).Value,
				bson.NewDocument(
					bson.EC.String("$regex", exp.Right.(ast.StringLiteral).Value),
				),
			),
		)
	case token.NLIKE:
		return bson.NewDocument(
			bson.EC.SubDocument(
				exp.Left.(ast.StringLiteral).Value,
				bson.NewDocument(
					bson.EC.String(
						"$regex",
						fmt.Sprintf(
							"((?!%s).)*",
							exp.Right.(ast.StringLiteral).Value,
						),
					),
				),
			),
		)
	case token.EQ:
		fieldName := exp.Left.(ast.StringLiteral).Value
		if fieldName == "completed" {
			var value time.Time
			if exp.Right.(ast.BooleanLiteral).Value {
				return bson.NewDocument(
					bson.EC.SubDocument(
						"completeddate",
						bson.NewDocument(
							bson.EC.Time(
								"$ne",
								value,
							),
						),
					),
				)
			}

			return bson.NewDocument(
				bson.EC.Time(
					"completeddate",
					value,
				),
			)
		}

		return bson.NewDocument(
			bson.EC.Interface(
				fieldName,
				exp.Right.(ast.Literal).GetValue(),
			),
		)
	case token.NE:
		fieldName := exp.Left.(ast.StringLiteral).Value
		if fieldName == "completed" {
			var value time.Time
			if exp.Right.(ast.BooleanLiteral).Value {
				return bson.NewDocument(
					bson.EC.Time(
						"completeddate",
						value,
					),
				)
			}

			return bson.NewDocument(
				bson.EC.SubDocument(
					"completeddate",
					bson.NewDocument(
						bson.EC.Time(
							"$ne",
							value,
						),
					),
				),
			)
		}

		return bson.NewDocument(
			bson.EC.SubDocument(
				fieldName,
				bson.NewDocument(
					bson.EC.Interface(
						"$ne",
						exp.Right.(ast.Literal).GetValue(),
					),
				),
			),
		)
	case token.GT:
		return bson.NewDocument(
			bson.EC.SubDocument(
				exp.Left.(ast.StringLiteral).Value,
				bson.NewDocument(
					bson.EC.Interface(
						"$gt",
						exp.Right.(ast.Literal).GetValue(),
					),
				),
			),
		)
	case token.GTE:
		return bson.NewDocument(
			bson.EC.SubDocument(
				exp.Left.(ast.StringLiteral).Value,
				bson.NewDocument(
					bson.EC.Interface(
						"$gte",
						exp.Right.(ast.Literal).GetValue(),
					),
				),
			),
		)
	case token.LT:
		return bson.NewDocument(
			bson.EC.SubDocument(
				exp.Left.(ast.StringLiteral).Value,
				bson.NewDocument(
					bson.EC.Interface(
						"$lt",
						exp.Right.(ast.Literal).GetValue(),
					),
				),
			),
		)
	case token.LTE:
		return bson.NewDocument(
			bson.EC.SubDocument(
				exp.Left.(ast.StringLiteral).Value,
				bson.NewDocument(
					bson.EC.Interface(
						"$lte",
						exp.Right.(ast.Literal).GetValue(),
					),
				),
			),
		)
	default:
		return bson.NewDocument()
	}
}
