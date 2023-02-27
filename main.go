package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type OperandType int

type Operand struct {
	value       int
	operandType OperandType
}

type Expression struct {
	operand1 Operand
	operand2 Operand
	operator func(int, int) int
}

type ArabicRoman struct {
	Arabic int
	Roman  string
}

const (
	ROMAN OperandType = iota
	ARABIC
	UNKNOWN
)

var tableArabicRoman = []ArabicRoman{
	{1, "I"},
	{4, "IV"},
	{5, "V"},
	{9, "IV"},
	{10, "X"},
	{50, "L"},
	{40, "XL"},
	{90, "XC"},
	{100, "C"},
}

var operators = map[string]func(int, int) int{
	"+": func(value1, value2 int) int { return value1 + value2 },
	"-": func(value1, value2 int) int { return value1 - value2 },
	"/": func(value1, value2 int) int { return value1 / value2 },
	"*": func(value1, value2 int) int { return value1 * value2 },
}

func main() {
	inputReader := bufio.NewScanner(os.Stdin)
	inputReader.Scan()
	input := inputReader.Text()
	inputArr := strings.Split(input, " ")

	if len(inputArr) != 3 {
		panic("формат математической операции не удовлетворяет условию — два операнда и один оператор (+, -, /, *)")
	}

	if ex, err := buildExpression(inputArr); err == nil {
		var result = ex.operator(ex.operand1.value, ex.operand2.value)

		if ex.operand1.operandType == ARABIC {
			fmt.Printf("%v", result)
		} else if ex.operand1.operandType == ROMAN {
			if result < 1 {
				panic("в римской системе нет отрицательных чисел")
			}
			fmt.Print(toRoman(result))
		}
	} else {
		panic(err)
	}
}

func buildExpression(inputArr []string) (Expression, error) {
	var value1 = buildInputValue(inputArr[0])
	var value2 = buildInputValue(inputArr[2])
	var operator func(int, int) int

	if ex, ok := operators[inputArr[1]]; ok {
		operator = ex
	} else {
		return Expression{}, errors.New("некорректный оператор")
	}

	if value1.operandType == UNKNOWN || value2.operandType == UNKNOWN {
		return Expression{}, errors.New("некорректный операнд")
	}

	if value1.operandType != value2.operandType {
		return Expression{}, errors.New("используются одновременно разные системы счисления")
	}

	if checkRangeValue(value1.value) || checkRangeValue(value2.value) {
		return Expression{}, errors.New("операнд вышел за пределы границ")
	}

	return Expression{value1, value2, operator}, nil
}

func checkRangeValue(value int) bool {
	return value < 1 || value > 10
}

func buildInputValue(value string) (result Operand) {
	if number, err := strconv.Atoi(value); err == nil {
		result = Operand{number, ARABIC}
	} else if number, err := toArabic(value); err == nil {
		result = Operand{number, ROMAN}
	} else {
		result = Operand{0, UNKNOWN}
	}

	return result
}

func toRoman(number int) (result string) {
	var N = number
	var i = 8

	for N > 0 {
		for tableArabicRoman[i].Arabic > N {
			i--
		}
		result += tableArabicRoman[i].Roman
		N -= tableArabicRoman[i].Arabic
	}

	return result
}

func toArabic(romanNumber string) (result int, err error) {
	var matchCounter = 0
	for i := 0; i < len(romanNumber); i++ {
		for j := 0; j < len(tableArabicRoman); j++ {
			if tableArabicRoman[j].Roman == romanNumber[i:i+1] {
				matchCounter++
				break
			}
		}
	}

	if matchCounter != len(romanNumber) {
		return 0, errors.New("некорректное выражение")
	}

	var i = 8
	var p = 1
	for p <= len(romanNumber) {
		for i > 0 {
			if p-1+len(tableArabicRoman[i].Roman) <= len(romanNumber) {
				if romanNumber[p-1:p-1+len(tableArabicRoman[i].Roman)] == tableArabicRoman[i].Roman {
					break
				}
			}

			i--
		}

		result += tableArabicRoman[i].Arabic
		p += len(tableArabicRoman[i].Roman)
	}
	return result, nil
}
