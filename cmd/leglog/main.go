package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/hslapr/leglog/pkg/data"
	"github.com/hslapr/leglog/pkg/model"
	"github.com/hslapr/leglog/pkg/parser"
	"github.com/hslapr/leglog/pkg/util"
)

func noescape(str string) template.JS {
	return template.JS(str)
}

const TEMPLATE_PATH = "../../web/template/"

type ViewModel struct {
	Text     *data.Text
	TextHtml template.HTML
}

var indexTemplate = template.Must(template.ParseFiles(
	TEMPLATE_PATH+"index.html",
	TEMPLATE_PATH+"partial/layout.html",
))

var readTemplate = template.Must(template.New("v2readTemplate").Funcs(template.FuncMap{
	"noescape": noescape,
}).ParseFiles(
	TEMPLATE_PATH+"partial/layout.html",
	TEMPLATE_PATH+"read.html"))

var statisticsTemplate = template.Must(template.ParseFiles(
	TEMPLATE_PATH+"partial/layout.html",
	TEMPLATE_PATH+"statistics.html"))

var entryTemplate = template.Must(template.ParseFiles(
	TEMPLATE_PATH+"partial/layout.html",
	TEMPLATE_PATH+"entry/entry.html"))

var entryListTemplate = template.Must(template.ParseFiles(
	TEMPLATE_PATH+"partial/layout.html",
	TEMPLATE_PATH+"entry/index.html"))

var textListTemplate = template.Must(template.ParseFiles(
	TEMPLATE_PATH+"partial/layout.html",
	TEMPLATE_PATH+"text/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *ViewModel) {
	/* err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} */
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		indexTemplate.ExecuteTemplate(w, "layout", nil)
	} else if r.Method == http.MethodPost {
		lang := r.FormValue("language")
		text := r.FormValue("text")
		title := r.FormValue("title")
		text = util.SanitizeText(text)
		t := model.NewText(lang, title)
		t.Parse(text)
		t.Save()
		http.Redirect(w, r, fmt.Sprintf("/read/%d", t.Id), http.StatusFound)
	}
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	lang := string(r.FormValue("language"))
	text := string(r.FormValue("text"))
	t := parser.Parse(text, lang)
	t.Save()
	http.Redirect(w, r, fmt.Sprintf("/read/%d", t.Id), http.StatusFound)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Path[len("/read/"):], 10, 64)
	t := model.LoadText(id)
	readTemplate.ExecuteTemplate(w, "layout", t)
}

func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	data := model.Statistics()
	statisticsTemplate.ExecuteTemplate(w, "layout", data)
}

func entryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/entry/"):], 10, 64)
	if err != nil {
		entries := model.Entries()
		entryListTemplate.ExecuteTemplate(w, "layout", entries)
	} else {
		entry := new(model.Entry)
		entry.Id = id
		entry.Load()
		entryTemplate.ExecuteTemplate(w, "layout", entry)
	}
}

func textHandler(w http.ResponseWriter, r *http.Request) {
	texts := model.Texts()
	textListTemplate.ExecuteTemplate(w, "layout", texts)
}

func main() {
	fs := http.FileServer(http.Dir("../../web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/statistics", statisticsHandler)
	http.HandleFunc("/text/", textHandler)
	http.HandleFunc("/read/get-lemmas", getLemmasHandler)
	http.HandleFunc("/entry/", entryHandler)
	http.HandleFunc("/note/create", createNoteHandler)
	http.HandleFunc("/read/", readHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/note/query", queryNotesHandler)
	http.HandleFunc("/note/load", loadNoteHandler)
	http.HandleFunc("/note/bind", bindNoteHandler)
	http.HandleFunc("/note/unbind", unbindNoteHandler)
	http.HandleFunc("/note/update", updateNoteHandler)
	http.HandleFunc("/note/delete", deleteNoteHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
