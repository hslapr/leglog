package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/hslapr/leglog/pkg/data"
// 	"github.com/hslapr/leglog/pkg/lemmatizer"
// )

// func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	noteId, _ := strconv.ParseInt(r.FormValue("noteId"), 10, 64)
// 	content := r.FormValue("content")
// 	lemma := r.FormValue("lemma")
// 	language := r.FormValue("language")
// 	var lemmaEntry *data.Entry
// 	if len(lemma) > 0 {
// 		lemmaEntry = data.GetEntry(lemma, language)
// 		if lemmaEntry == nil {
// 			lemmaEntry = data.InsertEntry(lemma, language)
// 		}
// 		data.UpdateNoteContentLemma(noteId, content, lemmaEntry.Id)
// 	} else {
// 		data.UpdateNote(noteId, content)
// 	}
// }

// func getEntriesHandler(w http.ResponseWriter, r *http.Request) {
// 	word := r.FormValue("word")
// 	language := r.FormValue("language")
// 	entries := data.GetEntries(word, language)
// 	js, _ := json.Marshal(entries)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

// func getNotesHandler(w http.ResponseWriter, r *http.Request) {
// 	entryId, _ := strconv.ParseInt(r.FormValue("entry-id"), 10, 64)
// 	notes := data.GetNotes(entryId)
// 	js, _ := json.Marshal(notes)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

// func loadNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	noteId, _ := strconv.ParseInt(r.FormValue("note-id"), 10, 64)
// 	note := data.LoadNote(noteId)
// 	js, _ := json.Marshal(note)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }
// func createNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	var note *data.Note
// 	language := r.FormValue("language")
// 	content := r.FormValue("content")
// 	word := r.FormValue("word")
// 	root := r.FormValue("root")
// 	segId, _ := strconv.ParseInt(r.FormValue("segId"), 10, 64)
// 	entryCase, _ := strconv.ParseInt(r.FormValue("case"), 10, 32)
// 	fmt.Println(language, content, word, root, entryCase)
// 	switch entryCase {
// 	case 0:
// 		word = strings.ToLower(word)
// 	case 1:
// 		word = strings.ToTitle(word)
// 	case 2:
// 		word = strings.ToUpper(word)
// 	}
// 	var rootEntry *data.Entry
// 	entry := data.GetEntry(word, language)
// 	if entry == nil {
// 		entry = data.InsertEntry(word, language)
// 	}
// 	if len(root) > 0 {
// 		rootEntry = data.GetEntry(root, language)
// 		if rootEntry == nil {
// 			rootEntry = data.InsertEntry(root, language)
// 		}
// 	}
// 	if rootEntry != nil {
// 		note = data.InsertNote(entry.Id, content, rootEntry.Id)
// 	} else {
// 		note = data.InsertNote(entry.Id, content, -1)
// 	}
// 	data.SetSegmentNote(segId, note.Id)
// 	// fmt.Println("segnoteid", segId, note.Id)
// 	js, _ := json.Marshal(note)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

// func groupSegmentsHandler(w http.ResponseWriter, r *http.Request) {
// 	startId, _ := strconv.ParseInt(r.FormValue("startId"), 10, 64)
// 	endId, _ := strconv.ParseInt(r.FormValue("endId"), 10, 64)
// 	id := data.GroupSegments(startId, endId)
// 	var o map[string]int64 = make(map[string]int64)
// 	o["seg-id"] = id
// 	js, _ := json.Marshal(o)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

// func setSegmentNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	segId, _ := strconv.ParseInt(r.FormValue("segId"), 10, 64)
// 	noteId, _ := strconv.ParseInt(r.FormValue("noteId"), 10, 64)
// 	data.SetSegmentNote(segId, noteId)
// 	var o map[string]int64 = make(map[string]int64)
// 	o["note-id"] = noteId
// 	js, _ := json.Marshal(o)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }
