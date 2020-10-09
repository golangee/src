package golang

import (
	"github.com/golangee/src/v2"
	"go/format"
	"unicode"
)

// Format tries to apply the gofmt rules to the given text. If it fails
// the error is returned and the string contains the text with line enumeration.
func Format(source []byte) ([]byte, error) {
	buf, err := format.Source(source)
	if err != nil {
		return []byte(src.WithLineNumbers(string(source))), err
	}

	return buf, nil
}

// MakePrivate converts ABc to aBc.
// Special cases:
//  * ID becomes id
func MakePrivate(str string) string {
	if len(str) == 0 {
		return str
	}

	switch str {
	case "ID":
		return "id"
	default:
		return string(unicode.ToLower(rune(str[0]))) + str[1:]
	}
}

// MakePublic converts aBc to ABc.
// Special cases:
//  * id becomes ID
func MakePublic(str string) string {
	if len(str) == 0 {
		return str
	}

	switch str {
	case "id":
		return "ID"
	default:
		return string(unicode.ToUpper(rune(str[0]))) + str[1:]
	}
}
