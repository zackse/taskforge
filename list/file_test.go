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
	"fmt"
	"os"
	"testing"

	"github.com/mitchellh/mapstructure"
)

func TestFileList(t *testing.T) {
	for _, test := range listTests {
		config := Config{
			"dir": fmt.Sprintf(".test.file.%s", test.name),
		}

		l := &File{}

		err := mapstructure.Decode(config, &l)
		if err != nil {
			t.Error(err)
			return
		}

		if err := l.Init(); err != nil {
			t.Errorf("unable to init list: %s", err)
			return
		}

		t.Run(test.name, func(t *testing.T) {
			defer os.RemoveAll(fmt.Sprintf(".test.file.%s", test.name))

			err := test.Test(l)
			if err != nil {
				t.Error(err)
				return
			}

		})
	}
}
