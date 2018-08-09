package backends

import "github.com/mongodb/mongo-go-driver/mongo"

type MongoBackend struct {
	Url string	
	DB string
	Collection string

	client *mongo.Client `mapstructure:"-"`
}

func (mb *MongoBackend) Init() error {
	if mb.Url == "" {
		mb.Url = "mongodb://localhost:27017"
	}

	if mb.DB == "" {
		mb.DB = "taskforge"
	}

	if mb.Collection == "" {
		mb.Collection = "tasks"
	}

	mb.client, err := mongo.NewClient(mb.Url)
	return err
}

func (mb *MongoBackend) Save() error {
	return nil
}

func (mb *MongoBackend) Load() error {
	return nil
}

// Evaluate the AST and return a List of matching results
func (mb *MongoBackend) Search(ast ast.AST) ([]Task, error) {

}

// Add a task to the List
func (mb *MongoBackend) Add(Task) error

// Add multiple tasks to the List, should be more efficient resource
// utilization.
func (mb *MongoBackend) AddMultiple(task []Task) error

// Return a slice of Tasks in this List
func (mb *MongoBackend) Slice() []Task

// Find a task by ID
func (mb *MongoBackend) FindByID(id string) (Task, error)

// Return the current task, meaning the oldest uncompleted task in the List
func (mb *MongoBackend) Current() (Task, error)

// Complete a task by id
func (mb *MongoBackend) Complete(id string) error

// Update a task in the list, finding the original by the ID of the given task
func (mb *MongoBackend) Update(Task) error

// Add note to a task by ID
func (mb *MongoBackend) AddNote(id string, note Note) error
