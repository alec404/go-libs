package stringutil

import (
	"reflect"
	"testing"
)

func TestSplitToStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "normal case",
			input:    "a,b,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with spaces",
			input:    "a, b, c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "with empty parts",
			input:    "a,,b,  ,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "single value",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			name:     "trailing comma",
			input:    "a,b,c,",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitToStringSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitToStringSlice(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitToIntSlice(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []int
		expectErr bool
	}{
		{
			name:      "normal case",
			input:     "1,2,3",
			expected:  []int{1, 2, 3},
			expectErr: false,
		},
		{
			name:      "with spaces",
			input:     "1, 2, 3",
			expected:  []int{1, 2, 3},
			expectErr: false,
		},
		{
			name:      "empty string",
			input:     "",
			expected:  []int{},
			expectErr: false,
		},
		{
			name:      "with empty parts",
			input:     "1,,2,  ,3",
			expected:  []int{1, 2, 3},
			expectErr: false,
		},
		{
			name:      "single value",
			input:     "42",
			expected:  []int{42},
			expectErr: false,
		},
		{
			name:      "negative numbers",
			input:     "-1, 0, 1",
			expected:  []int{-1, 0, 1},
			expectErr: false,
		},
		{
			name:      "invalid input",
			input:     "1,abc,3",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "trailing comma",
			input:     "1,2,3,",
			expected:  []int{1, 2, 3},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SplitToIntSlice(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("SplitToIntSlice(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("SplitToIntSlice(%q) unexpected error: %v", tt.input, err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("SplitToIntSlice(%q) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		list     []string
		expected bool
	}{
		{
			name:     "found in middle",
			target:   "b",
			list:     []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "found at beginning",
			target:   "a",
			list:     []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "found at end",
			target:   "c",
			list:     []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "not found",
			target:   "d",
			list:     []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "empty list",
			target:   "a",
			list:     []string{},
			expected: false,
		},
		{
			name:     "empty target in list",
			target:   "",
			list:     []string{"a", "", "b"},
			expected: true,
		},
		{
			name:     "empty target not in list",
			target:   "",
			list:     []string{"a", "b"},
			expected: false,
		},
		{
			name:     "single element found",
			target:   "hello",
			list:     []string{"hello"},
			expected: true,
		},
		{
			name:     "single element not found",
			target:   "hello",
			list:     []string{"world"},
			expected: false,
		},
		{
			name:     "case sensitive",
			target:   "Hello",
			list:     []string{"hello", "world"},
			expected: false,
		},
		{
			name:     "with spaces",
			target:   " hello ",
			list:     []string{"hello", " hello "},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringInSlice(tt.target, tt.list)
			if result != tt.expected {
				t.Errorf("StringInSlice(%q, %v) = %v, want %v", tt.target, tt.list, result, tt.expected)
			}
		})
	}
}

func TestHideStr(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		percent  int
		hide     string
		expected string
	}{
		{
			name:     "normal case 50%",
			str:      "1234567890",
			percent:  50,
			hide:     "*",
			expected: "123*****90", // hideLen=5, mid=5, start=2, end=7
		},
		{
			name:     "email address",
			str:      "test@example.com",
			percent:  50,
			hide:     "*",
			expected: "t**t@example.com", // only "test" is processed
		},
		{
			name:     "percent 0 - no hiding",
			str:      "1234567890",
			percent:  0,
			hide:     "*",
			expected: "1234567890",
		},
		{
			name:     "percent 100 - hide all",
			str:      "1234567890",
			percent:  100,
			hide:     "*",
			expected: "**********",
		},
		{
			name:     "percent 100 - email",
			str:      "test@example.com",
			percent:  100,
			hide:     "*",
			expected: "****@example.com",
		},
		{
			name:     "empty string",
			str:      "",
			percent:  50,
			hide:     "*",
			expected: "",
		},
		{
			name:     "single character",
			str:      "a",
			percent:  50,
			hide:     "*",
			expected: "a", // hideLen=0, no change
		},
		{
			name:     "two characters 50%",
			str:      "ab",
			percent:  50,
			hide:     "*",
			expected: "a*", // hideLen=1, mid=1, start=0, end=1
		},
		{
			name:     "unicode characters",
			str:      "你好世界",
			percent:  50,
			hide:     "*",
			expected: "你**界", // hideLen=2, mid=2, start=1, end=3
		},
		{
			name:     "phone number",
			str:      "13812345678",
			percent:  40,
			hide:     "*",
			expected: "138****5678", // hideLen=4, mid=5, start=3, end=7
		},
		{
			name:     "custom hide string",
			str:      "1234567890",
			percent:  50,
			hide:     "X",
			expected: "123XXXXX90",
		},
		{
			name:     "empty hide string",
			str:      "1234567890",
			percent:  50,
			hide:     "",
			expected: "12390", // removes middle chars without replacement
		},
		{
			name:     "percent 25",
			str:      "1234567890",
			percent:  25,
			hide:     "*",
			expected: "1234**7890", // hideLen=2, mid=5, start=4, end=6
		},
		{
			name:     "percent 75",
			str:      "1234567890",
			percent:  75,
			hide:     "*",
			expected: "12*******0", // hideLen=7, mid=5, start=1, end=8
		},
		{
			name:     "negative percent - treated as 0",
			str:      "1234567890",
			percent:  -10,
			hide:     "*",
			expected: "1234567890",
		},
		{
			name:     "percent over 100",
			str:      "1234567890",
			percent:  150,
			hide:     "*",
			expected: "**********",
		},
		{
			name:     "short string with high percent",
			str:      "abc",
			percent:  80,
			hide:     "*",
			expected: "**c", // hideLen=2, mid=1, start=0, end=2
		},
		{
			name:     "chinese email",
			str:      "张三@example.com",
			percent:  50,
			hide:     "*",
			expected: "张*@example.com", // hideLen=1, mid=1, start=0, end=1
		},
		{
			name:     "id card simulation",
			str:      "110101199001011234",
			percent:  40,
			hide:     "*",
			expected: "110101*******11234", // hideLen=7, mid=9, start=5, end=12
		},
		{
			name:     "three characters",
			str:      "abc",
			percent:  33,
			hide:     "*",
			expected: "abc", // hideLen=0 (3*33/100=0), no change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HideStr(tt.str, tt.percent, tt.hide)
			if result != tt.expected {
				t.Errorf("HideStr(%q, %d, %q) = %q, want %q", tt.str, tt.percent, tt.hide, result, tt.expected)
			}
		})
	}
}
