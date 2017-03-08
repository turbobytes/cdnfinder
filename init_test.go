package cdnfinder

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	Init()
	if len(cdnmatches) == 0 {
		t.Errorf("Expected length > 0")
	}
	if _, err := os.Stat(resourcefinderjs); os.IsNotExist(err) {
		t.Error(err)
	}
	stat, err := os.Stat(phantomjsbin)
	if os.IsNotExist(err) {
		t.Error(err)
	}
	if stat.Size() == 0 {
		t.Errorf("Binary size is 0")
	}
	//TODO: Other checks
}
