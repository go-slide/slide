package ferry

import (
	"testing"
)

func TestAllPaths(t *testing.T) {
	var paths []string
	err := getAllPaths("./", &paths)
	if err != nil {
		t.Errorf("error while getting paths %s", err.Error())
	}
}