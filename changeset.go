package okapi

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// Changeset specifies a specific migration step
type Changeset struct {
	ID          string
	SkipOnError bool
	Script      string
	Comment     string
}

// checksum calculates a checksum for a more simplified version of the actual
// script be ignoring all whitespace, so formatting changes won't break migrations
func (c Changeset) checksum() string {
	sanitized := strings.Join(strings.Fields(c.Script), " ")
	unique := fmt.Sprintf("%s::%s", c.ID, sanitized)
	return fmt.Sprintf("%x", md5.Sum([]byte(unique)))
}
