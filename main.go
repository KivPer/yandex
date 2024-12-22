package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
		if errors.Is(err, errors.New("Invalid expression")) {
			http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	response := Response{Result: strconv.FormatFloat(result, 'f', -1, 64)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

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
