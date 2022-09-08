package engine

import "strings"

// NameToTap returns the connected tap interface of an interface.
func NameToTap(name string) string {
	if IsTap(name) {
		return name
	}
	return name + "-tap"
}

// GetNameFromTap returns connected external interface from a tap interface.
func GetNameFromTap(name string) string {
	return strings.TrimSuffix(name, "-tap")
}

// IsTap returns wether the interface is a tap interface.
func IsTap(name string) bool {
	return strings.HasSuffix(name, "-tap")
}
