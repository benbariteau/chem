package chem

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type tweet struct {
	Text  string
	Likes uint
}

type tweetTableType struct {
	Id    IntegerColumn
	Text  BaseColumn
	Likes BaseColumn
}

func (_ tweetTableType) Name() string {
	return "tweet"
}

func (_ tweetTableType) Type() reflect.Type {
	return reflect.TypeOf(tweet{})
}

func (t tweetTableType) Columns() []Column {
	return []Column{
		t.Id,
		t.Text,
		t.Likes,
	}
}

var tweetTable = tweetTableType{
	Id: IntegerColumn{
		BaseColumn{
			table: tweetTableType{},
			name:  "id",
		},
	},
	Text: BaseColumn{
		table: tweetTableType{},
		name:  "text",
	},
	Likes: BaseColumn{
		table: tweetTableType{},
		name:  "likes",
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

	result, err := stmt.Values(tx, tweet{Text: "test tweet"})
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
	dbtweet := tweet{}
	err = row.Scan(&dbtweet.Text, &dbtweet.Likes)
	assert.Nil(t, err)
	assert.Equal(t, dbtweet, tweet{Text: "test tweet", Likes: 0})
}

func TestSelectStmtOne(t *testing.T) {
	db := setupDB()

	_, err := db.Exec("INSERT INTO tweet (id, text, likes) VALUES (?, ?, ?)", 1, "test tweet", 5)
	assert.Nil(t, err)

	tx, err := db.Begin()
	assert.Nil(t, err)

	var id int
	dbtweet := tweet{}
	err = Select(
		tweetTable,
	).Where(
		tweetTable.Id.Equals(1),
	).One(tx, &id, &dbtweet.Text, &dbtweet.Likes)

	assert.Nil(t, err)
	assert.Equal(t, id, 1)
	assert.Equal(t, dbtweet, tweet{"test tweet", 5})

	err = tx.Commit()
	assert.Nil(t, err)
}
