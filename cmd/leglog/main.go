package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/hslapr/leglog/pkg/data"
	"github.com/hslapr/leglog/pkg/model"
	"github.com/hslapr/leglog/pkg/util"
)

const (
	ENTRY_PER_PAGE = 10
	TEXT_PER_PAGE  = 5
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
	// entry index list all entries
	if err != nil {
		var page int64 = 1
		pageStr := r.FormValue("page")
		if len(pageStr) > 0 {
			page, _ = strconv.ParseInt(pageStr, 10, 64)
		}
		cntEntry := model.EntryCount()
		cntPages := int64(math.Ceil(float64(cntEntry) / ENTRY_PER_PAGE))
		if page > cntPages {
			page = cntPages
		}
		data := make(map[string]interface{})
		entries := model.Entries((page-1)*ENTRY_PER_PAGE, ENTRY_PER_PAGE, "creation_time DESC")
		data["Entries"] = entries
		data["Page"] = page
		if page < cntPages {
			data["NextPage"] = page + 1
		} else {
			data["NextPage"] = 0
		}
		if page > 1 {
			data["PrevPage"] = page - 1
		} else {
			data["PrevPage"] = 0
		}
		data["LastPage"] = cntPages
		entryListTemplate.ExecuteTemplate(w, "layout", data)
	} else {
		entry := new(model.Entry)
		entry.Id = id
		entry.Load()
		entryTemplate.ExecuteTemplate(w, "layout", entry)
	}
}

func textHandler(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	pageStr := r.FormValue("page")
	if len(pageStr) > 0 {
		page, _ = strconv.ParseInt(pageStr, 10, 64)
	}
	cntText := model.TextCount()
	cntPages := int64(math.Ceil(float64(cntText) / TEXT_PER_PAGE))
	if page > cntPages {
		page = cntPages
	}
	texts := model.Texts((page-1)*TEXT_PER_PAGE, TEXT_PER_PAGE, "creation_time DESC")
	data := make(map[string]interface{})
	data["Texts"] = texts
	data["Page"] = page
	if page < cntPages {
		data["NextPage"] = page + 1
	} else {
		data["NextPage"] = 0
	}
	if page > 1 {
		data["PrevPage"] = page - 1
	} else {
		data["PrevPage"] = 0
	}
	data["LastPage"] = cntPages
	textListTemplate.ExecuteTemplate(w, "layout", data)
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
