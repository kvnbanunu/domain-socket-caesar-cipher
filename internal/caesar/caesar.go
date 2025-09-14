package caesar

import (
	"log"
	"unicode"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
)

func cipher(c rune, shift rune, limit int) rune {
	var result rune

	// this block seems redundant, but we need to keep the case the same
	// otherwise it may overflow into being always uppercase
	if c > 'a' && c < 'z' {
		result = c + shift
		if result > 'z' {
			result -= rune(limit) // loops around
		} else if result < 'a' { // in the case of overflow from negative shifts
			result += rune(limit)
		}
	} else {
		result = c + shift
		if result > 'Z' {
			result -= rune(limit)
		} else if result < 'A' {
			result += rune(limit)
		}
	}
	return result
}

func decipher(c rune, shift rune, limit int) rune {
	var result rune
	if c > 'a' && c < 'z' {
		result = c - shift
		if result < 'a' {
			result += rune(limit) // loops around
		} else if result > 'z' {
			result -= rune(limit)
		}
	} else {
		result = c - shift
		if result < 'A' {
			result += rune(limit)
		} else if result > 'Z' {
			result -= rune(limit)
		}
	}
	return result
}

func Process(msg options.Message, task string, limit int) string {
	var output string
	shift := rune(msg.Shift % limit)     // ensure shift is within limit (26)
	for _, val := range msg.Content { // Go handles range string values as runes (UTF-9)
		if unicode.IsLetter(val) {
			var c rune
			switch task {
			case "cipher":
				c = cipher(val, shift, limit)
			case "decipher":
				c = decipher(val, shift, limit)
			default:
				log.Fatalf("Invalid task: %s", task)
			}
			output += string(c)
		} else {
			output += string(val)
		}
	}
	return output
}
