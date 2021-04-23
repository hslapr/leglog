package data

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const DB_PATH string = "../../assets/leglog.db"

func init() {
	db, _ = sql.Open("sqlite3", DB_PATH)
	if _, err := os.Stat(DB_PATH); os.IsNotExist(err) {
		os.Create(DB_PATH)
		db.Exec(`
		CREATE TABLE entry(
			ID        INTEGER PRIMARY KEY     NOT NULL,
			NAME      TEXT    NOT NULL,
			LANGUAGE        CHAR(3),
			CREATION_TIME NUMERIC NOT NULL
	 	);
		CREATE TABLE note(
			ID INTEGER PRIMARY KEY NOT NULL,
			ENTRY_ID INTEGER NOT NULL REFERENCES entry(ID) ON DELETE CASCADE,
			ROOT_ID INTEGER REFERENCES entry(ID) ON DELETE SET NULL,
			CONTENT TEXT,
			CREATION_TIME NUMERIC NOT NULL,
			UPDATE_TIME NUMERIC NOT NULL
		);
		CREATE TABLE text(
			ID INTEGER PRIMARY KEY NOT NULL,
			HEAD_ID INTEGER NOT NULL  REFERENCES segment(ID) ON DELETE CASCADE,
			LANGUAGE        CHAR(3),
			CREATION_TIME NUMERIC NOT NULL
		);
		CREATE TABLE segment(
			ID INTEGER PRIMARY KEY NOT NULL,
			NEXT_ID INTEGER REFERENCES segment(ID) ON DELETE SET NULL,
			CHILD_ID INTEGER REFERENCES segment(ID) ON DELETE SET NULL,
			TEXT TEXT,
			NOTE_ID INTEGER REFERENCES note(ID) ON DELETE SET NULL,
			IS_WORD NUMERIC
		);
		`)
	}
	// InsertNote(1, "wow", -1)
}

func InsertEntry(name string, language string) *Entry {
	stmt, _ := db.Prepare("insert into entry (NAME, LANGUAGE, CREATION_TIME) values (?, ?, ?)")
	now := time.Now()
	r, _ := stmt.Exec(name, language, now.Unix())
	entry := new(Entry)
	entry.Id, _ = r.LastInsertId()
	entry.Name = name
	entry.Language = Language(language)
	entry.CreationTime = now
	return entry
}

func InsertNote(entryId int64, content string, rootId int64) *Note {
	var r sql.Result
	note := new(Note)
	note.Content = content
	now := time.Now()
	note.Entry = LoadEntry(entryId)
	if rootId >= 0 {
		stmt, _ := db.Prepare("insert into note (ENTRY_ID, CONTENT, ROOT_ID, CREATION_TIME, UPDATE_TIME) values (?, ?, ?, ?, ?)")
		r, _ = stmt.Exec(entryId, content, rootId, now.Unix(), now.Unix())
		note.Root = LoadEntry(rootId)
	} else {
		stmt, _ := db.Prepare("insert into note (ENTRY_ID, CONTENT, CREATION_TIME, UPDATE_TIME) values (?, ?, ?, ?)")
		r, _ = stmt.Exec(entryId, content, now.Unix(), now.Unix())
	}
	note.Id, _ = r.LastInsertId()
	note.CreationTime = now
	note.UpdateTime = now
	return note
}

func SetSegmentNote(segId int64, noteId int64) {
	fmt.Println(segId, noteId)
	stmt, _ := db.Prepare("update segment set NOTE_ID = ? where ID = ?")
	stmt.Exec(noteId, segId)
}

func Entries() (entries []*Entry) {
	var (
		id       int64
		name     string
		language string
		unixTime int64
	)
	r, _ := db.Query("select id, name, language, creation_time from entry")
	for r.Next() {
		r.Scan(&id, &name, &language, &unixTime)
		entries = append(entries, &Entry{Id: id, Name: name, Language: Language(language), CreationTime: time.Unix(unixTime, 0)})
	}
	return entries
}

func Texts() (texts []*Text) {
	var (
		id       int64
		headId   int64
		language string
		unixTime int64
		t        *Text
	)
	r, _ := db.Query("SELECT ID, HEAD_ID, LANGUAGE, CREATION_TIME FROM text")
	for r.Next() {
		r.Scan(&id, &headId, &language, &unixTime)
		t = new(Text)
		t.Id = id
		t.Head = LoadSegment(headId)
		t.Language = Language(language)
		t.CreationTime = time.Unix(unixTime, 0)
		texts = append(texts, t)
	}
	return texts
}
