package aspell

import (
	"strings"
	"testing"

//	"fmt"
)

// This is a test for basic Aspell initialization
// and simple word check.
func TestBasic(t *testing.T) {
	// Initialization
	opts := map[string]string{
		"lang": "en_US",
	}
	speller, err := NewSpeller(opts)
	if err != nil {
		t.Errorf("Aspell error: %s", err.Error())
		return
	}
	defer speller.Delete()

	// Test config getter
	enc := speller.Config("encoding")
	if enc != "utf-8" {
		t.Errorf("Invalid aspell encoding: %s", enc)
	}

	// Test Check() against a dictionary
	dict := map[string]bool{
		"bottom":    true,
		"xyzzyw":    false,
		"operation": true,
		"rooby":     false,
		"go":        true,
	}

	for word, res := range dict {
		if speller.Check(word) != res {
			t.Errorf("Incorrect result for '%s', expected: %t", word, res)
		}
	}
}

// This is a test for the list of suggestions
func TestSuggestReplace(t *testing.T) {
	// Initialization
	opts := map[string]string{
		"lang":     "en_US",
		"sug-mode": "slow",
	}
	speller, err := NewSpeller(opts)
	if err != nil {
		t.Errorf("Aspell error: %s", err.Error())
		return
	}
	defer speller.Delete()

	// A "must have" dictionary
	dict := map[string]string{
		"choise":       "choice",
		"soem":         "some",
		"paerticulaur": "particular",
		"unessessay":   "unnecessary",
		"lauf":         "laugh",
		"voteing":      "voting",
		"xiaom":        "Xiaomi",
		// "juse":         "juice", // aspell fails at this
	}
	speller.AddToPersonal("Xiaomi")

	// Search for correct values among suggestions
	for incorrect, correct := range dict {
		suggs := speller.Suggest(incorrect)
		found := false
		for _, word := range suggs {
			if word == correct {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing suggestion for '%s': expected '%s', suggested '%s'", incorrect, correct, strings.Join(suggs, ", "))
		}
	}

	// Store and test a new replacement
	if speller.Replace("juse", "juice") {
		sugJuse := speller.Suggest("juse")
		found := false
		for _, word := range sugJuse {
			if word == "juice" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing replacement for 'juse': expected 'juse', suggested '%s'", strings.Join(sugJuse, ", "))
		}

	} else {
		t.Error("Storing a replacement failed")
	}

	// // Print dict list
	// dicts := Dicts()
	// fmt.Printf("Dicts count: %d\n", len(dicts))
	// for _, dict := range dicts {
	// 	fmt.Printf("Name: %s\nCode: %s\nJargon: %s\nSize: %s\nModule: %s\n\n", dict.name, dict.code, dict.jargon, dict.size, dict.module)
	// }

	// // Print main word list
	// words, err := speller.MainWordList()
	// if err != nil {
	// 	t.Error(err.Error())
	// } else {
	// 	for _, word := range words {
	// 		fmt.Printf("%s, ", word)
	// 	}
	// }
}
