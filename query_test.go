package chem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicSelectQueryString(t *testing.T) {
	expected := "SELECT id, name, is_admin FROM user"
	got, err := BasicSelect{
		columns: []string{"id", "name", "is_admin"},
		table:   "user",
	}.QueryString()
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
