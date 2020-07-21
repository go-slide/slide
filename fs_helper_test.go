package ferry

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestAllPaths(t *testing.T) {
	var paths []string
	err := getAllPaths("example", &paths)
	if err != nil {
		t.Errorf("error while getting paths %s", err.Error())
	}
	assert.Equal(t, len(paths), 1)
}
