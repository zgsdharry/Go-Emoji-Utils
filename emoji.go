package emoji

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/tmdvs/Go-Emoji-Utils/utils"
)

// Emoji - Struct representing Emoji
type Emoji struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Descriptor string `json:"descriptor"`
}

// Unmarshal the emoji JSON into the Emojis map
func init() {
	// Work out where we are in relation to the caller
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	// Open the Emoji definition JSON and Unmarshal into map
	jsonFile, err := os.Open(path.Dir(filename) + "/data/emoji.json")
	if jsonFile != nil {
		defer jsonFile.Close()
	}
	if err != nil && len(Emojis) < 1 {
		fmt.Println(err)
	}

	byteValue, e := ioutil.ReadAll(jsonFile)
	if e != nil {
		if len(Emojis) > 0 { // Use build-in emojis data (from emojidata.go)
			return
		}
		panic(e)
	}

	err = json.Unmarshal(byteValue, &Emojis)
	if err != nil {
		panic(e)
	}
}

// LookupEmoji - Lookup a single emoji definition
func LookupEmoji(emojiString string) (emoji Emoji, err error) {

	hexKey := utils.StringToHexKey(emojiString)

	// If we have a definition for this string we'll return it,
	// else we'll return an error
	if e, ok := Emojis[hexKey]; ok {
		emoji = e
	} else {
		err = fmt.Errorf("No record for \"%s\" could be found", emojiString)
	}

	return emoji, err
}

// LookupEmojis - Lookup definitions for each emoji in the input
func LookupEmojis(emoji []string) (matches []interface{}) {
	for _, emoji := range emoji {
		if match, err := LookupEmoji(emoji); err == nil {
			matches = append(matches, match)
		} else {
			matches = append(matches, err)
		}
	}

	return
}

// RemoveAll - Remove all emoji
func RemoveAll(input string) string {

	// Find all the emojis in this string
	matches := FindAll(input)

	// Make a list of the indexes of all the runes used for emoji characters
	emojiRunes := []int{}
	for _, match := range matches {
		for _, loc := range match.Locations {
			for i := loc[0]; i <= loc[1]; i++ {
				emojiRunes = append(emojiRunes, i)
			}
		}
	}

	// Loop over the input strings runes
	runes := []rune(input)
	for i := len(runes); i >= 0; i-- {

		// Loop through the runes indexes used for emoji
		for _, e := range emojiRunes {
			// If the current rune is an emoji rune we'll remove it
			if i == e {
				runes = append(runes[:i-1], runes[i:]...)
			}
		}
	}

	// Remove and trim and left over whitespace
	return strings.TrimSpace(strings.Join(strings.Fields(string(runes)), " "))
	//return input
}
