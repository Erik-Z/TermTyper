package words

import (
	_ "embed"
	"encoding/json"
	"math/rand/v2"
	"strings"
)

//TODO: Add capitalization and punctuation to the random words
//TODO: Add quote test.
//TODO: Add zen mode

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
	Count     int
	poolsJson map[string]WordList
}

func NewGenerator() WordGenerator {
	var gen WordGenerator
	gen.Count = 300
	gen.poolsJson = addEmbededSources(make(map[string]WordList, 0))

	return gen
}

func (gen WordGenerator) Generate(wordListName string) []rune {
	pool := gen.poolsJson[wordListName].Words
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	amount := min(gen.Count, len(pool))
	words := pool[0:amount]

	return []rune(strings.Join(words, " "))
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
