package hw03frequencyanalysis

import (
	"slices"
	"strings"
)

func Top10(str string) []string {
	words := strings.Fields(str)

	wordsMap := make(map[string]int, len(words))

	for _, word := range words {
		wordsMap[word]++
	}

	//тут я долго не верила, что в GO нет сортировки по значениям мапы!!!
	type wordCount struct {
		word  string
		count int
	}
	resultStruct := make([]wordCount, 0, len(wordsMap))

	for w, c := range wordsMap {
		resultStruct = append(resultStruct, wordCount{w, c})
	}

	slices.SortFunc(resultStruct, func(a, b wordCount) int {
		if b.count != a.count {
			return b.count - a.count //a - b - сортировка по возрастанию, b - a - по убыванию, НЕ ЗНАЮ, КАК ЕЩЁ ЗАПОМНИТЬ
		}
		//"Если слова имеют одинаковую частоту, то должны быть отсортированы лексикографически"
		return strings.Compare(a.word, b.word)
	})

	result := make([]string, 0, 10)
	for i := 0; i < 10 && i < len(resultStruct); i++ {
		result = append(result, resultStruct[i].word)
	}
	return result
}
