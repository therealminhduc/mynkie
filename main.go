package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	os.MkdirAll("public", 0755)

	staticOut := filepath.Join("public", "static")
	os.RemoveAll(staticOut)
	cpDir("static", staticOut)

	if err := renderPage("content/about.md", "public/index.html", "mynkie – about", "about"); err != nil {
		log.Fatal(err)
	}

	if err := renderPage("content/contact.md", "public/contact.html", "mynkie – contact", "contact"); err != nil {
		log.Fatal(err)
	}

	posts, err := loadBlogPosts("content/blog")
	if err != nil {
		log.Fatal(err)
	}

	if err := renderBlogIndex(posts); err != nil {
		log.Fatal(err)
	}
}
