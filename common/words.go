package common

import (
	"embed"
	"math/rand"
	"strings"
)

//go:embed words.txt
var wordsFile embed.FS

var dictionary []string

func RandomWord() string {
	return dictionary[rand.Intn(len(dictionary))]
}

func GeneratePhrase(numWords int) string {
	var words []string
	for i := 0; i < numWords; i++ {
		words = append(words, RandomWord())
	}
	return strings.Join(words, " ")
}

func init() {
	data, err := wordsFile.ReadFile("words.txt")
	if err != nil {
		panic(err)
	}

	trimmed := strings.TrimSpace((string(data)))
	dictionary = strings.Split(trimmed, " ")
}
