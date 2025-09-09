
package util

import "github.com/microcosm-cc/bluemonday"

// SanitizeHTML removes potentially harmful HTML from a string,
// preventing XSS attacks.
func SanitizeHTML(html string) string {
	// Use the UGC policy as a starting point. It allows a good range of
	// HTML elements and attributes, but strips anything dangerous.
	p := bluemonday.UGCPolicy()
	return p.Sanitize(html)
}
