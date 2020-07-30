package slide

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllPaths(t *testing.T) {
	var paths []string
	err := getAllPaths(".github", &paths)
	if err != nil {
		t.Errorf("error while getting paths %s", err.Error())
	}
	assert.Equal(t, len(paths), 3)
}

func TestAllPathsError(t *testing.T) {
	var paths []string
	err := getAllPaths("lol", &paths)
	assert.NotNil(t, err)
}
