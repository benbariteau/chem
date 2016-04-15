package chem

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type Tweet struct {
	Text  string
	Likes uint
}

type TweetRow struct {
	Id int
	Tweet
}

type TweetTableType struct {
	Id    IntegerColumn
	Text  StringColumn
	Likes IntegerColumn
}

func (_ TweetTableType) Name() string {
	return "tweet"
}

func (_ TweetTableType) Type() reflect.Type {
	return reflect.TypeOf(Tweet{})
}

func (t TweetTableType) Columns() []Column {
	return []Column{
		t.Id,
		t.Text,
		t.Likes,
	}
}

var tweetTable = TweetTableType{
	Id: IntegerColumn{
		BaseColumn{
			table: TweetTableType{},
			name:  "id",
		},
	},
	Text: StringColumn{
		BaseColumn{
			table: TweetTableType{},
			name:  "text",
		},
	},
	Likes: IntegerColumn{
		BaseColumn{
			table: TweetTableType{},
			name:  "likes",
		},
	},
}

const tweetTableDef = `
CREATE TABLE tweet (
	id INTEGER PRIMARY KEY,
	text TEXT NOT NULL,
	likes INTEGER NOT NULL DEFAULT 0
);`

func setupDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(tweetTableDef)
	if err != nil {
		panic(err)
	}
	return db
}

func TestInsertStmtValues(t *testing.T) {
	db := setupDB()

	stmt := Insert(tweetTable)

	tx, err := db.Begin()
	assert.Nil(t, err)

	result, err := stmt.Values(tx, Tweet{Text: "test tweet"})
	assert.Nil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	id, err := result.LastInsertId()
	assert.Nil(t, err)
	assert.Equal(t, id, int64(1))

	rowsAffected, err := result.RowsAffected()
	assert.Nil(t, err)
	assert.Equal(t, rowsAffected, int64(1))

	row := db.QueryRow("SELECT text, likes FROM tweet WHERE text = ?", "test tweet")
	dbtweet := Tweet{}
	err = row.Scan(&dbtweet.Text, &dbtweet.Likes)
	assert.Nil(t, err)
	assert.Equal(t, dbtweet, Tweet{Text: "test tweet", Likes: 0})
}

func TestSelectStmtFirst(t *testing.T) {
	db := setupDB()

	_, err := db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 1, "test tweet", 5)
	assert.Nil(t, err)

	tx, err := db.Begin()
	assert.Nil(t, err)

	dbtweet := TweetRow{}
	err = Select(
		tweetTable,
	).Where(
		tweetTable.Id.Equals(1),
	).First(tx, &dbtweet)

	assert.Nil(t, err)
	assert.Equal(t, dbtweet, TweetRow{1, Tweet{"test tweet", 5}})

	err = tx.Commit()
	assert.Nil(t, err)
}

func TestSelectStmtAll(t *testing.T) {
	db := setupDB()

	_, err := db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 1, "test tweet 1", 1)
	assert.Nil(t, err)
	_, err = db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 2, "test tweet 2", 2)
	assert.Nil(t, err)
	_, err = db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 3, "test tweet 3", 3)
	assert.Nil(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)

	dbtweets := make([]TweetRow, 0, 3)
	err = Select(
		tweetTable,
	).All(tx, &dbtweets)

	// TODO this is fragile since it depends on ordering
	// either refactor to just care about values or just add ordering when it's implemented
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]TweetRow{
			TweetRow{1, Tweet{"test tweet 1", 1}},
			TweetRow{2, Tweet{"test tweet 2", 2}},
			TweetRow{3, Tweet{"test tweet 3", 3}},
		},
		dbtweets,
	)

	err = tx.Commit()
	assert.Nil(t, err)
}
