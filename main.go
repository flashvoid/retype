package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bendahl/uinput"
)

type keyPress struct {
	code int
	mode int
}

func NewKeyPress(code int, mode int) keyPress {
	return keyPress{code, mode}
}

func (k keyPress) Press(keyboard uinput.Keyboard) {
	if k.code == 0 {
		return
	}

	if k.mode != 0 {
		keyboard.KeyDown(k.mode)
	}

	keyboard.KeyPress(k.code)

	if k.mode != 0 {
		keyboard.KeyUp(k.mode)
	}
}
func convertToLetter(word string) string {
	var letter string

	switch word {
	case "exclam":
		letter = "!"
	case "at":
		letter = "@"
	case "numbersign":
		letter = "#"
	case "dollar":
		letter = "$"
	case "percent":
		letter = "%"
	case "asciicircum":
		letter = "^"
	case "ampersand":
		letter = "&"
	case "asterisk":
		letter = "*"
	case "parenleft":
		letter = "("
	case "parenright":
		letter = ")"
	case "minus":
		letter = "-"
	case "underscore":
		letter = "_"
	case "equal":
		letter = "="
	case "plus":
		letter = "+"
	case "bracketleft":
		letter = "["
	case "braceleft":
		letter = "{"
	case "bracketright":
		letter = "]"
	case "braceright":
		letter = "}"
	case "semicolon":
		letter = ";"
	case "colon":
		letter = ":"
	case "apostrophe":
		letter = "'"
	case "quotedbl":
		letter = string('"')
	case "grave":
		letter = "`"
	case "asciitilde":
		letter = "~"
	case "backslash":
		letter = "\\"
	case "bar":
		letter = "|"
	case "comma":
		letter = ","
	case "less":
		letter = "<"
	case "period":
		letter = "."
	case "greater":
		letter = ">"
	case "slash":
		letter = "/"
	case "question":
		letter = "?"
	case "space":
		letter = " "
	case "Return":
		letter = "\n"
	}

	return letter
}

func loadKeymaps(file string) (map[string]keyPress, error) {
	keymap := make(map[string]keyPress)
	var data []byte
	var err error

	if file == "" {
		data = DefaultKeyMap
	} else {
		data, err = ioutil.ReadFile(file)
	}
	if err != nil {
		return keymap, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") || strings.HasPrefix(text, " ") {
			continue
		}

		if !strings.Contains(text, "0x") {
			continue
		}

		// here `text` is expected to have 1-3 words
		// 1 letter or name of key
		// 2 hex code of key
		// 3 optional modifier

		textA := strings.Split(text, " ")
		if len(textA) < 2 || len(textA) > 3 {
			continue
		}

		letter := textA[0]
		if len(letter) > 1 {
			letter = convertToLetter(letter)
			if letter == "" {
				continue
			}

		}

		// keyCode format expected to be 0xXX
		if len(textA[1]) < 3 {
			continue
		}
		keyCode, err := strconv.ParseInt(textA[1][2:], 16, 32)
		if err != nil {
			continue
		}

		var mode int = 0
		if len(textA) == 3 {
			mode = uinput.KeyLeftshift
		}

		keymap[letter] = NewKeyPress(int(keyCode), mode)

	}

	return keymap, err
}

func main() {
	flagDelay := flag.Int("delay", 0, "delay typing (seconds)")
	flagInterval := flag.Int("interval", 1, "interval between keypresses (milliseconds)")
	flagKeymapFile := flag.String("keymap", "", "use alternative keymap file (or set KEYMAPS variable)")
	flagDumpKeymap := flag.Bool("dump-keymap", false, "dumps embedded keymap to stdout (use as template to craft your own")
	flag.Parse()

	if *flagDumpKeymap {
		os.Stdout.Write(DefaultKeyMap)
		os.Exit(0)
	}

	keyMapsFile := os.Getenv("KEYMAPS")
	if *flagKeymapFile != "" {
		keyMapsFile = *flagKeymapFile
	}

	var err error
	var chars map[string]keyPress
	chars, err = loadKeymaps(keyMapsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "retype failed to load keymaps file, err=%s\n", err)
		os.Exit(1)
	}

	if *flagDelay != 0 {
		time.Sleep(time.Second * time.Duration(*flagDelay))
	}

	// initialize keyboard and check for possible errors
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("retype-keyboard"))
	if err != nil {
		fmt.Printf("retype failed to open /dev/uinput: %s", err)
		os.Exit(1)
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer keyboard.Close()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		chr := scanner.Text()
		key, ok := chars[chr]
		if !ok {
			continue
		}

		key.Press(keyboard)
		time.Sleep(time.Millisecond * time.Duration(*flagInterval))
	}
}
