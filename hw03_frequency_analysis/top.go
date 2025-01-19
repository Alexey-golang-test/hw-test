package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	// Пустой текст, включая пробелы
	if strings.TrimSpace(str) == "" {
		return []string{}
	}

	// Разбивка текста на слова и подсчет частоты встречи каждого слова
	counter := make(map[string]uint)
	for _, s := range strings.Fields(str) {
		counter[s]++
	}

	// Переложим значения из мапы в слайс строк, чтобы в дальнейшем можно было отсортировать
	wordList := make([]string, 0, len(counter))
	for key := range counter {
		wordList = append(wordList, key)
	}

	// Сортировка
	sort.Slice(wordList, func(ii, jj int) bool {
		// При равенстве частоты слов, сортируется лексикографически
		if counter[wordList[ii]] == counter[wordList[jj]] {
			return wordList[ii] < wordList[jj]
		}

		return counter[wordList[ii]] > counter[wordList[jj]]
	})

	// Итоговый слайс
	recCount := 10
	if len(wordList) < 10 {
		recCount = len(wordList)
	}

	result := make([]string, 0, recCount)

	for ii := 0; ii < recCount; ii++ {
		result = append(result, wordList[ii])
	}

	return result
}
