package helper

import (
	"fmt"
	"net/url"
	"strings"
)

func IsValueMatchPattern(value, pattern string) bool {
	if len(pattern) == 0 {
		return false
	}
	if valuePattern, err := GetIdPattern(value); err == nil {
		return strings.EqualFold(valuePattern, pattern)
	}
	return false
}

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

func GetResourceType(id string) string {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return ""
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	resourceType := ""
	provider := ""
	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return ""
		}

		if key == "providers" {
			provider = value
			resourceType = provider
		} else if len(provider) > 0 {
			resourceType += "/" + key
		}
	}
	return resourceType
}
