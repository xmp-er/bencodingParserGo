package validators

import (
	"strconv"
)

func IsValidInp(str []string) bool {
	return len(str) == 1
}

func IsValidString(str []string) bool {
	//string array length must be 3 in the format "<int>:<string>"

	//checking if the length is three
	if len(str) != 2 {
		return false
	}

	//checking if the str[0] is int or not
	if _, err := strconv.Atoi(str[0]); err != nil {
		return false
	}

	//since all the validations are correct, returning the string

	return true

}

func IsValidInt(str string) bool {
	// if the first digit is not i and the last digit is not e and the string between them is not a integer, return false
	//tried something funny kek
	if str[0] != 'i' || str[len(str)-1] != 'e' || (func() bool {
		_, err := strconv.Atoi(str[1 : len(str)-1])
		return err != nil
	}()) {
		return false
	}
	return true
}
