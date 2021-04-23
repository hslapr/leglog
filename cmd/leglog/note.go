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
