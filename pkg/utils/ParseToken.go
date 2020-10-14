package utils

import (
	"errors"
	"strings"
)

func ParseToken(str string) (string, error) {
	// get token from session id
	authList := strings.Split(str, " ")
	if len(authList) != 2 {
		return "", errors.New("failed to parse token")
	}

	return authList[1], nil
}
