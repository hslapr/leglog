package data

import (
	"database/sql"
	"fmt"
	"html/template"
	"strings"
	"time"
)

type Text struct {
	Id           int64
	Head         *Segment
	Tail         *Segment
	Language     Language
	CreationTime time.Time
}

func (t *Text) Append(seg *Segment) {
	if t.Head == nil {
		t.Head = seg
		t.Tail = seg
	} else {
		t.Tail.Next = seg
		t.Tail = seg
	}
}

func (seg *Segment) SetNext(s *Segment) {
	seg.Next = s
	s.Prev = seg
}

func LoadText(id int64) (t *Text) {
	var (
		headId   int64
		language string
		unixSec  int64
	)
	t = &Text{Id: id}
	stmt, _ := db.Prepare("SELECT HEAD_ID, LANGUAGE, CREATION_TIME FROM text WHERE ID=?")
	stmt.QueryRow(id).Scan(&headId, &language, &unixSec)
	t.Head = LoadSegment(headId)
	t.Language = Language(language)
	t.CreationTime = time.Unix(unixSec, 0)
	return t
}

func LoadSegment(id int64) (seg *Segment) {
	// stmt, _ := db.Prepare("SELECT NEXT_ID, CHILD_ID, TEXT, NOTE_ID, IS_WORD FROM segment WHERE ID=?")
	var next_id sql.NullInt64
	var child_id sql.NullInt64
	var text string
	var noteId sql.NullInt64
	var is_word bool
	db.QueryRow("SELECT NEXT_ID, CHILD_ID, NOTE_ID, IS_WORD, TEXT FROM segment WHERE ID=?", id).Scan(&next_id, &child_id, &noteId, &is_word, &text)
	// fmt.Println("load:", id, next_id, child_id, text, noteId, is_word)
	seg = NewSegment(text, is_word)
	if next_id.Valid {
		seg.SetNext(LoadSegment(next_id.Int64))
	}
	if child_id.Valid {
		seg.Child = LoadSegment(child_id.Int64)
	}
	seg.Id = id
	if noteId.Valid {
		seg.NoteId = noteId.Int64
	} else {
		seg.NoteId = -1
	}
	return seg
}

func (t *Text) Save() {
	seg := t.Head
	for seg != nil {
		seg.Save()
		child_seg := seg.Child
		for child_seg != nil {
			child_seg.Save()
			child_seg = child_seg.Next
		}
		seg = seg.Next
	}
	stmt, _ := db.Prepare("INSERT INTO text(HEAD_ID, LANGUAGE, CREATION_TIME) VALUES(?,?, ?)")
	now := time.Now()
	r, _ := stmt.Exec(t.Head.Id, t.Language, now.Unix())
	t.CreationTime = now
	t.Id, _ = r.LastInsertId()
	seg = t.Head
	for seg != nil {
		seg.UpdateNextId()
		seg.UpdateChild()
		child_seg := seg.Child
		for child_seg != nil {
			child_seg.UpdateNextId()
			child_seg = child_seg.Next
		}
		seg = seg.Next
	}
}

type Segment struct {
	Id     int64
	Prev   *Segment
	Next   *Segment
	Child  *Segment
	Text   string
	note   *Note
	NoteId int64
	IsWord bool
}

func (seg *Segment) Note() *Note {
	if seg.note == nil {
		var noteId sql.NullInt64
		db.QueryRow("SELECT NOTE_ID FROM segment WHERE id = ?", seg.Id).Scan(&noteId)
		if noteId.Valid {
			seg.note = LoadNote(noteId.Int64)
		}
	}
	return seg.note
}

func (seg *Segment) Save() {
	stmt, _ := db.Prepare("INSERT INTO segment(TEXT, IS_WORD) VALUES (?,?)")
	r, _ := stmt.Exec(seg.Text, seg.IsWord)
	seg.Id, _ = r.LastInsertId()
}

func (seg *Segment) UpdateNextId() {
	if seg.Next != nil {
		stmt, _ := db.Prepare("UPDATE segment SET NEXT_ID = ? WHERE ID = ?")
		stmt.Exec(seg.Next.Id, seg.Id)
	}
}

func (seg *Segment) UpdateChild() {
	if seg.Child != nil {
		stmt, _ := db.Prepare("UPDATE segment SET CHILD_ID = ? WHERE ID = ?")
		stmt.Exec(seg.Child.Id, seg.Id)
	}
}

func (t *Text) InsertBefore(seg *Segment, newSeg *Segment) {
	if seg.Prev == nil {
		t.Head = newSeg
		newSeg.Next = seg
		seg.Prev = newSeg
	} else {
		seg.Prev.Next = newSeg
		newSeg.Prev = seg.Prev
		newSeg.Next = seg
		seg.Prev = newSeg
	}
}

func (sentence *Text) Remove(seg *Segment) {
	if seg.Prev == nil {
		sentence.Head = seg.Next
	} else {
		seg.Prev.Next = seg.Next
	}
	if seg.Next == nil {
		sentence.Tail = seg.Prev
	} else {
		seg.Next.Prev = seg.Prev
	}
}

func (sentence *Text) RemoveBetween(seg1 *Segment, seg2 *Segment) {
	s := seg1
	for {
		sentence.Remove(s)
		if s == seg2 {
			break
		}
		s = s.Next
	}
}

func NewSegment(t string, isWord bool) *Segment {
	return &Segment{Text: t, IsWord: isWord, NoteId: -1}
}

func (t *Text) Excerpt(max int) string {
	var builder strings.Builder
	seg := t.Head
	for seg != nil {
		if seg.Child != nil {
			s := seg.Child
			for s != nil {
				builder.WriteString(s.Text)
				s = s.Next
			}
		} else {
			builder.WriteString(seg.Text)
		}
		if builder.Len() >= max {
			builder.WriteString("...")
			break
		}
		seg = seg.Next
	}
	return builder.String()
}

func (t *Text) Html() template.HTML {
	var builder strings.Builder
	seg := t.Head
	for seg != nil {
		if seg.Text == "\n" {
			builder.WriteString("<br />")
		} else if seg.Child != nil {
			if seg.NoteId > -1 {
				builder.WriteString(fmt.Sprintf("<span class=\"phrase\" seg-id=\"%d\" note-id=\"%d\">", seg.Id, seg.NoteId))
			} else {
				builder.WriteString(fmt.Sprintf("<span class=\"phrase\" seg-id=\"%d\">", seg.Id))
			}
			s := seg.Child
			for s != nil {
				if s.IsWord {
					if s.NoteId > -1 {
						builder.WriteString(fmt.Sprintf("<span class=\"word\" seg-id=\"%d\" note-id=\"%d\">%s</span>", s.Id, s.NoteId, s.Text))
					} else {
						builder.WriteString(fmt.Sprintf("<span class=\"word\" seg-id=\"%d\">%s</span>", s.Id, s.Text))
					}
				} else {
					builder.WriteString(s.Text)
				}
				s = s.Next
			}
			builder.WriteString("</span>")
		} else if seg.IsWord {
			if seg.NoteId > -1 {
				builder.WriteString(fmt.Sprintf("<span class=\"word\" seg-id=\"%d\" note-id=\"%d\">%s</span>", seg.Id, seg.NoteId, seg.Text))
			} else {
				builder.WriteString(fmt.Sprintf("<span class=\"word\" seg-id=\"%d\">%s</span>", seg.Id, seg.Text))
			}
		} else {
			builder.WriteString(seg.Text)
		}
		seg = seg.Next
	}
	return template.HTML(builder.String())
}

func (sentence *Text) Group(first *Segment, last *Segment) {
	g := NewSegment("", true)
	g.Child = first
	sentence.InsertBefore(first, g)
	sentence.RemoveBetween(first, last)
	last.Next = nil
}

func GroupSegments(startId int64, endId int64) (id int64) {
	var nextId sql.NullInt64
	var r sql.Result
	db.QueryRow("select next_id from segment where id = ?", endId).Scan(&nextId)
	if nextId.Valid {
		stmt, _ := db.Prepare("insert into segment(next_id, child_id, is_word) values (?, ?, ?)")
		r, _ = stmt.Exec(nextId.Int64, startId, true)
	} else {
		stmt, _ := db.Prepare("insert into segment(child_id, is_word) values (?, ?)")
		r, _ = stmt.Exec(startId, true)
	}
	id, _ = r.LastInsertId()
	stmt, _ := db.Prepare("update segment set next_id = ? where next_id = ?")
	stmt.Exec(id, startId)
	stmt, _ = db.Prepare("update segment set next_id = null where id = ?")
	stmt.Exec(endId)
	return id
}

func CountEntry() int {
	var count int
	db.QueryRow("select count(*) from entry").Scan(&count)
	return count
}

func CountText() int {
	var count int
	db.QueryRow("select count(*) from text").Scan(&count)
	return count
}

func CountNote() int {
	var count int
	db.QueryRow("select count(*) from note").Scan(&count)
	return count
}

func CountSegment() int {
	var count int
	db.QueryRow("select count(*) from segment").Scan(&count)
	return count
}
