package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Post struct {
	Title   string
	Slug    string
	Date    string
	Content template.HTML
}

func slugify(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

func loadBlogPosts(dir string) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var posts []Post
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		slug := slugify(entry.Name())

		mdBytes, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		fm, content := parseFrontmatter(mdBytes)

		htmlBuf, err := convertMarkdown(content)
		if err != nil {
			return nil, err
		}

		title := fm.Title
		if title == "" {
			title = strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		}

		posts = append(posts, Post{
			Title:   title,
			Slug:    slug,
			Date:    fm.Date,
			Content: template.HTML(htmlBuf),
		})
	}

	return posts, nil
}

func renderBlogIndex(posts []Post) error {
	sort.Slice(posts, func(i, j int) bool {
		t1, _ := time.Parse("2006-01-02", posts[i].Date)
		t2, _ := time.Parse("2006-01-02", posts[j].Date)
		return t1.After(t2)
	})

	for _, post := range posts {
		os.MkdirAll(filepath.Join("public", "posts"), 0755)
		outPath := filepath.Join("public", "posts", post.Slug+".json")
		if err := os.WriteFile(outPath, []byte(post.Content), 0644); err != nil {
			return err
		}
		log.Println("generated", outPath)
	}

	var buf bytes.Buffer
	for i, post := range posts {
		date, _ := time.Parse("2006-01-02", post.Date)
		buf.WriteString("<article class=\"blog-post\">\n")
		buf.WriteString("<h2 class=\"blog-title\" data-slug=\"")
		buf.WriteString(post.Slug)
		buf.WriteString("\">")
		buf.WriteString(post.Title)
		buf.WriteString(" <span class=\"blog-date\">")
		buf.WriteString(date.Format("Jan 02, 2006"))
		buf.WriteString("</span></h2>\n")
		buf.WriteString("<div class=\"blog-content\" id=\"post-")
		buf.WriteString(post.Slug)
		buf.WriteString("\"")
		if i == 0 {
			buf.WriteString(" data-expanded=\"true\"")
		}
		buf.WriteString("></div>\n")
		buf.WriteString("</article>\n")
	}

	return executeTemplate("public/blog.html", "mynkie – blog", "blog", template.HTML(buf.String()))
}

func renderPage(mdPath, outPath, title, current string) error {
	mdBytes, err := os.ReadFile(mdPath)
	if err != nil {
		return err
	}

	htmlBuf, err := convertMarkdown(mdBytes)
	if err != nil {
		return err
	}

	return executeTemplate(outPath, title, current, template.HTML(htmlBuf))
}
