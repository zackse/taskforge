package backends

import (
	"os"
	"testing"
)

func TestFileBackend(t *testing.T) {
	runBackendTest(t, backendTest{
		cleanup: func(b Backend) error {
			fileB := b.(*File)
			return os.RemoveAll(fileB.Dir)
		},
		config: Config{
			"dir": ".test.file.d",
		},
		backend: &File{},
		empty:   &File{},
	})
}
