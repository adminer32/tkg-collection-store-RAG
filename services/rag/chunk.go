package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	maxChunkRunes    = 420
	minChunkRunes    = 80
	chunkOverlapSent = 1
)

type Chunk struct {
	ID         string
	SourceID   string
	SourceKind string
	Title      string
	Section    string
	Text       string
	Tags       []string
	Category   string
	Series     string
	PriceCny   int
	Stock      *int
	Product    *CatalogItem
}

func chunkMarkdown(sourceID, title string, tags []string, markdown string) []Chunk {
	markdown = strings.ReplaceAll(markdown, "\r\n", "\n")
	lines := strings.Split(markdown, "\n")
	type section struct {
		heading string
		body    []string
	}
	sections := []section{{heading: title}}
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			heading := strings.TrimSpace(strings.TrimPrefix(line, "# "))
			if heading != "" && sections[len(sections)-1].heading == title && len(sections[len(sections)-1].body) == 0 {
				sections[len(sections)-1].heading = heading
				continue
			}
		}
		if strings.HasPrefix(line, "## ") {
			heading := strings.TrimSpace(strings.TrimPrefix(line, "## "))
			sections = append(sections, section{heading: heading})
			continue
		}
		sections[len(sections)-1].body = append(sections[len(sections)-1].body, line)
	}

	var chunks []Chunk
	for _, sec := range sections {
		body := strings.TrimSpace(strings.Join(sec.body, "\n"))
		if body == "" {
			continue
		}
		parts := splitLongText(body, maxChunkRunes)
		for i, part := range parts {
			id := fmt.Sprintf("%s#%s", sourceID, slug(sec.heading))
			if len(parts) > 1 {
				id = fmt.Sprintf("%s-%d", id, i+1)
			}
			chunks = append(chunks, Chunk{
				ID:         id,
				SourceID:   sourceID,
				SourceKind: "knowledge",
				Title:      title,
				Section:    sec.heading,
				Text:       part,
				Tags:       tags,
			})
		}
	}
	return chunks
}

func chunkProduct(item CatalogItem) []Chunk {
	stock := item.Stock
	meta := strings.Join([]string{
		item.Title, item.Category, item.Series, item.Summary, strings.Join(item.Tags, " "),
		fmt.Sprintf("售价 %d 元 库存 %d", item.PriceCny, item.Stock),
	}, " ")
	chunks := []Chunk{{
		ID:         item.ID + "#summary",
		SourceID:   item.ID,
		SourceKind: "product",
		Title:      item.Title,
		Section:    "商品摘要",
		Text:       strings.TrimSpace(meta),
		Tags:       item.Tags,
		Category:   item.Category,
		Series:     item.Series,
		PriceCny:   item.PriceCny,
		Stock:      &stock,
		Product:    &item,
	}}
	if desc := strings.TrimSpace(item.Description); desc != "" && desc != item.Summary {
		chunks = append(chunks, Chunk{
			ID:         item.ID + "#description",
			SourceID:   item.ID,
			SourceKind: "product",
			Title:      item.Title,
			Section:    "商品说明",
			Text:       desc,
			Tags:       item.Tags,
			Category:   item.Category,
			Series:     item.Series,
			PriceCny:   item.PriceCny,
			Stock:      &stock,
			Product:    &item,
		})
	}
	return chunks
}

func chunkFAQ(item FAQItem) []Chunk {
	text := strings.TrimSpace(item.Title + "\n" + item.Body)
	return []Chunk{{
		ID:         item.ID + "#body",
		SourceID:   item.ID,
		SourceKind: "faq",
		Title:      item.Title,
		Section:    item.Title,
		Text:       text,
		Tags:       item.Tags,
	}}
}

func splitLongText(text string, limit int) []string {
	paras := splitParagraphs(text)
	var packed []string
	var buf strings.Builder
	flush := func() {
		s := strings.TrimSpace(buf.String())
		if s != "" {
			packed = append(packed, s)
		}
		buf.Reset()
	}
	for _, p := range paras {
		if utf8.RuneCountInString(p) > limit {
			flush()
			packed = append(packed, splitBySentences(p, limit)...)
			continue
		}
		next := p
		if buf.Len() > 0 {
			next = buf.String() + "\n\n" + p
		}
		if utf8.RuneCountInString(next) > limit && buf.Len() > 0 {
			flush()
			buf.WriteString(p)
			continue
		}
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(p)
	}
	flush()

	if chunkOverlapSent > 0 && len(packed) > 1 {
		for i := 1; i < len(packed); i++ {
			prevSents := splitSentences(packed[i-1])
			if len(prevSents) == 0 {
				continue
			}
			overlap := strings.TrimSpace(prevSents[len(prevSents)-1])
			if overlap != "" && !strings.HasPrefix(strings.TrimSpace(packed[i]), overlap) {
				packed[i] = overlap + "\n" + packed[i]
			}
		}
	}

	out := packed[:0]
	for _, p := range packed {
		p = strings.TrimSpace(p)
		if utf8.RuneCountInString(p) >= minChunkRunes || len(packed) == 1 {
			out = append(out, p)
		} else if len(out) > 0 {
			out[len(out)-1] = out[len(out)-1] + "\n" + p
		} else {
			out = append(out, p)
		}
	}
	return out
}

func splitParagraphs(text string) []string {
	raw := strings.Split(text, "\n")
	var paras []string
	var buf []string
	flush := func() {
		s := strings.TrimSpace(strings.Join(buf, "\n"))
		if s != "" {
			paras = append(paras, s)
		}
		buf = buf[:0]
	}
	for _, line := range raw {
		if strings.TrimSpace(line) == "" {
			flush()
			continue
		}
		buf = append(buf, line)
	}
	flush()
	return paras
}

func splitBySentences(text string, limit int) []string {
	sents := splitSentences(text)
	if len(sents) == 0 {
		return []string{clip(text, limit)}
	}
	var out []string
	var buf strings.Builder
	for _, s := range sents {
		next := s
		if buf.Len() > 0 {
			next = buf.String() + s
		}
		if utf8.RuneCountInString(next) > limit && buf.Len() > 0 {
			out = append(out, strings.TrimSpace(buf.String()))
			buf.Reset()
			buf.WriteString(s)
			continue
		}
		buf.WriteString(s)
	}
	if buf.Len() > 0 {
		out = append(out, strings.TrimSpace(buf.String()))
	}
	return out
}

func slug(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	if s == "" {
		return "chunk"
	}
	return s
}
