package webui

import (
	"embed"
	"html/template"
	"net/http"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

//go:embed assets/*
var assets embed.FS

var libraryTemplate = template.Must(template.ParseFS(assets, "assets/library.html"))

type LibraryPageData struct {
	Bookmarks []bookmarkcore.Bookmark
	StoreMode string
	Error     string
}

func RenderLibrary(w http.ResponseWriter, status int, data LibraryPageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := libraryTemplate.ExecuteTemplate(w, "library.html", data); err != nil {
		http.Error(w, "Unable to render bookmark library", http.StatusInternalServerError)
	}
}

func Styles(w http.ResponseWriter, _ *http.Request) {
	content, err := assets.ReadFile("assets/library.css")
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}
