package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	runesIn := []rune(str)

	var resStr strings.Builder
	isPrevNum := true

	for ii, ss := range runesIn {
		curSymbol := string(ss)

		// Текущий символ число
		if _, err := strconv.Atoi(curSymbol); err == nil {
			// Если предыдущий символ был числом, или это первый символ в переданной строке, то ошибка
			if isPrevNum {
				return "", ErrInvalidString
			}

			// Если предыдущий символ не был числом, то распаковка уже была произведена и это число можно пропустить
			isPrevNum = true
			continue
		}

		// Здесь и далее текущий символ не является числом
		isPrevNum = false

		// Обработка не последнего символа
		if ii < len(runesIn)-1 {
			num, err := strconv.Atoi(string(runesIn[ii+1]))
			if err == nil {
				resStr.WriteString(strings.Repeat(curSymbol, num))
			} else {
				resStr.WriteString(curSymbol)
			}
		} else {
			// Последний символ
			resStr.WriteString(curSymbol)
		}
	}

	return resStr.String(), nil
}
