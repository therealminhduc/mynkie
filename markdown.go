package main

import (
	"bytes"
	"log"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
}

type bytesBuffer struct {
	b *[]byte
}

func (w *bytesBuffer) Write(p []byte) (int, error) {
	*w.b = append(*w.b, p...)
	return len(p), nil
}

func parseFrontmatter(mdBytes []byte) (Frontmatter, []byte) {
	if !bytes.HasPrefix(mdBytes, []byte("---")) {
		return Frontmatter{}, mdBytes
	}

	afterFirst := mdBytes[3:]
	endIdx := bytes.Index(afterFirst, []byte("---"))
	if endIdx == -1 {
		return Frontmatter{}, mdBytes
	}

	var fm Frontmatter
	if err := yaml.Unmarshal(afterFirst[:endIdx], &fm); err != nil {
		log.Printf("warning: failed to parse frontmatter: %v", err)
		return Frontmatter{}, mdBytes
	}

	return fm, bytes.TrimLeft(afterFirst[endIdx+3:], "\n")
}

func convertMarkdown(mdBytes []byte) ([]byte, error) {
	var htmlBuf []byte
	bufWriter := &bytesBuffer{b: &htmlBuf}
	if err := goldmark.Convert(mdBytes, bufWriter); err != nil {
		return nil, err
	}
	return htmlBuf, nil
}
