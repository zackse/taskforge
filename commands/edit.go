package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/chasinglogic/tsk/task"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var edit = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit a task as YAML",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("TASK_ID is a required argument")
			os.Exit(1)
		}

		backend, err := config.backend()
		if err != nil {
			fmt.Println("ERROR Unable to load backend:", err)
			os.Exit(1)
		}

		tsk, err := backend.FindByID(args[0])
		if err != nil && err == task.ErrNotFound {
			fmt.Println(err)
			os.Exit(0)
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file, err := ioutil.TempFile("", "")
		if err != nil {
			fmt.Println("ERROR Unable to create temp file:", err)
			os.Exit(1)
		}

		defer file.Close()
		defer os.Remove(file.Name())

		yml, err := yaml.Marshal(tsk)
		if err != nil {
			fmt.Println("ERROR Unable to serialize task into yaml:", err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(file.Name(), yml, 0600)
		if err != nil {
			fmt.Println("ERROR Unable to generate temporary file:", err)
			os.Exit(1)
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		editorCmd := exec.Command(editor, file.Name())
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		err = editorCmd.Wait()
		if err != nil {
			fmt.Println("ERROR Unexpected error from editor:", err)
			os.Exit(1)
		}

		content, err := ioutil.ReadFile(file.Name())
		if err != nil {
			fmt.Println("ERROR Unable to read file:", err)
			os.Exit(1)
		}

		var updatedTask task.Task

		err = yaml.Unmarshal(content, &updatedTask)
		if err != nil {
			fmt.Println("ERROR Unable to parse yaml", err)
			os.Exit(1)
		}

		if err := backend.Update(updatedTask); err != nil {
			fmt.Println("ERROR Updating task:", err)
			os.Exit(1)
		}
	},
}
