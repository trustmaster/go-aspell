// Package aspell is simplified bindings to GNU Aspell spell checking library.
package aspell

/*
#cgo LDFLAGS: -laspell
#include <stdlib.h>
#include "aspell.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Speller is a type that encapsulates Aspell internals.
type Speller struct {
	config  *C.AspellConfig
	speller *C.AspellSpeller
}

// NewSpeller creates a new speller instance with configuration options
// given as a map. At least the language option should be specified
// (see example below).
//
// The returned value is a speller struct. The second returned value
// contains error data in case of error or nil if NewSpeller succeeded.
//
// In the most common case you would like to pass the language option
// which accepts two letter ISO 639 language code and an optional
// two letter ISO 3166 country code after a dash or underscore:
//
// 		opts := map[string] string {
// 			"lang": "en_US", // American English
// 		}
// 		speller, err := aspell.NewSpeller(opts)
// 		if err != nil {
// 			panic("Aspell error: " + err.Error())
// 		}
// 		defer speller.Delete()
//
// See available options at http://aspell.net/man-html/The-Options.html
//
// Because aspell package is a binding to Aspell C library, memory
// allocated by NewSpeller() call has to be disposed explicitly.
// This is why the above example contains a deferred call to Delete().
func NewSpeller(options map[string]string) (Speller, error) {
	var s Speller

	// Pass configuration options
	s.config = C.new_aspell_config()
	if _, hasEnc := options["encoding"]; !hasEnc {
		options["encoding"] = "utf-8"
	}
	for k, v := range options {
		optName := C.CString(k)
		optValue := C.CString(v)
		res := C.aspell_config_replace(s.config, optName, optValue)
		C.free(unsafe.Pointer(optName))
		C.free(unsafe.Pointer(optValue))
		if res == 0 {
			msg := C.aspell_config_error_message(s.config)
			err := errors.New(C.GoString(msg))
			C.free(unsafe.Pointer(msg))
			return s, err
		}
	}

	// Attempt to initialize the speller
	var probErr *C.AspellCanHaveError
	probErr = C.new_aspell_speller(s.config)
	C.delete_aspell_config(s.config)
	if C.aspell_error_number(probErr) != 0 {
		msg := C.aspell_error_message(probErr)
		err := errors.New(C.GoString(msg))
		C.free(unsafe.Pointer(msg))
		C.delete_aspell_can_have_error(probErr)
		return s, err
	}

	// Successful speller initialization
	s.speller = C.to_aspell_speller(probErr)
	s.config = C.aspell_speller_config(s.speller)

	return s, nil
}

// SetOption sets Aspell configuration option for the speller.
// All options should be changed before calling spell checking methods.
// See available options at http://aspell.net/man-html/The-Options.html
func (s Speller) Config(name string) string {
	cName := C.CString(name)
	cVal := C.aspell_config_retrieve(s.config, cName)
	val := C.GoString(cVal)
	C.free(unsafe.Pointer(cName))
	C.free(unsafe.Pointer(cVal))
	return val
}

// Check looks the word up in the spell checker dictionary
// and returns true if the word is found there or false
// otherwise.
func (s Speller) Check(word string) bool {
	cword := C.CString(word)
	defer C.free(unsafe.Pointer(cword))
	res := C.aspell_speller_check(s.speller, cword, -1)
	return res != 0
}

// Delete frees memory allocated by Aspell for the speller.
func (s Speller) Delete() {
	// For some reason this breaks everything
	// if s.speller != nil {
	// 	C.delete_aspell_speller(s.speller)
	// }
	// s.config is deleted automatically
}

// Suggest returns a slice of possible suggestions for the given word.
func (s Speller) Suggest(word string) []string {
	cword := C.CString(word)
	defer C.free(unsafe.Pointer(cword))
	suggestions := C.aspell_speller_suggest(s.speller, cword, -1)
	count := uint(C.aspell_word_list_size(suggestions))
	result := make([]string, count)
	elements := C.aspell_word_list_elements(suggestions)
	for i := 0; ; i++ {
		alt := C.aspell_string_enumeration_next(elements)
		if alt == nil {
			break
		}
		result[i] = C.GoString(alt)
		// C.free(unsafe.Pointer(alt))
	}
	C.delete_aspell_string_enumeration(elements)
	return result
}
