package types

import (
	"fmt"
	"net/url"
	"strings"
)

func GetIdPattern(id string) (string, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return "", fmt.Errorf("cannot parse Azure ID, id: %s: %+v", id, err)
	}
	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	pattern := ""
	for current := 0; current <= len(components)-2; current += 2 {
		key := components[current]
		value := components[current+1]
		// Check key/value for empty strings.
		if key == "" || value == "" {
			return "", fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s', id: %s", key, value, id)
		}
		pattern += "/" + key
		if key == "providers" {
			// TODO: add validation on value, it should be like Microsoft.Something and case sensitive
			pattern += "/" + value
		}
	}

	return pattern, nil
}
