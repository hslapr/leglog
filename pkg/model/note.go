package model

import (
	"database/sql"
	"log"
	"time"
)

type Note struct {
	Id                int64
	EntryId           int64
	Content           string
	CreationTimestamp int64 // Unix timestamp
	UpdateTimestamp   int64
	EntryText         string // no column
}

func NewNote(entryId int64, content string) *Note {
	return &Note{EntryId: entryId, Content: content}
}

func (note *Note) Save() {
	note.save(db)
}

func (note *Note) save(db *sql.DB) {
	if note.Id < 1 {
		note.CreationTimestamp = time.Now().Unix()
		note.UpdateTimestamp = note.CreationTimestamp
		r, _ := db.Exec("INSERT INTO note (entry_id, content, creation_time, update_time) VALUES (?, ?, ?, ?)",
			note.EntryId, note.Content, note.CreationTimestamp, note.UpdateTimestamp)
		note.Id, _ = r.LastInsertId()
	}
}

func (note *Note) SaveContent() {
	note.saveContent(db)
}

func (note *Note) saveContent(db *sql.DB) {
	if note.Id > 0 {
		note.UpdateTimestamp = time.Now().Unix()
		_, e := db.Exec("UPDATE note SET content = ?, update_time = ? WHERE id = ?", note.Content, note.UpdateTimestamp, note.Id)
		if e != nil {
			log.Println(e)
		}
	}
}

func (note *Note) Load() {
	note.load(db)
}

func (note *Note) load(db *sql.DB) {
	if note.Id > 0 {
		db.QueryRow(`
		SELECT note.entry_id, note.creation_time, note.update_time, note.content, entry.text
		FROM note INNER JOIN entry ON note.entry_id = entry.id WHERE note.id = ?
		`, note.Id).Scan(&(note.EntryId), &note.CreationTimestamp, &note.UpdateTimestamp, &note.Content, &note.EntryText)
	}
}

func QueryNotes(s string, language string) (notes []*Note) {
	var (
		id                int64
		entryId           int64
		creationTimestamp int64
		updateTimestamp   int64
		content           string
		entryText         string
	)
	r, e := db.Query(`
	WITH ids AS (SELECT id FROM entry WHERE text LIKE ? AND language = ?)
	SELECT note.id, note.entry_id, note.creation_time, note.update_time, note.content, entry.text FROM note
	INNER JOIN entry ON note.entry_id = entry.id
	WHERE note.entry_id IN (select id from ids UNION SELECT lemma_id AS id FROM lemmatization WHERE entry_id IN ids)
	`, s, language)
	if e != nil {
		log.Println(e)
	}
	for r.Next() {
		r.Scan(&id, &entryId, &creationTimestamp, &updateTimestamp, &content, &entryText)
		notes = append(notes, &Note{Id: id, EntryId: entryId, CreationTimestamp: creationTimestamp, UpdateTimestamp: updateTimestamp, Content: content, EntryText: entryText})
	}
	return notes
}
