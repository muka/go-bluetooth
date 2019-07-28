package api

import "strings"

// UUIDToPath convert an UUID to a DBus path
func UUIDToPath(uuid string) string {
	return strings.Replace(uuid, "-", "_", -1)
}
