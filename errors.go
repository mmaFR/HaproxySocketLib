package HaproxySocketLib

import "errors"

var unknownConnectionType func(string) error = func(s string) error {
	return errors.New(s + ": unknown connection type")
}
