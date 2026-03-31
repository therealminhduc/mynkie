package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"
)

type PageData struct {
	Title   string
	Content template.HTML
	Current string
	Year    int
}

func executeTemplate(outPath, title, current string, content template.HTML) error {
	tmpl, err := template.ParseFiles("template/base.html")
	if err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := PageData{
		Title:   title,
		Content: content,
		Current: current,
		Year:    time.Now().Year(),
	}

	log.Println("generated", outPath)
	return tmpl.ExecuteTemplate(f, "base", data)
}

func cpDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(target, data, 0644)
	})
}
