package model

import (
	"database/sql"
	"time"
)

type Entry struct {
	Id                int64
	Text              string
	Language          string
	CreationTimestamp int64
}

type Derive struct {
	Entry
	Comment string
}

func (entry *Entry) CreationTime(layout string) string {
	return time.Unix(entry.CreationTimestamp, 0).Format(layout)
}

func (entry *Entry) AddLemma(lemma *Entry, comment string) {
	db.Exec("INSERT INTO lemmatization (entry_id, lemma_id, comment) VALUES (?, ?, ?)", entry.Id, lemma.Id, comment)
}

func NewEntry(text string, language string) *Entry {
	return &Entry{Text: text, Language: language}
}

func GetEntries(s string, language string) []*Entry {
	entries := make([]*Entry, 0)
	var id int64
	var text string
	var timestamp int64
	r, _ := db.Query("select id,text,CREATION_TIME from ENTRY where text like ? and language = ?", s, language)
	for r.Next() {
		r.Scan(&id, &text, &timestamp)
		entries = append(entries, &Entry{Id: id, Text: text, Language: language, CreationTimestamp: timestamp})
	}
	return entries
}

func GetEntry(s string, language string) *Entry {
	var id int64
	var text string
	var timestamp int64
	err := db.QueryRow("select id,text,CREATION_TIME from ENTRY where text = ? and language = ?", s, language).
		Scan(&id, &text, &timestamp)
	if err != nil {
		return nil
	}
	return &Entry{Id: id, Text: text, Language: language, CreationTimestamp: timestamp}
}

func (entry *Entry) Save() {
	entry.save(db)
}

func (entry *Entry) save(db *sql.DB) {
	if entry.Id < 1 {
		entry.CreationTimestamp = time.Now().Unix()
		r, _ := db.Exec("INSERT INTO entry (text, language, creation_time) VALUES (?, ?, ?)",
			entry.Text, entry.Language, entry.CreationTimestamp)
		entry.Id, _ = r.LastInsertId()
	}
}

func (entry *Entry) Load() {
	entry.load(db)
}

func (entry *Entry) load(db *sql.DB) {
	if entry.Id > 0 {
		db.QueryRow("SELECT text, language, creation_time FROM entry WHERE id = ?", entry.Id).Scan(
			&entry.Text, &entry.Language, &entry.CreationTimestamp,
		)
	}
}

func (entry *Entry) Notes() []*Note {
	return entry.notes(db)
}

func (entry *Entry) notes(db *sql.DB) []*Note {
	var notes = make([]*Note, 0)
	r, _ := db.Query("SELECT id, content, creation_time, update_time FROM note WHERE entry_id = ?", entry.Id)
	for r.Next() {
		note := new(Note)
		note.EntryId = entry.Id
		note.EntryText = entry.Text
		r.Scan(&note.Id, &note.Content, &note.CreationTimestamp, &note.UpdateTimestamp)
		notes = append(notes, note)
	}
	return notes
}

func (entry *Entry) Lemmas() []*Entry {
	return entry.lemmas(db)
}

func (entry *Entry) lemmas(db *sql.DB) []*Entry {
	var entries = make([]*Entry, 0)
	r, _ := db.Query("SELECT id, text, language, creation_time FROM entry WHERE id IN (SELECT lemma_id FROM lemmatization WHERE entry_id = ?)", entry.Id)
	for r.Next() {
		e := new(Entry)
		r.Scan(&e.Id, &e.Text, &e.Language, &e.CreationTimestamp)
		entries = append(entries, e)
	}
	return entries
}

func (entry *Entry) Derives() []*Derive {
	return entry.derives(db)
}

func (entry *Entry) derives(db *sql.DB) []*Derive {
	var derives = make([]*Derive, 0)
	r, _ := db.Query("SELECT entry.id, entry.text, entry.language, entry.creation_time, lemmatization.comment FROM entry INNER JOIN lemmatization ON entry.id = lemmatization.entry_id WHERE lemmatization.lemma_id = ?", entry.Id)
	for r.Next() {
		d := new(Derive)
		r.Scan(&d.Id, &d.Text, &d.Language, &d.CreationTimestamp, &d.Comment)
		derives = append(derives, d)
	}
	return derives
}
