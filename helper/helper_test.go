package helper

// Source for unit test: https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318
// change folder from main to cd helper before execute tests
import (
	"testing"
)

func TestIfFolderExist(t *testing.T) {
	var shouldExist = Exists("./testIfExist")
	var shouldNotExist = Exists("./testIfNotExist")

	if !shouldExist {
		t.Errorf("Existing folder ./testIfExist not found")
	}
	if shouldNotExist {
		t.Errorf("None Exsiting folder found")
	}

}
