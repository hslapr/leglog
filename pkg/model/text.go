package model

import (
	"database/sql"
	"errors"
	"html/template"
	"strings"
	"time"
)

type Text struct {
	Id                int64
	RootId            int64
	CreationTimestamp int64
	Language          string
	Title             string
}

func (text *Text) Excerpt(max int) string {
	root := text.Root()
	stack := make([]*Node, 0)
	stack = append(stack, root)
	var builder strings.Builder
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		children := node.Children()
		if len(children) > 0 {
			for i, j := 0, len(children)-1; i < j; i, j = i+1, j-1 {
				children[i], children[j] = children[j], children[i]
			}
			stack = append(stack, children...)
		} else {
			builder.WriteString(node.Text)
			if builder.Len() > max {
				builder.WriteString("...")
				break
			}
		}
	}
	return builder.String()
}

func (text *Text) CreationTime(layout string) string {
	return time.Unix(text.CreationTimestamp, 0).Format(layout)
}

func NewText(language string, title string) *Text {
	return &Text{Language: language, Title: title}
}

func (text *Text) Save() {
	text.save(db)
}
func (text *Text) Js() template.JS {
	return template.JS(text.Root().Js())
}

func (text *Text) save(db *sql.DB) {
	if text.Id < 1 && text.RootId > 0 {
		r, _ := db.Exec("INSERT INTO text (root_id, creation_time, language, title) VALUES (?, ?, ?, ?)",
			text.RootId, time.Now().Unix(), text.Language, text.Title)
		text.Id, _ = r.LastInsertId()
	}
}

func (text *Text) Parse(s string) {
	scanner := NewScanner(strings.NewReader(s))
	root := NewNode("", ROOT)
	root.Save()
	text.RootId = root.Id
	p := NewNode("", PARAGRAPH)
	p.ParentId = root.Id
	var prev *Node
	for scanner.Next() {
		t, nodeType := scanner.Scan()
		if nodeType == PARAGRAPH {
			if p.Id < 1 {
				continue
			}
			prevParagraphId := p.Id
			p = NewNode("", PARAGRAPH)
			p.ParentId = root.Id
			p.PrevId = prevParagraphId
			prev = nil
		} else {
			if p.Id < 1 {
				p.Save()
			}
			node := NewNode(t, nodeType)
			node.ParentId = p.Id
			if prev == nil {
				node.PrevId = 0
			} else {
				node.PrevId = prev.Id
			}
			node.Save()
			prev = node
		}
	}
}

func LoadText(id int64) *Text {
	text := new(Text)
	db.QueryRow("select id, root_id, creation_time, language, title from text where id = ?", id).
		Scan(&(text.Id), &(text.RootId), &(text.CreationTimestamp), &(text.Language), &(text.Title))
	return text
}

func (text *Text) Load() {
	text.load(db)
}

func (text *Text) load(db *sql.DB) {
	if text.Id > 1 {
		db.QueryRow("select  root_id, creation_time, language, title from text where id = ?", text.Id).Scan(
			&(text.RootId), &(text.CreationTimestamp), &(text.Language), &(text.Title),
		)
	}
}

func (text *Text) Root() *Node {
	return text.root(db)
}

func (text *Text) root(db *sql.DB) *Node {
	return LoadNode(db, text.RootId)
}

func CreatePhrase(nodes []*Node, note *Note) (parent *Node, err error) {
	if len(nodes) < 1 {
		return nil, errors.New("nodes is empty")
	}
	for _, node := range nodes {
		if node.ParentId != nodes[0].ParentId {
			return nil, errors.New("not same parent")
		}
	}
	parent = &Node{
		ParentId: nodes[0].ParentId,
		NodeType: PHRASE,
		PrevId:   nodes[0].PrevId,
		NoteId:   note.Id,
	}
	parent.Save()
	nodes[0].PrevId = 0
	for _, node := range nodes {
		node.ParentId = parent.Id
		node.Save()
	}
	// db.Exec("UPDATE node SET parent_id = ? WHERE id = ?", parent.Id, nodes[len(nodes)-1].Id)
	db.Exec("UPDATE node SET prev_id = ? WHERE prev_id = ?", parent.Id, nodes[len(nodes)-1].Id)
	return parent, nil
}
