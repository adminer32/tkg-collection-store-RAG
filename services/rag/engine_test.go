package main

import (
	"strings"
	"testing"
)

func TestSearchAndAsk(t *testing.T) {
	engine, err := LoadEngine("data")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if engine.ChunkCount() < 20 {
		t.Fatalf("expected knowledge chunks, got %d", engine.ChunkCount())
	}

	hits := engine.search("日本明信片", 3)
	if len(hits) == 0 {
		t.Fatal("expected hits for 日本明信片")
	}
	if !strings.HasPrefix(hits[0].SourceID, "jp-postcard") && hits[0].SourceID != "vintage-postcard" && !strings.Contains(hits[0].ID, "japan") {
		t.Fatalf("unexpected top hit %s source=%s", hits[0].ID, hits[0].SourceID)
	}

	ship := engine.search("快递怎么寄", 5)
	if len(ship) == 0 {
		t.Fatal("expected shipping faq hit")
	}
	foundFAQ := false
	for _, h := range ship {
		if h.Kind == "faq" {
			foundFAQ = true
		}
	}
	if !foundFAQ {
		t.Fatalf("expected faq citation, got %#v", ship)
	}

	ans := composeAnswer("代寄多少钱", ship)
	if !strings.Contains(ans, "中国邮政") && !strings.Contains(ans, "代寄") {
		t.Fatalf("weak answer: %s", ans)
	}
}

func TestKnowledgeChunks(t *testing.T) {
	engine, err := LoadEngine("data")
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	hits := engine.search("齿孔品相", 5)
	if len(hits) == 0 {
		t.Fatal("expected condition knowledge hits")
	}
	foundKnowledge := false
	for _, h := range hits {
		if h.Kind == "knowledge" {
			foundKnowledge = true
		}
	}
	if !foundKnowledge {
		t.Fatalf("expected knowledge chunk, got %#v", hits)
	}

	ans := composeAnswer("怎么看齿孔品相", hits)
	if !strings.Contains(ans, "【知识】") {
		t.Fatalf("expected knowledge citation in answer: %s", ans)
	}

	europe := engine.search("欧洲散票 1950", 4)
	if len(europe) == 0 {
		t.Fatal("expected europe vintage hits")
	}
}

func TestChunkMarkdownSplitsHeadings(t *testing.T) {
	chunks := chunkMarkdown("demo", "集邮入门", []string{"邮票"}, "# 集邮入门\n\n导语段落。\n\n## 新票与旧票\n\n新票指未经实寄的邮票。旧票带有销戳。这是足够长的一段说明文字，用来超过最小切块长度，避免被合并掉。\n\n## 套票\n\n套票是同一主题一次发行的若干枚，适合成组收藏，也适合入门者先选一个国家或一个年代。")
	if len(chunks) < 2 {
		t.Fatalf("expected heading-based chunks, got %d: %#v", len(chunks), chunks)
	}
	if chunks[0].Section == "" {
		t.Fatal("missing section")
	}
}
