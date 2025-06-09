package textsplitter

import (
	"strings"
	"unicode/utf8"
)

type RecursiveCharacter struct {
	Separators   []string
	ChunkSize    int
	ChunkOverlap int
}

// Split takes a long text and returns chunks of at most ChunkSize characters,
// with ChunkOverlap characters of overlap between adjacent chunks.
func (rs *RecursiveCharacter) Split(text string) []string {
	pieces := rs.recursiveSplit(text, rs.Separators)
	return rs.mergeWithOverlap(pieces)
}

// recursiveSplit breaks text into pieces ≤ ChunkSize using separators.
func (rs *RecursiveCharacter) recursiveSplit(text string, seps []string) []string {
	if utf8.RuneCountInString(text) <= rs.ChunkSize {
		return []string{text}
	}
	if len(seps) == 0 {
		return rs.hardSplit(text)
	}

	sep := seps[0]
	parts := strings.Split(text, sep)
	if len(parts) == 1 {
		// Separator not found; try next one
		return rs.recursiveSplit(text, seps[1:])
	}

	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		sub := rs.recursiveSplit(part, seps[1:])
		result = append(result, sub...)
	}
	return result
}

// hardSplit force‑splits text into ChunkSize‑sized pieces (no overlap here).
func (rs *RecursiveCharacter) hardSplit(text string) []string {
	var pieces []string
	runes := []rune(text)
	for start := 0; start < len(runes); start += rs.ChunkSize {
		end := start + rs.ChunkSize
		if end > len(runes) {
			end = len(runes)
		}
		pieces = append(pieces, string(runes[start:end]))
	}
	return pieces
}

// mergeWithOverlap stitches atomic pieces into final chunks with overlap.
func (rs *RecursiveCharacter) mergeWithOverlap(pieces []string) []string {
	var chunks []string
	var current []rune

	pushChunk := func() {
		if len(current) > 0 {
			chunks = append(chunks, string(current))
		}
	}

	for _, piece := range pieces {
		pr := []rune(piece)
		if len(current)+len(pr) > rs.ChunkSize {
			// emit current
			pushChunk()
			// retain overlap
			if rs.ChunkOverlap > 0 && rs.ChunkOverlap < len(current) {
				current = current[len(current)-rs.ChunkOverlap:]
			} else {
				current = []rune{}
			}
		}
		current = append(current, pr...)
	}
	pushChunk()
	return chunks
}
