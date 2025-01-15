package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCounter struct {
	word  string
	count uint
}

func Top10(str string) []string {
	// Пустой текст, включая пробелы
	if strings.TrimSpace(str) == "" {
		return []string{}
	}

	// Разбивка текста на слова и подсчет частоты встречи каждого слова
	counter := make(map[string]uint)
	for _, s := range strings.Fields(str) {
		// Проверка ключа
		_, exist := counter[s]
		if !exist {
			counter[s] = 0
		}

		// Слово уже есть в мапе
		counter[s]++
	}

	// Переложим полученную мапу в слайс структур wordCounter, чтобы в дальнейшем можно было отсортировать
	wordCounterList := make([]wordCounter, 0, len(counter))
	for key, value := range counter {
		wordCounterList = append(wordCounterList, wordCounter{key, value})
	}

	// Сортировка
	sort.Slice(wordCounterList, func(ii, jj int) bool {
		// Слово ii встречается чаще слова jj
		if wordCounterList[ii].count > wordCounterList[jj].count {
			return true
		}

		// При равенстве частоты слов, сортируется лексикографически
		if wordCounterList[ii].count == wordCounterList[jj].count && wordCounterList[ii].word < wordCounterList[jj].word {
			return true
		}

		// Для остальных случаев false
		return false
	})

	// Итоговый слайс
	recCount := 10
	if len(wordCounterList) < 10 {
		recCount = len(wordCounterList)
	}

	result := make([]string, 0, recCount)

	for ii := 0; ii < recCount; ii++ {
		result = append(result, wordCounterList[ii].word)
	}

	return result
}
