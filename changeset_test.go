package okapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Changeset_Checksum(t *testing.T) {
	// ARRANGE
	c := Changeset{
		ID:     "001",
		Script: "CREATE TABLE test ( id INTEGER );",
	}

	// ACT
	md5 := c.checksum()

	// ASSERT
	assert.Equal(t, "48e8d99ae12125e25a65f9465c15468a", md5)
}

func Test_Changeset_Checksum_IgnoresWhitespace(t *testing.T) {
	// ARRANGE
	c1 := Changeset{
		ID:     "001",
		Script: "\nCREATE \t\tTABLE test (\n\tid    INTEGER\n);",
	}
	c2 := Changeset{
		ID:     "001",
		Script: "CREATE TABLE test ( id INTEGER );",
	}

	// ACT
	md5c1 := c1.checksum()
	md5c2 := c2.checksum()

	// ASSERT
	assert.Equal(t, md5c1, md5c2)
}

func Test_Changeset_ChecksumDiffersByID(t *testing.T) {
	// ARRANGE
	c1 := Changeset{
		ID:     "001",
		Script: "CREATE TABLE test ( id INTEGER );",
	}
	c2 := Changeset{
		ID:     "002",
		Script: "CREATE TABLE test ( id INTEGER );",
	}

	// ACT
	md5c1 := c1.checksum()
	md5c2 := c2.checksum()

	// ASSERT
	assert.NotEqual(t, md5c1, md5c2)
}

func Test_Changeset_ChecksumDiffersByScript(t *testing.T) {
	// ARRANGE
	c1 := Changeset{
		ID:     "001",
		Script: "CREATE TABLE test ( id INTEGER );",
	}
	c2 := Changeset{
		ID:     "001",
		Script: "CREATE TABLE test ( id2 INTEGER );",
	}

	// ACT
	md5c1 := c1.checksum()
	md5c2 := c2.checksum()

	// ASSERT
	assert.NotEqual(t, md5c1, md5c2)
}
