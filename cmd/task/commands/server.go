package commands

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chasinglogic/taskforge/server"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd.AddCommand(serverRun)
	serverCmd.AddCommand(serverGenToken)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage a task server",
}

var serverRun = &cobra.Command{
	Use:   "run",
	Short: "Run a task server",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := config.Server.list()
		if err != nil {
			fmt.Println("ERROR Unable to load list:", err)
			os.Exit(1)
		}

		tokens, err := readTokenFile()
		if err != nil && os.IsNotExist(err) {
			fmt.Println("ERROR you need to generate a token before running the task server")
			os.Exit(1)
		} else if err != nil {
			fmt.Println("ERROR unable to load tokens:", err)
			os.Exit(1)
		}

		s := server.New(l, tokens...)

		if config.Server.Addr != "" {
			s.Addr = config.Server.Addr
		}

		if config.Server.Port != 0 {
			s.Port = config.Server.Port
		}

		s.Listen()
	},
}

var serverGenToken = &cobra.Command{
	Use:   "gen-token",
	Short: "Generate an API authentication token for this task server",
	Run: func(cmd *cobra.Command, args []string) {
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			fmt.Println("ERROR generating token", err)
			os.Exit(1)
		}

		token := fmt.Sprintf("%x", b)
		err = saveTokenFile([]string{token})
		if err != nil {
			fmt.Println("ERROR unable to save token file:", err)
			os.Exit(1)
		}

		fmt.Println("successfully created token:", token)
	},
}

func tokenFile() string {
	dir := defaultDir()
	return filepath.Join(dir, "server.state.json")
}

func readTokenFile() ([]string, error) {
	content, err := ioutil.ReadFile(tokenFile())
	if err != nil {
		return nil, err
	}

	var tokens []string
	err = json.Unmarshal(content, &tokens)
	return tokens, err
}

func saveTokenFile(newTokens []string) error {
	tokens, err := readTokenFile()
	if err != nil && os.IsNotExist(err) {
		tokens = []string{}
	} else if err != nil {
		return err
	}

	tokens = append(tokens, newTokens...)
	jsn, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(tokenFile(), jsn, 0600)
}
