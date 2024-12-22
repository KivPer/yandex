package main

import (
	"errors"
	"strconv"
)

func Calc(expression string) (float64, error) {
	stackNum := []float64{}
	stackOp := []rune{}

	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	for i, char := range expression {
		if char == ' ' {
			continue
		}

		if char >= '0' && char <= '9' {
			start := i
			for i < len(expression) && (expression[i] >= '0' && expression[i] <= '9' || expression[i] == '.') {
				i++
			}
			numStr := expression[start:i]
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, errors.New("Invalid expression")
			}
			stackNum = append(stackNum, num)
			i--
		} else if char == '(' {
			stackOp = append(stackOp, '(')
		} else if char == ')' {
			for len(stackOp) > 0 && stackOp[len(stackOp)-1] != '(' {
				if err := calculate(&stackNum, &stackOp); err != nil {
					return 0, err
				}
			}
			if len(stackOp) == 0 || stackOp[len(stackOp)-1] != '(' {
				return 0, errors.New("Invalid expression")
			}
			stackOp = stackOp[:len(stackOp)-1]
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			for len(stackOp) > 0 && precedence[stackOp[len(stackOp)-1]] >= precedence[char] {
				if err := calculate(&stackNum, &stackOp); err != nil {
					return 0, err
				}
			}
			stackOp = append(stackOp, char)
		} else {
			return 0, errors.New("Invalid expression")
		}
	}

	for len(stackOp) > 0 {
		if stackOp[len(stackOp)-1] == '(' {
			return 0, errors.New("Invalid expression")
		}
		if err := calculate(&stackNum, &stackOp); err != nil {
			return 0, err
		}
	}

	if len(stackNum) != 1 {
		return 0, errors.New("Invalid expression")
	}

	return stackNum[0], nil
}

func calculate(stackNum *[]float64, stackOp *[]rune) error {
	if len(*stackNum) < 2 {
		return errors.New("Invalid expression")
	}

	num2 := (*stackNum)[len(*stackNum)-1]
	num1 := (*stackNum)[len(*stackNum)-2]
	*stackNum = (*stackNum)[:len(*stackNum)-2]

	op := (*stackOp)[len(*stackOp)-1]
	*stackOp = (*stackOp)[:len(*stackOp)-1]

	var result float64
	switch op {
	case '+':
		result = num1 + num2
	case '-':
		result = num1 - num2
	case '*':
		result = num1 * num2
	case '/':
		if num2 == 0 {
			return errors.New("Division by zero")
		}
		result = num1 / num2
	}

	*stackNum = append(*stackNum, result)

	return nil
}
