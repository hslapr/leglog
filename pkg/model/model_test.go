package model

import (
	"fmt"
	"testing"
)

func TestModel(t *testing.T) {
	text := NewText("en", "test")
	text.Parse(`The crafting of clear, coherent paragraphs is the subject of considerable stylistic debate. The form varies among different types of writing. For example, newspapers, scientific journals, and fictional essays have somewhat different conventions for the placement of paragraph breaks.`)
	text.Save()
	fmt.Println(text.Id)
	// text := LoadText(db, 5)
	// fmt.Println(text.Root().Children())
	fmt.Println(text.Root().Js())

	// var (
	// 	id       int64
	// 	nodeType int8
	// 	parentId int64
	// 	prevId   int64
	// 	noteId   int64
	// 	text     string
	// )
	// r, _ := db.Query("select id, type, parent_id, prev_id, note_id, text from node where parent_id = 70")
	// for r.Next() {
	// 	r.Scan(&id, &nodeType, &parentId, &prevId, &noteId, &text)
	// 	fmt.Println(id, nodeType, parentId, prevId, noteId, text)
	// }
}
