package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type CatalogItem struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Series      string   `json:"series"`
	PriceCny    int      `json:"priceCny"`
	Stock       int      `json:"stock"`
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
}

type FAQItem struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Body  string   `json:"body"`
}

type Document struct {
	ID       string       `json:"id"`
	Kind     string       `json:"kind"`
	Title    string       `json:"title"`
	Section  string       `json:"section,omitempty"`
	SourceID string       `json:"sourceId,omitempty"`
	Category string       `json:"category,omitempty"`
	Series   string       `json:"series,omitempty"`
	PriceCny int          `json:"priceCny,omitempty"`
	Stock    *int         `json:"stock,omitempty"`
	Tags     []string     `json:"tags"`
	Summary  string       `json:"summary"`
	Text     string       `json:"-"`
	Product  *CatalogItem `json:"product,omitempty"`
}

type Hit struct {
	Document
	Score      float64  `json:"score"`
	Highlights []string `json:"highlights"`
}

type AskRequest struct {
	Question string `json:"question"`
	Limit    int    `json:"limit"`
}

type AskResponse struct {
	Question  string   `json:"question"`
	Answer    string   `json:"answer"`
	Citations []Hit    `json:"citations"`
	Mode      string   `json:"mode"`
}

type Engine struct {
	docs     []Document
	chunks   []Chunk
	catalog  []CatalogItem
	index    map[string][]posting
	docTerms []map[string]int
}

type posting struct {
	doc int
	tf  int
}

func LoadEngine(dataDir string) (*Engine, error) {
	catalog, err := loadCatalog(filepath.Join(dataDir, "catalog.json"))
	if err != nil {
		return nil, err
	}
	faqs, err := loadFAQ(filepath.Join(dataDir, "faq.json"))
	if err != nil {
		return nil, err
	}

	var chunks []Chunk
	for _, item := range catalog {
		chunks = append(chunks, chunkProduct(item)...)
	}
	for _, faq := range faqs {
		chunks = append(chunks, chunkFAQ(faq)...)
	}
	knowledge, err := loadKnowledgeChunks(filepath.Join(dataDir, "docs"))
	if err != nil {
		return nil, err
	}
	chunks = append(chunks, knowledge...)

	docs := make([]Document, 0, len(chunks))
	for _, c := range chunks {
		docs = append(docs, documentFromChunk(c))
	}

	e := &Engine{docs: docs, chunks: chunks, catalog: catalog}
	e.buildIndex()
	return e, nil
}

func documentFromChunk(c Chunk) Document {
	summary := c.Text
	if c.Section != "" && !strings.HasPrefix(summary, c.Section) {
		summary = c.Section + "：" + clip(c.Text, 160)
	} else {
		summary = clip(c.Text, 180)
	}
	return Document{
		ID:       c.ID,
		Kind:     c.SourceKind,
		Title:    c.Title,
		Section:  c.Section,
		SourceID: c.SourceID,
		Category: c.Category,
		Series:   c.Series,
		PriceCny: c.PriceCny,
		Stock:    c.Stock,
		Tags:     c.Tags,
		Summary:  summary,
		Text:     strings.Join([]string{c.Title, c.Section, c.Text, strings.Join(c.Tags, " ")}, " "),
		Product:  c.Product,
	}
}

func (e *Engine) DocCount() int   { return len(e.docs) }
func (e *Engine) ChunkCount() int { return len(e.chunks) }

func (e *Engine) handleCatalog(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": e.catalog})
}

func (e *Engine) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	limit := parseLimit(r.URL.Query().Get("limit"), 6)
	hits := e.search(q, limit)
	writeJSON(w, http.StatusOK, map[string]any{"query": q, "hits": hits})
}

func (e *Engine) handleAsk(w http.ResponseWriter, r *http.Request) {
	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	q := strings.TrimSpace(req.Question)
	if q == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "question is required"})
		return
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 4
	}
	hits := e.search(q, limit)
	writeJSON(w, http.StatusOK, AskResponse{
		Question:  q,
		Answer:    composeAnswer(q, hits),
		Citations: hits,
		Mode:      "chunked-extractive-rag",
	})
}

func (e *Engine) buildIndex() {
	e.index = map[string][]posting{}
	e.docTerms = make([]map[string]int, len(e.docs))
	for i, doc := range e.docs {
		counts := map[string]int{}
		addTokens(counts, tokenize(doc.Title), 3)
		addTokens(counts, tokenize(doc.Section), 3)
		addTokens(counts, tokenize(strings.Join(doc.Tags, " ")), 2)
		addTokens(counts, tokenize(doc.Text), 1)
		e.docTerms[i] = counts
		for term, tf := range counts {
			e.index[term] = append(e.index[term], posting{doc: i, tf: tf})
		}
	}
}

func (e *Engine) search(query string, limit int) []Hit {
	if strings.TrimSpace(query) == "" {
		hits := make([]Hit, 0, min(limit, len(e.docs)))
		for i, doc := range e.docs {
			if i >= limit {
				break
			}
			hits = append(hits, Hit{Document: doc, Score: 0, Highlights: nil})
		}
		return hits
	}

	qTokens := tokenize(query)
	if len(qTokens) == 0 {
		return nil
	}

	scores := map[int]float64{}
	nDocs := float64(len(e.docs))
	for _, term := range qTokens {
		postings := e.index[term]
		if len(postings) == 0 {
			continue
		}
		idf := 1.0 + math.Log(1+nDocs/float64(len(postings)))
		for _, p := range postings {
			scores[p.doc] += float64(p.tf) * idf
		}
	}

	// Phrase / substring boost for Chinese queries that may not tokenize cleanly.
	foldedQ := strings.ToLower(query)
	for i, doc := range e.docs {
		hay := strings.ToLower(doc.Title + " " + doc.Summary + " " + doc.Text)
		if strings.Contains(hay, foldedQ) {
			scores[i] += 8
		}
		for _, tag := range doc.Tags {
			if strings.Contains(foldedQ, strings.ToLower(tag)) || strings.Contains(strings.ToLower(tag), foldedQ) {
				scores[i] += 3
			}
		}
	}

	type ranked struct {
		idx   int
		score float64
	}
	rank := make([]ranked, 0, len(scores))
	for idx, score := range scores {
		if score <= 0 {
			continue
		}
		rank = append(rank, ranked{idx: idx, score: score})
	}
	sort.Slice(rank, func(i, j int) bool {
		if rank[i].score == rank[j].score {
			left, right := e.docs[rank[i].idx], e.docs[rank[j].idx]
			if left.Title == right.Title {
				return left.ID < right.ID
			}
			return left.Title < right.Title
		}
		return rank[i].score > rank[j].score
	})
	if limit > len(rank) {
		limit = len(rank)
	}
	hits := make([]Hit, 0, limit)
	for _, item := range rank[:limit] {
		doc := e.docs[item.idx]
		hits = append(hits, Hit{
			Document:   doc,
			Score:      round2(item.score),
			Highlights: highlights(doc, qTokens, 2),
		})
	}
	return hits
}

func composeAnswer(question string, hits []Hit) string {
	if len(hits) == 0 {
		return "暂时没有在店内商品、FAQ 和邮票知识库里找到匹配内容。你可以换个说法，例如「日本明信片」「代寄怎么收费」「齿孔品相」「欧洲散票」。"
	}

	var b strings.Builder
	intent := classifyIntent(question)
	switch intent {
	case "price":
		b.WriteString("根据店内目录，和你问题最相关的价格如下：\n")
	case "stock":
		b.WriteString("当前可查到的相关库存如下：\n")
	case "shipping":
		b.WriteString("关于寄送与代寄，店内说明是：\n")
	case "condition":
		b.WriteString("关于品相、齿孔与保存，知识库里相关段落如下：\n")
	default:
		b.WriteString("结合店内商品、FAQ 与邮票知识库切块，建议这样理解：\n")
	}

	used := 0
	seen := map[string]bool{}
	for _, hit := range hits {
		if used >= 4 {
			break
		}
		key := hit.ID
		if seen[key] {
			continue
		}
		seen[key] = true
		switch hit.Kind {
		case "product":
			fmt.Fprintf(&b, "\n- 【商品】%s", hit.Title)
			if hit.Section != "" {
				fmt.Fprintf(&b, " · %s", hit.Section)
			}
			fmt.Fprintf(&b, "（%s / %s）售价 %d 元", hit.Category, hit.Series, hit.PriceCny)
			if hit.Stock != nil {
				fmt.Fprintf(&b, "，库存 %d。", *hit.Stock)
			} else {
				b.WriteString("。")
			}
			if hit.Summary != "" {
				b.WriteString(" ")
				b.WriteString(hit.Summary)
			}
		case "faq":
			fmt.Fprintf(&b, "\n- 【FAQ】%s：%s", hit.Title, hit.Summary)
		default:
			fmt.Fprintf(&b, "\n- 【知识】%s", hit.Title)
			if hit.Section != "" && hit.Section != hit.Title {
				fmt.Fprintf(&b, " / %s", hit.Section)
			}
			fmt.Fprintf(&b, "：%s", clip(hit.Summary, 220))
		}
		used++
	}
	b.WriteString("\n\n以上回答来自切块检索（chunked extractive RAG）。知识库用于解释集邮概念；价格与库存以商品条目为准。实寄一律走中国邮政。")
	return b.String()
}

func classifyIntent(q string) string {
	switch {
	case containsAny(q, "价格", "多少钱", "售价", "费用", "收费"):
		return "price"
	case containsAny(q, "库存", "有货", "卖完", "还有吗"):
		return "stock"
	case containsAny(q, "寄", "邮", "快递", "代寄", "发货", "物流"):
		return "shipping"
	case containsAny(q, "品相", "齿孔", "背胶", "保存", "护票", "赝品", "伪票", "目录"):
		return "condition"
	default:
		return "general"
	}
}

func highlights(doc Document, tokens []string, n int) []string {
	src := doc.Summary
	if src == "" {
		src = doc.Text
	}
	parts := splitSentences(src)
	var out []string
	for _, part := range parts {
		low := strings.ToLower(part)
		matched := false
		for _, tok := range tokens {
			if tok == "" {
				continue
			}
			if strings.Contains(low, tok) {
				matched = true
				break
			}
		}
		if matched {
			out = append(out, strings.TrimSpace(part))
			if len(out) >= n {
				break
			}
		}
	}
	if len(out) == 0 && strings.TrimSpace(src) != "" {
		out = append(out, clip(src, 80))
	}
	return out
}

func tokenize(text string) []string {
	text = strings.ToLower(text)
	var tokens []string
	var buf strings.Builder
	flush := func() {
		if buf.Len() == 0 {
			return
		}
		tokens = append(tokens, buf.String())
		buf.Reset()
	}
	for _, r := range text {
		switch {
		case unicode.Is(unicode.Han, r):
			flush()
			tokens = append(tokens, string(r))
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			buf.WriteRune(r)
		default:
			flush()
		}
	}
	flush()

	// Character bigrams help Chinese phrase matching.
	han := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if len([]rune(t)) == 1 && unicode.Is(unicode.Han, []rune(t)[0]) {
			han = append(han, t)
		}
	}
	for i := 0; i+1 < len(han); i++ {
		tokens = append(tokens, han[i]+han[i+1])
	}
	return tokens
}

func addTokens(dst map[string]int, tokens []string, weight int) {
	for _, t := range tokens {
		if t == "" {
			continue
		}
		dst[t] += weight
	}
}

func loadCatalog(path string) ([]CatalogItem, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read catalog: %w", err)
	}
	var items []CatalogItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("parse catalog: %w", err)
	}
	return items, nil
}

func loadFAQ(path string) ([]FAQItem, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read faq: %w", err)
	}
	var items []FAQItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("parse faq: %w", err)
	}
	return items, nil
}

func loadKnowledgeChunks(dir string) ([]Chunk, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read docs dir: %w", err)
	}
	var chunks []Chunk
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		sourceID := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		title := sourceID
		body := string(raw)
		if i := strings.Index(body, "\n"); i > 0 && strings.HasPrefix(body, "# ") {
			title = strings.TrimSpace(strings.TrimPrefix(body[:i], "# "))
		}
		tags := []string{"集邮", "邮票", "知识库"}
		switch {
		case strings.Contains(sourceID, "europe"):
			tags = append(tags, "欧洲", "散票", "复古")
		case strings.Contains(sourceID, "japan"):
			tags = append(tags, "日本", "明信片")
		case strings.Contains(sourceID, "forward"):
			tags = append(tags, "代寄", "中国邮政")
		case strings.Contains(sourceID, "condition"), strings.Contains(sourceID, "care"):
			tags = append(tags, "品相", "保存")
		}
		chunks = append(chunks, chunkMarkdown(sourceID, title, tags, body)...)
	}
	return chunks, nil
}

func parseLimit(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}
	if n > 20 {
		return 20
	}
	return n
}

func containsAny(s string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

var sentenceSplit = regexp.MustCompile(`[。！？!?\n]`)

func splitSentences(s string) []string {
	parts := sentenceSplit.Split(s, -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func clip(s string, n int) string {
	rs := []rune(strings.TrimSpace(s))
	if len(rs) <= n {
		return string(rs)
	}
	return string(rs[:n]) + "…"
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
