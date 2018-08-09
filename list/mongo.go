package list

import (
	"context"

	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/task"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

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
	mb.coll = mb.client.Database(mb.DB).Collection(mb.Collection)
	return err
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
	return task.Task{}, nil
}

// Complete a task by id
func (mb *MongoList) Complete(id string) error {
	return nil
}

// Update a task in the listist, finding the original by the ID of the given task
func (mb *MongoList) Update(task.Task) error {
	return nil
}

// AddNote to a task by ID
func (mb *MongoList) AddNote(id string, note task.Note) error {
	return nil
}

// Search will evaluate the AST and return a List of matching results
func (mb *MongoList) Search(ast ast.AST) ([]task.Task, error) {
	return nil, nil
}
