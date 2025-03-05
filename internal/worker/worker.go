package worker

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type Task struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result,omitempty"`
	Status     string `json:"status"`
}

func HandleTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := EvaluateExpression(task.Expression)
	if err != nil {
		http.Error(w, "Invalid expression", http.StatusBadRequest)
		return
	}

	task.Result = strconv.FormatFloat(result, 'f', -1, 64)
	task.Status = "completed"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func EvaluateExpression(expr string) (float64, error) {
	expr = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, expr)

	return parseExpression(&expr)
}

func parseExpression(expr *string) (float64, error) {
	var values []float64
	var ops []byte

	for len(*expr) > 0 {
		c := (*expr)[0]

		if unicode.IsLetter(rune(c)) {
			funcName, rest := parseFunctionName(*expr)
			*expr = rest
			if len(*expr) == 0 || (*expr)[0] != '(' {
				return 0, errors.New("expected '(' after function name")
			}
			*expr = (*expr)[1:]
			arg, err := parseExpression(expr)
			if err != nil {
				return 0, err
			}
			values = append(values, applyFunction(funcName, arg))
		} else if unicode.IsDigit(rune(c)) || c == '.' {
			num, rest, err := parseNumber(*expr)
			if err != nil {
				return 0, err
			}
			values = append(values, num)
			*expr = rest
		} else if c == '(' {
			*expr = (*expr)[1:]
			num, err := parseExpression(expr)
			if err != nil {
				return 0, err
			}
			values = append(values, num)
		} else if c == ')' {
			*expr = (*expr)[1:]
			break
		} else if strings.ContainsRune("+-*/", rune(c)) {
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(c) {
				values = applyOperator(ops[len(ops)-1], values)
				ops = ops[:len(ops)-1]
			}
			ops = append(ops, c)
			*expr = (*expr)[1:]
		} else {
			return 0, errors.New("invalid character in expression")
		}
	}

	for len(ops) > 0 {
		values = applyOperator(ops[len(ops)-1], values)
		ops = ops[:len(ops)-1]
	}

	if len(values) != 1 {
		return 0, errors.New("invalid expression")
	}
	return values[0], nil
}

func parseFunctionName(expr string) (string, string) {
	for i := 0; i < len(expr); i++ {
		if !unicode.IsLetter(rune(expr[i])) {
			return expr[:i], expr[i:]
		}
	}
	return expr, ""
}

func applyFunction(name string, arg float64) float64 {
	switch name {
	case "sqrt":
		return math.Sqrt(arg)
	case "sin":
		return math.Sin(arg)
	case "cos":
		return math.Cos(arg)
	case "log":
		return math.Log(arg)
	default:
		log.Printf("Unknown function: %s", name)
		return 0
	}
}

func parseNumber(expr string) (float64, string, error) {
	i := 0
	for i < len(expr) && (unicode.IsDigit(rune(expr[i])) || expr[i] == '.') {
		i++
	}
	num, err := strconv.ParseFloat(expr[:i], 64)
	if err != nil {
		return 0, "", errors.New("invalid number format")
	}
	return num, expr[i:], nil
}

func precedence(op byte) int {
	if op == '+' || op == '-' {
		return 1
	}
	if op == '*' || op == '/' {
		return 2
	}
	return 0
}

func applyOperator(op byte, values []float64) []float64 {
	if len(values) < 2 {
		return values
	}
	b := values[len(values)-1]
	a := values[len(values)-2]
	values = values[:len(values)-2]

	var result float64
	switch op {
	case '+':
		result = a + b
	case '-':
		result = a - b
	case '*':
		result = a * b
	case '/':
		if b == 0 {
			log.Println("Division by zero")
			result = 0
		} else {
			result = a / b
		}
	}
	return append(values, result)
}
