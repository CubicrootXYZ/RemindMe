package format

import (
	"errors"
	"strconv"
	"strings"
)

// GetSuffixInt returns a suffixed integer in the given string value
func GetSuffixInt(value string) (int, error) {
	splitUp := strings.Split(value, " ")
	if len(splitUp) == 0 {
		return 0, errors.New("empty string does not contain integer")
	}

	integerString := splitUp[len(splitUp)-1]

	return strconv.Atoi(integerString)
}
