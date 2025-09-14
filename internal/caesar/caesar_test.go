package caesar

import (
	"testing"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
)

func TestCipher(t *testing.T) {
	cases := [11]int{0,1,-1,5,-5,10,-10,25,-25,52,-52}
	msg := options.Message{
		Content: "Hello, World. 123!@#",
		Shift: 0,
	}
	limit := 26
	success := [11]string{
		"Hello, World. 123!@#",
		"Ifmmp, Xpsme. 123!@#",
		"Gdkkn, Vnqkc. 123!@#",
		"Mjqqt, Btwqi. 123!@#",
		"Czggj, Rjmgy. 123!@#",
		"Rovvy, Gybvn. 123!@#",
		"Xubbe, Mehbt. 123!@#",
		"Gdkkn, Vnqkc. 123!@#",
		"Ifmmp, Xpsme. 123!@#",
		"Hello, World. 123!@#",
		"Hello, World. 123!@#",
	}

	for i, c := range cases {
		msg.Shift = c
		res := Process(msg, "cipher", limit)
		if res != success[i] {
			t.Errorf("[case %d] Hello, World. 123!@# = %q, want %s, error", i, res, success[i])
		}
	}
}

func TestDecipher(t *testing.T) {
	shifts := [11]int{0,1,-1,5,-5,10,-10,25,-25,52,-52}
	original := "Hello, World. 123!@#"
	msg := options.Message{
		Content: "Hello, World. 123!@#",
		Shift: 0,
	}
	limit := 26
	cases := [11]string{
		"Hello, World. 123!@#",
		"Ifmmp, Xpsme. 123!@#",
		"Gdkkn, Vnqkc. 123!@#",
		"Mjqqt, Btwqi. 123!@#",
		"Czggj, Rjmgy. 123!@#",
		"Rovvy, Gybvn. 123!@#",
		"Xubbe, Mehbt. 123!@#",
		"Gdkkn, Vnqkc. 123!@#",
		"Ifmmp, Xpsme. 123!@#",
		"Hello, World. 123!@#",
		"Hello, World. 123!@#",
	}

	for i, c := range cases {
		msg.Shift = shifts[i]
		msg.Content = c
		res := Process(msg, "decipher", limit)
		if res != original {
			t.Errorf("[case %d] %s = %q, want %s, error", i, c, res, original)
		}
	}
}
