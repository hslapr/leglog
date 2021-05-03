package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/hslapr/leglog/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var dbPath string = config.Config.DatabasePath
	db, _ = sql.Open("sqlite3", dbPath)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		os.Create(dbPath)
		_, err = db.Exec(`
		CREATE TABLE entry(
			id        INTEGER PRIMARY KEY     NOT NULL,
			text      TEXT    NOT NULL UNIQUE,
			language        CHAR(8),
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
			parent_id INTEGER REFERENCES node(id) ON DELETE SET NULL,
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
		if err != nil {
			log.Printf("model.init: %s", err)
		}
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
	db.QueryRow("SELECT COUNT(id) FROM node").Scan(&cnt)
	data["NodeCount"] = cnt
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

func Entries(offset int64, limit int64, lang string, order string) []*Entry {
	var queryString string
	var r *sql.Rows
	var e error
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
	if len(lang) > 0 {
		queryString = fmt.Sprintf("SELECT id, text, language, creation_time FROM entry WHERE language = ? %s LIMIT ? OFFSET ?", order)
		r, e = db.Query(queryString, lang, limit, offset)
		if e != nil {
			log.Printf("model.Entries: %s", e)
		}
	} else {
		queryString = fmt.Sprintf("SELECT id, text, language, creation_time FROM entry %s LIMIT ? OFFSET ?", order)
		r, e = db.Query(queryString, limit, offset)
		if e != nil {
			log.Printf("model.Entries: %s", e)
		}
	}
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
	r, e := db.Query(queryString, limit, offset)
	if e != nil {
		log.Printf("model.Texts: %s", e)
	}
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
	Delete("note", id)
}

func DeleteEntry(id int64) {
	r, _ := db.Query("SELECT id FROM note WHERE entry_id = ?", id)
	noteIds := make([]int64, 0)
	var noteId int64
	for r.Next() {
		r.Scan(&noteId)
		noteIds = append(noteIds, noteId)
	}
	for _, noteId = range noteIds {
		DeleteNote(noteId)
	}
	db.Exec("DELETE FROM lemmatization WHERE entry_id = ? OR lemma_id = ?", id, id)
	Delete("entry", id)
	log.Printf("model.DeleteEntry: id = %d", id)
}

func DeleteText(id int64, rootId int64) {
	log.Printf("model.DeleteText: id = %d, rootId = %d", id, rootId)
	DeleteBranch(rootId)
	Delete("text", id)
}

func DeleteBranch(rootId int64) {
	db.Exec("UPDATE node SET prev_id = (SELECT prev_id FROM node WHERE id = ?) WHERE prev_id = ?", rootId, rootId)
	var stack []int64
	stack = make([]int64, 0)
	stack = append(stack, rootId)
	var nodeId int64
	var childId int64
	for len(stack) > 0 {
		nodeId = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		r, e := db.Query("SELECT id FROM node WHERE parent_id = ?", nodeId)
		if e != nil {
			log.Printf("model.DeleteBranch: %s", e)
		}
		for r.Next() {
			r.Scan(&childId)
			stack = append(stack, childId)
		}
		Delete("node", nodeId)
	}
}

func RemoveLemma(entryId int64, lemmaId int64) {
	_, e := db.Exec("DELETE FROM lemmatization WHERE entry_id = ? AND lemma_id = ?", entryId, lemmaId)
	if e != nil {
		log.Printf("model.RemoveLemma: %s", e)
	}
}
