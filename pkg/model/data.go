package model

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const DB_PATH string = "../../assets/leglog_v2.db"

func init() {
	db, _ = sql.Open("sqlite3", DB_PATH)
	if _, err := os.Stat(DB_PATH); os.IsNotExist(err) {
		os.Create(DB_PATH)
		db.Exec(`
		CREATE TABLE entry(
			id        INTEGER PRIMARY KEY     NOT NULL,
			text      TEXT    NOT NULL,
			language        CHAR(3),
			creation_time NUMERIC NOT NULL
	 	);
		CREATE TABLE note(
			id INTEGER PRIMARY KEY NOT NULL,
			entry_id INTEGER NOT NULL REFERENCES entry(id) ON DELETE CASCADE,
			creation_time NUMERIC NOT NULL,
			update_time NUMERIC NOT NULL,
			content TEXT
		);
		CREATE TABLE text(
			id INTEGER PRIMARY KEY NOT NULL,
			root_id INTEGER NOT NULL REFERENCES node(id) ON DELETE CASCADE,
			creation_time NUMERIC NOT NULL,
			language        CHAR(3),
			title TEXT
		);
		CREATE TABLE node(
			id INTEGER PRIMARY KEY NOT NULL,
			type INTEGER NOT NULL,
			parent_id INTEGER REFERENCES node(id) ON DELETE CASCADE,
			prev_id INTEGER REFERENCES node(id) ON DELETE SET NULL,
			note_id INTEGER REFERENCES note(id) ON DELETE SET NULL,
			text TEXT
		);
		CREATE TABLE lemmatization(
			id INTEGER PRIMARY KEY NOT NULL,
			entry_id INTEGER NOT NULL REFERENCES entry(id) ON DELETE CASCADE,
			lemma_id INTEGER NOT NULL REFERENCES entry(id) ON DELETE CASCADE,
			comment TEXT
		);
		`)
	}
}

func Statistics() map[string]int64 {
	var cnt int64
	var data = make(map[string]int64)
	db.QueryRow("SELECT COUNT(id) FROM text").Scan(&cnt)
	data["TextCount"] = cnt
	db.QueryRow("SELECT COUNT(id) FROM entry").Scan(&cnt)
	data["EntryCount"] = cnt
	db.QueryRow("SELECT COUNT(id) FROM note").Scan(&cnt)
	data["NoteCount"] = cnt
	return data
}

func EntryCount() int64 {
	var cnt int64
	db.QueryRow("SELECT COUNT(id) FROM entry").Scan(&cnt)
	return cnt
}

func TextCount() int64 {
	var cnt int64
	db.QueryRow("SELECT COUNT(id) FROM text").Scan(&cnt)
	return cnt
}

func Entries(offset int64, limit int64, order string) []*Entry {
	entries := make([]*Entry, 0)
	var (
		id                int64
		text              string
		language          string
		creationTimestamp int64
	)
	if len(order) > 0 {
		order = "ORDER BY " + order
	}
	queryString := fmt.Sprintf("SELECT id, text, language, creation_time FROM entry %s LIMIT ? OFFSET ?", order)
	r, _ := db.Query(queryString, limit, offset)
	for r.Next() {
		r.Scan(&id, &text, &language, &creationTimestamp)
		entries = append(entries, &Entry{Id: id, Text: text, Language: language, CreationTimestamp: creationTimestamp})
	}
	return entries
}

func Texts(offset int64, limit int64, order string) []*Text {
	texts := make([]*Text, 0)
	if len(order) > 0 {
		order = "ORDER BY " + order
	}
	queryString := fmt.Sprintf("SELECT id, root_id, creation_time, language, title FROM text %s LIMIT ? OFFSET ?", order)
	r, _ := db.Query(queryString, limit, offset)
	for r.Next() {
		text := new(Text)
		r.Scan(&(text.Id), &(text.RootId), &(text.CreationTimestamp), &(text.Language), &(text.Title))
		texts = append(texts, text)
	}
	return texts
}

func Delete(table string, id int64) {
	db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", table), id)
}

func DeleteNote(id int64) {
	Delete("note", id)
	r, _ := db.Query("SELECT id, type, parent_id, prev_id FROM node WHERE note_id = ?", id)
	nodes := make([]*Node, 0)
	for r.Next() {
		node := new(Node)
		r.Scan(&node.Id, &node.NodeType, &node.ParentId, &node.PrevId)
		node.NoteId = id
		nodes = append(nodes, node)
	}
	for _, node := range nodes {
		node.UnbindNote()
	}
}
