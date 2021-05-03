package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/hslapr/leglog/pkg/lemmatizer"
	"github.com/hslapr/leglog/pkg/model"
	"github.com/hslapr/leglog/pkg/util"
)

type CreateNotePostData struct {
	Nodes     []*model.Node
	Content   string
	EntryText string
	Language  string
	Lemma     string
	Comment   string
}

type BindNotePostData struct {
	Nodes []*model.Node
	Note  *model.Note
}

func bindNoteHandler(w http.ResponseWriter, r *http.Request) {
	var rdata = make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	var data BindNotePostData
	decoder.Decode(&data)
	if len(data.Nodes) > 1 {
		phraseNode, _ := model.CreatePhrase(data.Nodes, data.Note)
		rdata["phraseNode"] = phraseNode
	} else {
		data.Nodes[0].NoteId = data.Note.Id
		data.Nodes[0].SaveNoteId()
	}
	rdata["note"] = data.Note
	js, _ := json.Marshal(rdata)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func unbindNoteHandler(w http.ResponseWriter, r *http.Request) {
	var node = new(model.Node)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&node)
	node.UnbindNote()
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note *model.Note
	decoder := json.NewDecoder(r.Body)
	var data CreateNotePostData
	decoder.Decode(&data)
	entry := model.GetEntry(data.EntryText, data.Language)
	if entry == nil {
		entry = model.NewEntry(data.EntryText, data.Language)
		entry.Save()
	}
	if len(data.Lemma) > 0 {
		lemma := model.GetEntry(data.Lemma, data.Language)
		if lemma == nil {
			lemma = model.NewEntry(data.Lemma, data.Language)
			lemma.Save()
		}
		note = model.NewNote(lemma.Id, data.Content)
		note.Save()
		note.EntryText = lemma.Text
		entry.AddLemma(lemma, data.Comment)
	} else {
		note = model.NewNote(entry.Id, data.Content)
		note.Save()
		note.EntryText = entry.Text
	}
	var rdata = make(map[string]interface{})
	// var rdata = make(map[string]int64)
	if len(data.Nodes) > 1 {
		phraseNode, _ := model.CreatePhrase(data.Nodes, note)
		rdata["phraseNode"] = phraseNode
	} else {
		data.Nodes[0].NoteId = note.Id
		data.Nodes[0].SaveNoteId()
	}
	// rdata["noteId"] = note.Id
	rdata["note"] = note
	js, _ := json.Marshal(rdata)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func loadNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteId, _ := strconv.ParseInt(r.FormValue("noteId"), 10, 64)
	note := &model.Note{Id: noteId}
	note.Load()
	js, _ := json.Marshal(note)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note = new(model.Note)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(note)
	note.SaveContent()
	js, _ := json.Marshal(note)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func queryNotesHandler(w http.ResponseWriter, r *http.Request) {
	language := r.FormValue("language")
	text := r.FormValue("text")
	notes := model.QueryNotes(text, language)
	js, _ := json.Marshal(notes)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getLemmasHandler(w http.ResponseWriter, r *http.Request) {
	language := r.FormValue("language")
	word := r.FormValue("word")
	lemmas := lemmatizer.Lemmatizers[language].Lemmatize(word)
	js, _ := json.Marshal(lemmas)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func removeLemmaHandler(w http.ResponseWriter, r *http.Request) {
	entryId, _ := strconv.ParseInt(r.FormValue("entryId"), 10, 64)
	lemmaId, _ := strconv.ParseInt(r.FormValue("lemmaId"), 10, 64)
	log.Printf("removeLemmaHandler: entryId = %d, lemmaId = %d", entryId, lemmaId)
	model.RemoveLemma(entryId, lemmaId)
}

func deleteEntryHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	model.DeleteEntry(id)
	http.Redirect(w, r, "/entry/", http.StatusFound)
}

func deleteTextHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("textId"), 10, 64)
	rootId, _ := strconv.ParseInt(r.FormValue("rootId"), 10, 64)
	model.DeleteText(id, rootId)
	http.Redirect(w, r, "/text/", http.StatusFound)
}

func addParagraphHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("textId"), 10, 64)
	rootId, _ := strconv.ParseInt(r.FormValue("rootId"), 10, 64)
	s := r.FormValue("s")
	s = util.SanitizeText(s)
	text := &model.Text{Id: id, RootId: rootId}
	text.ParseAppendParagraphs(s)
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func deleteParagraphHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("paragraphId"), 10, 64)
	model.DeleteBranch(id)
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func changeTitleHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	title := r.FormValue("title")
	text := &model.Text{Id: id, Title: title}
	text.SaveTitle()
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}
