package model

import (
	"database/sql"
	"fmt"
	"strings"
)

// NodeType
const (
	ROOT      = 0
	PARAGRAPH = 1
	PHRASE    = 2
	WORD      = 3
	NUMBER    = 4
	SPACE     = 5
	OTHER     = 6
)

type Node struct {
	Id       int64
	NodeType int8
	ParentId int64
	PrevId   int64
	NoteId   int64
	Text     string
}

func NewNode(text string, nodeType int8) *Node {
	return &Node{Text: text, NodeType: nodeType}
}

func (node *Node) Save() {
	node.save(db)
}

func (node *Node) save(db *sql.DB) {
	// update node.Id
	if node.Id < 1 {
		r, _ := db.Exec("INSERT INTO node (type, parent_id, prev_id, note_id, text) VALUES (?, ?, ?, ?, ?)",
			node.NodeType, node.ParentId, node.PrevId, node.NoteId, node.Text)
		node.Id, _ = r.LastInsertId()
	} else {
		db.Exec("UPDATE node SET type = ?, parent_id = ?, prev_id = ?, note_id = ?, text = ? WHERE id = ?",
			node.NodeType, node.ParentId, node.PrevId, node.NoteId, node.Text, node.Id)
	}
}

func (node *Node) SaveNoteId() {
	node.saveNoteId(db)
}

func (node *Node) saveNoteId(db *sql.DB) {
	if node.Id > 0 {
		db.Exec("UPDATE node SET note_id = ? WHERE id = ?", node.NoteId, node.Id)
	}
}

func (node *Node) UnbindNote() {
	node.unbindNote(db)
}

func (node *Node) unbindNote(db *sql.DB) {
	if node.NodeType == PHRASE {
		db.Exec("UPDATE node SET prev_id = ? WHERE parent_id = ? AND (prev_id IS NULL OR prev_id = 0)", node.PrevId, node.Id)
		db.Exec("UPDATE node SET prev_id = (SELECT id FROM node WHERE parent_id = ? AND id NOT IN (SELECT prev_id FROM node)) WHERE prev_id = ?", node.Id, node.Id)
		db.Exec("UPDATE node SET parent_id = ? WHERE parent_id = ?", node.ParentId, node.Id)
		db.Exec("DELETE FROM node WHERE id = ?", node.Id)
	} else {
		node.NoteId = 0
		node.SaveNoteId()
	}
}

func (node *Node) Children() (children []*Node) {
	return node.children(db)
}

func (node *Node) children(db *sql.DB) (children []*Node) {
	children = make([]*Node, 0)
	var (
		id       int64
		nodeType int8
		parentId int64
		prevId   int64
		noteId   int64
		text     string
	)
	for {
		err := db.QueryRow("select id, type, parent_id, prev_id, note_id, text from node where parent_id = ? and prev_id = ?", node.Id, id).Scan(&id, &nodeType, &parentId, &prevId, &noteId, &text)
		if err != nil {
			break
		}
		children = append(children,
			&Node{Id: id, NodeType: nodeType, ParentId: parentId, PrevId: prevId, NoteId: noteId, Text: text})
	}
	return children
}

func (node *Node) Prev() *Node {
	return node.prev(db)
}

func (node *Node) prev(db *sql.DB) *Node {
	if node.PrevId < 1 {
		return nil
	}
	return LoadNode(db, node.PrevId)
}

func (node *Node) Js() string {
	var builder strings.Builder
	// builder.WriteString(fmt.Sprintf("let node%d = new Node(%d, %d,%d,%d, %d, \"%s\");\n",
	// 	node.Id, node.Id, node.NodeType, node.ParentId, node.PrevId, node.NoteId, node.Text))
	builder.WriteString(fmt.Sprintf("window.Leglog.nodes[%d] = new Node(%d, %d,%d,%d, %d, \"%s\");\n",
		node.Id, node.Id, node.NodeType, node.ParentId, node.PrevId, node.NoteId, node.Text))
	children := node.Children()
	if len(children) > 0 {
		var lastChildId int64
		for _, child := range children {
			builder.WriteString(child.Js())
			// builder.WriteString(fmt.Sprintf("node%d.children.push(node%d);\n",
			// 	node.Id, child.Id))
			// builder.WriteString(fmt.Sprintf("node%d.parent = node%d;\n",
			// 	child.Id, node.Id))
			builder.WriteString(fmt.Sprintf("window.Leglog.nodes[%d].children.push(window.Leglog.nodes[%d]);\n",
				node.Id, child.Id))
			builder.WriteString(fmt.Sprintf("window.Leglog.nodes[%d].parent = window.Leglog.nodes[%d];\n",
				child.Id, node.Id))
			if lastChildId > 0 {
				// builder.WriteString(fmt.Sprintf("node%d.next = node%d;\n",
				// 	lastChildId, child.Id))
				// builder.WriteString(fmt.Sprintf("node%d.prev = node%d;\n",
				// 	child.Id, lastChildId))
				builder.WriteString(fmt.Sprintf("window.Leglog.nodes[%d].next = window.Leglog.nodes[%d];\n",
					lastChildId, child.Id))
				builder.WriteString(fmt.Sprintf("window.Leglog.nodes[%d].prev = window.Leglog.nodes[%d];\n",
					child.Id, lastChildId))
			}
			lastChildId = child.Id
		}
	}
	return builder.String()
}

func LoadNode(db *sql.DB, id int64) *Node {
	node := new(Node)
	db.QueryRow("select id, type, parent_id, prev_id, note_id, text from node where id = ?", id).
		Scan(&(node.Id), &(node.NodeType), &(node.ParentId), &(node.PrevId), &(node.NoteId), &(node.Text))
	return node
}
