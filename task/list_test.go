package task

import (
	"testing"
	"time"
)

func compareLists(list1, list2 MemoryList) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i := range list1 {
		if list1[i].ID != list2[i].ID {
			return false
		}
	}

	return true
}

func TestMemoryList_Completed(t *testing.T) {
	type args struct {
		completed bool
	}

	fixture := MemoryList{
		Task{Title: "completed 1", CompletedDate: time.Now()},
		Task{Title: "not completed"},
		Task{Title: "completed 2", CompletedDate: time.Now()},
	}

	tests := []struct {
		name string
		ml   MemoryList
		args args
		want []Task
	}{
		{
			name: "find completed",
			ml:   fixture,
			args: args{true},
			want: []Task{
				fixture[0],
				fixture[2],
			},
		},
		{
			name: "find incomplete",
			ml:   fixture,
			args: args{false},
			want: []Task{
				fixture[1],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ml.Completed(tt.args.completed)

			if !compareLists(got, tt.want) {
				t.Errorf("MemoryList.Completed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryList_Context(t *testing.T) {
	type args struct {
		context string
	}
	tests := []struct {
		name string
		ml   MemoryList
		args args
		want []Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ml.Context(tt.args.context)
			if len(got) != len(tt.want) {
				t.Errorf("MemoryList.Completed() = %v, want %v", got, tt.want)
				return
			}

			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("MemoryList.Completed() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestMemoryList_Add(t *testing.T) {
	list := MemoryList{
		New("task 1"),
		New("task 2"),
	}

	task3 := New("task 3")

	err := list.Add(task3)
	if err != nil {
		t.Error(err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 tasks got %d: %v", len(list), list)
		return
	}

	if list[2].ID != task3.ID {
		t.Errorf("Expected %v Got %v", task3, list[2])
	}
}

func TestMemoryList_AddMultiple(t *testing.T) {
	list := MemoryList{
		New("task 1"),
		New("task 2"),
	}

	task3 := New("task 3")
	task4 := New("task 4")

	other := []Task{
		task3,
		task4,
	}

	err := list.AddMultiple(other)
	if err != nil {
		t.Error(err)
	}

	if len(list) != 4 {
		t.Errorf("Expected 4 tasks got %d: %v", len(list), list)
		return
	}

	if list[2].ID != task3.ID {
		t.Errorf("Expected %v Got %v", task3, list[2])
	}

	if list[3].ID != task4.ID {
		t.Errorf("Expected %v Got %v", task4, list[3])
	}
}

func TestMemoryList_FindById(t *testing.T) {
	task2 := New("task 2")
	list := MemoryList{
		New("task 1"),
		task2,
		New("task 3"),
	}

	found := list.FindById(task2.ID)
	if task2.ID != found.ID || task2.Title != found.Title {
		t.Errorf("MemoryList.FindById() = %v, want %v", found, task2)
	}
}

func TestMemoryList_Current(t *testing.T) {
	tests := []struct {
		name string
		ml   MemoryList
		want *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ml.Current()

			if got == nil && tt.want != nil {
				t.Errorf("MemoryList.Current() = got nil when expected %v", tt.want)
				return
			}

			if tt.want.ID != got.ID {
				t.Errorf("MemoryList.Current() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryList_Complete(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		ml      MemoryList
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ml.Complete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("MemoryList.Complete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryList_Update(t *testing.T) {
	type args struct {
		other Task
	}
	tests := []struct {
		name    string
		ml      MemoryList
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ml.Update(tt.args.other); (err != nil) != tt.wantErr {
				t.Errorf("MemoryList.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryList_AddNote(t *testing.T) {
	type args struct {
		id   string
		note Note
	}
	tests := []struct {
		name    string
		ml      MemoryList
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ml.AddNote(tt.args.id, tt.args.note); (err != nil) != tt.wantErr {
				t.Errorf("MemoryList.AddNote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
