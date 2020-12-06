package helper

// Source for unit test: https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318
// change folder from main to cd helper before execute tests
import (
	"os/user"
	"testing"
)

func TestUserValidationSuccess(t *testing.T) {

	var myTestUser *user.User = &(user.User{})

	myTestUser.Username = "wenzela"
	validLicence1 := CheckLicence(myTestUser)

	myTestUser.Username = "atlabor"
	validLicence2 := CheckLicence(myTestUser)

	if !validLicence1 {
		t.Errorf("ValidationSuccess failed, wzl not accept")
	}
	if !validLicence2 {
		t.Errorf("ValidationSuccess failed, at not accept")
	}

}

func TestUserValidationFailure(t *testing.T) {

	var myTestUser *user.User = &(user.User{})

	myTestUser.Username = "hase"
	validLicence1 := CheckLicence(myTestUser)

	myTestUser.Username = "peter"
	validLicence2 := CheckLicence(myTestUser)

	if validLicence1 {
		t.Errorf("ValidationFailure failed, hans was accept")
	}
	if validLicence2 {
		t.Errorf("ValidationFailure failed, peter was accept")
	}

}

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
