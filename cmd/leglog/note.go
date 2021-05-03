package main

import (
	"net/http"
	"strconv"

	"github.com/hslapr/leglog/pkg/model"
)

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	model.DeleteNote(id)
}

func addNoteHandler(w http.ResponseWriter, r *http.Request) {
	entryId, _ := strconv.ParseInt(r.FormValue("entryId"), 10, 64)
	content := r.FormValue("content")
	note := model.NewNote(entryId, content)
	note.Save()
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}
