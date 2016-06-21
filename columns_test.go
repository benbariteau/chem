package chem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumns(t *testing.T) {
	cols := []Column{
		StringColumn{
			BaseColumn{
				Container: tweetTable,
				Name:      "text",
			},
		},
		IntegerColumn{
			BaseColumn{
				Container: tweetTable,
				Name:      "id",
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
					Container: tweetTable,
					Name:      "id",
				},
			},
			StringColumn{
				BaseColumn{
					Container: tweetTable,
					Name:      "text",
				},
			},
		},
	)
}
