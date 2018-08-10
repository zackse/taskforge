package list_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/chasinglogic/taskforge/list"
	"github.com/chasinglogic/taskforge/server"
	"github.com/mitchellh/mapstructure"
)

func TestServerList(t *testing.T) {
	port := 9090
	for _, test := range list.ListTests {
		t.Run(test.Name, func(t *testing.T) {
			fileList := &list.File{
				Dir: fmt.Sprintf(".test.server.%s", test.Name),
			}

			defer os.RemoveAll(fileList.Dir)

			s := server.New(fileList, "testToken")
			s.Port = port
			go func() {
				fmt.Println("starting server")
				s.Listen()
			}()

			fmt.Println("server started")

			config := list.Config{
				"serverURL": fmt.Sprintf("http://localhost:%d", port),
				"token":     "testToken",
			}

			l := &list.ServerList{}
			err := mapstructure.Decode(config, &l)
			if err != nil {
				t.Error(err)
				return
			}

			defer s.Shutdown()
			if err := l.Init(); err != nil {
				t.Errorf("unable to init list: %s", err)
				return
			}

			err = test.Test(l)
			if err != nil {
				t.Error(err)
				return
			}
		})

		port++
	}
}
