package chem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumns(t *testing.T) {
	cols := []Column{
		StringColumn{
			BaseColumn{
				table: tweetTable,
				name:  "text",
			},
		},
		IntegerColumn{
			BaseColumn{
				table: tweetTable,
				name:  "id",
			},
		},
	}

	original := cols[:]

	sortedCols := sortColumns(cols)

	assert.Equal(t, cols, original)
	assert.Equal(
		t,
		sortedCols,
		[]Column{
			IntegerColumn{
				BaseColumn{
					table: tweetTable,
					name:  "id",
				},
			},
			StringColumn{
				BaseColumn{
					table: tweetTable,
					name:  "text",
				},
			},
		},
	)
}
