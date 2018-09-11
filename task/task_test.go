package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUniqueIDs(t *testing.T) {
	t1 := New("task 1")
	t2 := New("task 2")
	require.NotEqual(t, t1.ID, t2.ID)
}

func TestComplete(t *testing.T) {
	t1 := New("task 1")

	var zeroTime time.Time
	require.Equal(t, t1.CompletedDate, zeroTime)
	t1.Complete()
	require.NotEqual(t, t1.CompletedDate, zeroTime)
	require.True(t, t1.IsComplete())
	require.True(t, t1.IsCompleted())
}
