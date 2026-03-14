package main

import (
	"html/template"
	"log"
	"os"
	"time"

	"github.com/yuin/goldmark"
)

type PageData struct {
	Title string
	Content template.HTML
	Current string
	Year int
}

func renderPage(mdPath, outPath, title, current string) error {
	mdBytes, err := os.ReadFile(mdPath)
	if err != nil { return err }

	var htmlBuf []byte
	var bufWriter = &bytesBuffer{b: &htmlBuf}
	if err := goldmark.Convert(mdBytes, bufWriter); err != nil { return err }

	tmpl, err := template.ParseFiles("template/base.html")
    if err != nil { return err }

	f, err := os.Create(outPath)
    if err != nil { return err }
    defer f.Close()

	data := PageData{
        Title:   title,
        Content: template.HTML(htmlBuf),
        Current: current,
		Year: time.Now().Year(),
    }

	log.Println("generated ", outPath)

    return tmpl.ExecuteTemplate(f, "base", data)
}

func main() {
	os.MkdirAll("public", 0755)

    if err := renderPage("content/about.md", "public/index.html", "mynkie – about", "about"); err != nil {
        log.Fatal(err)
    }

	if err := renderPage("content/contact.md", "public/contact.html", "mynkie – contact", "contact"); err != nil {
        log.Fatal(err)
    }

	if err := renderPage("content/blog.md", "public/blog.html", "mynkie – blog", "blog"); err != nil {
        log.Fatal(err)
    }
}

type bytesBuffer struct {
	b *[]byte
}

func (w *bytesBuffer) Write(p []byte) (int, error) {
	*w.b = append(*w.b, p...)
	return len(p), nil
}