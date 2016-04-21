package chem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateStmtSet(t *testing.T) {
	db := setupDB()

	tx, err := db.Begin()
	assert.NoError(t, err)

	result, err := Insert(
		tweetTable,
	).Values(
		tx,
		Tweet{"test tweet", 3},
	)
	assert.NoError(t, err)
	id, err := result.LastInsertId()
	assert.NoError(t, err)

	result, err = Update(
		tweetTable,
	).Where(
		tweetTable.Id.Equals(int(id)),
	).Set(tx, map[Column]interface{}{
		tweetTable.Text:  "buttninjas",
		tweetTable.Likes: 666,
	})
	assert.NoError(t, err)
	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, int64(1))

	dbtweet := TweetRow{}
	err = Select(
		tweetTable,
	).Where(
		tweetTable.Id.Equals(int(id)),
	).First(tx, &dbtweet)
	assert.NoError(t, err)

	assert.Equal(t, dbtweet, TweetRow{1, Tweet{"buttninjas", 666}})

	err = tx.Commit()
	assert.NoError(t, err)
}
