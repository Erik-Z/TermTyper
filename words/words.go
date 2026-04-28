package words

import (
	_ "embed"
	"encoding/json"
	"math/rand/v2"
	"strings"
	"unicode"
	"unicode/utf8"
)

//TODO: Add quote test.

//go:embed embeds/common-english.json
var commonWords string

type MetaData struct {
	Name string
	Size int
}

type WordList struct {
	MetaData MetaData
	Words    []string
}

type WordGenerator struct {
	Count       int
	Punctuation bool
	poolsJson   map[string]WordList
	currentPool []string
}

func NewGenerator() WordGenerator {
	var gen WordGenerator
	gen.Count = 300
	gen.poolsJson = addEmbededSources(make(map[string]WordList, 0))

	return gen
}

func (gen *WordGenerator) Generate(wordListName string) []rune {
	pool := gen.poolsJson[wordListName].Words
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	amount := min(gen.Count, len(pool))
	words := pool[0:amount]

	gen.currentPool = pool

	if !gen.Punctuation {
		return []rune(strings.Join(words, " "))
	}

	return gen.generateWithPunctuation(words)
}

func (gen *WordGenerator) generateWithPunctuation(words []string) []rune {
	var result []rune
	sentenceStart := true

	for i, word := range words {
		currentWord := word
		if sentenceStart {
			currentWord = capitalizeFirst(word)
			sentenceStart = false
		}

		result = append(result, []rune(currentWord)...)

		if i == len(words)-1 {
			result = append(result, '.')
			break
		}

		nextPunct := gen.getNextPunctuation(i, len(words))
		result = append(result, nextPunct...)

		if containsEndPunct(nextPunct) {
			sentenceStart = true
		}
	}

	return result
}

func (gen *WordGenerator) getNextPunctuation(currentIdx, totalWords int) []rune {
	remaining := totalWords - currentIdx - 1

	if remaining <= 0 {
		return []rune{'.'}
	}

	roll := rand.Float64()

	switch {
	case roll < 0.15 && remaining > 5:
		return gen.generateParenthesis()
	case roll < 0.25 && remaining > 3:
		return gen.generateQuotedPhrase()
	case roll < 0.35:
		return []rune{';', ' '}
	case roll < 0.45:
		return []rune{',', ' '}
	case roll < 0.55 && remaining > 1:
		return []rune{'.', ' '}
	case roll < 0.65 && remaining > 1:
		return []rune{'!', ' '}
	case roll < 0.75 && remaining > 1:
		return []rune{'?', ' '}
	default:
		return []rune{' '}
	}
}

func (gen *WordGenerator) generateParenthesis() []rune {
	var words []rune
	words = append(words, ' ')
	words = append(words, '(')
	wordCount := rand.IntN(3) + 1
	for i := 0; i < wordCount; i++ {
		if i > 0 {
			words = append(words, ' ')
		}
		words = append(words, []rune(gen.randomWord())...)
	}
	words = append(words, ')')
	words = append(words, ' ')
	return words
}

func (gen *WordGenerator) generateQuotedPhrase() []rune {
	var words []rune
	words = append(words, ' ')
	words = append(words, '"')
	wordCount := rand.IntN(4) + 1
	for i := 0; i < wordCount; i++ {
		if i > 0 {
			words = append(words, ' ')
		}
		words = append(words, []rune(gen.randomWord())...)
	}
	words = append(words, '"')
	words = append(words, ' ')
	return words
}

func (gen *WordGenerator) randomWord() string {
	if len(gen.currentPool) == 0 {
		return "the"
	}
	return gen.currentPool[rand.IntN(len(gen.currentPool))]
}

func containsEndPunct(r []rune) bool {
	for _, c := range r {
		if c == '.' || c == '!' || c == '?' {
			return true
		}
	}
	return false
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

func addEmbededSources(sources map[string]WordList) map[string]WordList {
	var wordList WordList
	err := json.Unmarshal([]byte(commonWords), &wordList)

	if err != nil {
		panic(err)
	}

	sources["Common words"] = wordList

	return sources
}
