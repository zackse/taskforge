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
