package helper

import (
	"os/user"
	"strings"
)

// get local user, separated based on testing issues
func GetUser() *user.User {
	localUser, err := user.Current()

	if err != nil {
		panic(err)
	}
	return localUser
}

// check if it is a valid user
func CheckLicence(localUser *user.User) bool {
	validLicence := false

	if strings.Contains(localUser.Username, "wenzela") || strings.Contains(localUser.Username, "atlabor") {
		validLicence = true
	}

	return validLicence
}
