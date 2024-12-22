package calculator

import (
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		expectErr  bool
	}{
		{"2 + 2", 4, false},
		{"2 - 1", 1, false},
		{"2 * 3", 6, false},
		{"6 / 2", 3, false},
		{"2 + 3 * 4", 14, false},
		{"(1 + 2) * 3", 9, false},
		{"2 + (3 * 4)", 14, false},
		{"(2 + 3) * (4 + 1)", 25, false},
		{"10 / 0", 0, true}, // Деление на ноль
		{"2 + x", 0, true},  // Некорректное выражение
		{"", 0, true},       // Пустое выражение
	}

	for _, test := range tests {
		result, err := Calc(test.expression)

		if (err != nil) != test.expectErr {
			t.Errorf("Calc(%q) error = %v, expectErr %v", test.expression, err, test.expectErr)
		}
		if !test.expectErr && result != test.expected {
			t.Errorf("Calc(%q) = %v, want %v", test.expression, result, test.expected)
		}
	}
}
