package chem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectStmtFirst(t *testing.T) {
	db := setupDB()

	_, err := db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 1, "test tweet", 5)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)

	dbtweet := TweetRow{}
	err = Select(
		tweetTable,
	).Where(
		tweetTable.Id.Equals(1),
	).First(tx, &dbtweet)

	assert.NoError(t, err)
	assert.Equal(t, TweetRow{1, Tweet{"test tweet", 5}}, dbtweet)

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestSelectStmtAllOrderBy(t *testing.T) {
	db := setupDB()

	_, err := db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 1, "test tweet 1", 1)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 2, "test tweet 2", 5)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 3, "test tweet 3", 3)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)

	dbtweets := make([]TweetRow, 0, 3)
	err = Select(
		tweetTable,
	).OrderBy(
		tweetTable.Likes.Asc(),
	).All(tx, &dbtweets)

	assert.NoError(t, err)
	assert.Equal(
		t,
		[]TweetRow{
			TweetRow{1, Tweet{"test tweet 1", 1}},
			TweetRow{3, Tweet{"test tweet 3", 3}},
			TweetRow{2, Tweet{"test tweet 2", 5}},
		},
		dbtweets,
	)

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestSelectStmtOrderByLimit(t *testing.T) {
	db := setupDB()

	_, err := db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 1, "test tweet 1", 1)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 2, "test tweet 2", 5)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 3, "test tweet 3", 3)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)

	dbtweets := make([]TweetRow, 0, 3)
	err = Select(
		tweetTable,
	).OrderBy(
		tweetTable.Likes.Asc(),
	).Limit(2).All(tx, &dbtweets)

	// TODO this is fragile since it depends on ordering
	// either refactor to just care about values or just add ordering when it's implemented
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]TweetRow{
			TweetRow{1, Tweet{"test tweet 1", 1}},
			TweetRow{3, Tweet{"test tweet 3", 3}},
		},
		dbtweets,
	)

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestSelectStmtOrderByLimitOffset(t *testing.T) {
	db := setupDB()

	_, err := db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 1, "test tweet 1", 1)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 2, "test tweet 2", 5)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 3, "test tweet 3", 3)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)

	dbtweets := make([]TweetRow, 0, 3)
	err = Select(
		tweetTable,
	).OrderBy(
		tweetTable.Likes.Asc(),
	).Limit(2).Offset(2).All(tx, &dbtweets)

	// TODO this is fragile since it depends on ordering
	// either refactor to just care about values or just add ordering when it's implemented
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]TweetRow{
			TweetRow{2, Tweet{"test tweet 2", 5}},
		},
		dbtweets,
	)

	err = tx.Commit()
	assert.NoError(t, err)
}
