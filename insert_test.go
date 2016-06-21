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
			Container: TweetTableType{},
			Name:      "id",
		},
	},
	Text: StringColumn{
		BaseColumn{
			Container: TweetTableType{},
			Name:      "text",
		},
	},
	Likes: IntegerColumn{
		BaseColumn{
			Container: TweetTableType{},
			Name:      "likes",
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
	assert.NoError(t, err)

	result, err := stmt.Values(tx, Tweet{Text: "test tweet"})
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	id, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	row := db.QueryRow("SELECT text, likes FROM tweet WHERE text = ?", "test tweet")
	dbtweet := Tweet{}
	err = row.Scan(&dbtweet.Text, &dbtweet.Likes)
	assert.NoError(t, err)
	assert.Equal(t, Tweet{Text: "test tweet", Likes: 0}, dbtweet)
}
