
package util

import (
	"strings"

	"github.com/gosimple/slug"
	"github.com/gosimple/unidecode"
)

// Slugify converts a string to a slug.
func Slugify(s string) string {
	return slug.Make(s)
}

// Normalize removes diacritics from a string and converts it to lowercase.
func Normalize(s string) string {
	return strings.ToLower(unidecode.Unidecode(s))
}
