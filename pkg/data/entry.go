package data

import (
	"database/sql"
	"time"
)

type Entry struct {
	Id           int64
	Name         string
	notes        []*Note `json:"-"`
	CreationTime time.Time
	Language     Language
}

type Note struct {
	Id           int64
	Entry        *Entry
	Content      string
	Root         *Entry
	CreationTime time.Time `json:"-"`
	UpdateTime   time.Time `json:"-"`
}

func (entry *Entry) Notes() (notes []*Note) {
	var id int64
	r, _ := db.Query("select id from note where entry_id = ?", entry.Id)
	for r.Next() {
		r.Scan(&id)
		notes = append(notes, LoadNote(id))
	}
	return notes
}
func (entry *Entry) Derives() (derives []*Entry) {
	var id int64
	r, _ := db.Query("select entry_id from note where root_id = ?", entry.Id)
	for r.Next() {
		r.Scan(&id)
		derives = append(derives, LoadEntry(id))
	}
	return derives
}

func GetEntries(word string, language string) (entries []*Entry) {
	var id int64
	var name string
	var unixSec int64
	r, _ := db.Query("select id,name,CREATION_TIME from ENTRY where name like ? and language = ?", word, language)
	for r.Next() {
		r.Scan(&id, &name, &unixSec)
		entries = append(entries, &Entry{Id: id, Name: name, Language: Language(language), CreationTime: time.Unix(unixSec, 0)})
	}
	return entries
}

func GetEntry(word string, language string) (entry *Entry) {
	var id int64
	var name string
	var unixSec int64
	err := db.QueryRow("select id,name,CREATION_TIME from ENTRY where name = ? and language = ?", word, language).Scan(&id, &name, &unixSec)
	if err != nil {
		return nil
	}
	return &Entry{Id: id, Name: name, Language: Language(language), CreationTime: time.Unix(unixSec, 0)}
}

func GetNotes(entryId int64) (notes []*Note) {
	var id int64
	var rootId sql.NullInt64
	var content string
	var note *Note
	var unixSecCreation int64
	var unixSecUpdate int64
	r, _ := db.Query("select id, root_id, content,creation_time,update_time from NOTE where ENTRY_ID = ?", entryId)
	for r.Next() {
		note = new(Note)
		note.Entry = LoadEntry(entryId)
		r.Scan(&id, &rootId, &content, &unixSecCreation, &unixSecUpdate)
		note.Id = id
		note.Content = content
		note.CreationTime = time.Unix(unixSecCreation, 0)
		note.UpdateTime = time.Unix(unixSecUpdate, 0)
		if rootId.Valid {
			note.Root = LoadEntry(rootId.Int64)
		}
		notes = append(notes, note)
	}
	return notes
}

func LoadNote(id int64) (note *Note) {
	note = new(Note)
	var (
		entry_id        int64
		root_id         sql.NullInt64
		content         string
		unixSecCreation int64
		unixSecUpdate   int64
	)
	db.QueryRow("select entry_id, root_id, content,creation_time,update_time from note where id = ?", id).Scan(&entry_id, &root_id, &content, &unixSecCreation, &unixSecUpdate)
	note.Id = id
	note.Entry = LoadEntry(entry_id)
	if root_id.Valid {
		note.Root = LoadEntry(root_id.Int64)
	}
	note.Content = content
	note.CreationTime = time.Unix(unixSecCreation, 0)
	note.UpdateTime = time.Unix(unixSecUpdate, 0)
	return note
}

func LoadEntry(id int64) (entry *Entry) {
	var (
		name     string
		language string
		sec      int64
	)
	db.QueryRow("select name, language, creation_time from entry where id = ?", id).Scan(&name, &language, &sec)
	return &Entry{Id: id, Name: name, Language: Language(language), CreationTime: time.Unix(sec, 0)}
}

type Language string

const (
	ENGLISH Language = "en"
	ITALIAN Language = "it"
	GERMAN  Language = "de"
	FRENCH  Language = "fr"
)

func UpdateNote(id int64, content string) {
	stmt, _ := db.Prepare("update note set CONTENT = ?, UPDATE_TIME=? where ID = ?")
	stmt.Exec(content, time.Now().Unix(), id)
}

func UpdateNoteContentLemma(id int64, content string, lemmaId int64) {
	stmt, _ := db.Prepare("update note set CONTENT = ?, UPDATE_TIME=?, ROOT_ID=? where ID = ?")
	stmt.Exec(content, time.Now().Unix(), lemmaId, id)
}
