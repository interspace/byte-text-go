// Package validate provides routines for validating texts
package validate

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/interspace/byte-text-go/extract"
	"golang.org/x/text/unicode/norm"
)

const (
	invalidChars = "\uFFFE\uFEFF\uFFFF\u202A\u202B\u202C\u202D\u202E"
)

var formC = norm.NFC

// TooLongError is returned when text is too long to be valid.
// The value of the error is the actual length of the input string
type TooLongError struct {
	length    int
	maxLength int
}

func (e TooLongError) Error() string {
	return fmt.Sprintf("Length %d exceeds %d characters", e.length, e.maxLength)
}

// EmptyError is returned when text is empty
type EmptyError struct{}

func (e EmptyError) Error() string {
	return "Text may not be empty"
}

// InvalidCharacterError is returned when text contains an invalid character.
// This error embeds the value of the invalid character, and its byte-offset
// within the input string
type InvalidCharacterError struct {
	Character rune
	Offset    int
}

func (e InvalidCharacterError) Error() string {
	return fmt.Sprintf(
		"Invalid character [%s] found at byte offset %d",
		string(e.Character),
		e.Offset,
	)
}

// TextLength returns the length of the string as it would be displayed. This
// is equivalent to the length of the Unicode NFC
// (See: http://www.unicode.org/reports/tr15). This is needed in order to
// consistently calculate the length of a string no matter which actual form
// was transmitted. For example:
//     U+0065  Latin Small Letter E
// +   U+0301  Combining Acute Accent
// ----------
// =   2 bytes, 2 characters, displayed as é (1 visual glyph)
// … The NFC of {U+0065, U+0301} is {U+00E9}, which is a single character and a
// +display_length+ of 1
// The string could also contain U+00E9 already, in which case the
// canonicalization will not change the value.
func TextLength(text string) int {
	length := utf8.RuneCountInString(formC.String(text))
	return length
}

type ValidationArgs struct {
	maxLength  int
	canBeEmpty bool
}

// TextIsValid checks whether a string is a valid text and returns true or false
func TextIsValid(text string, args ValidationArgs) bool {
	err := TextValidate(text, args)
	return err == nil
}

// TextValidate checks whether a string is a valid text. Returns nil if the
// string is valid. Otherwise, it returns an error in the following cases:
// - The text is too long
// - The text is empty
// - The text contains invalid characters
func TextValidate(text string, args ValidationArgs) error {
	if !args.canBeEmpty && text == "" {
		return EmptyError{}
	} else if length := TextLength(text); length > args.maxLength {
		return TooLongError{length: length, maxLength: args.maxLength}
	} else if i := strings.IndexAny(text, invalidChars); i > -1 {
		r, _ := utf8.DecodeRuneInString(text[i:])
		return InvalidCharacterError{Offset: i, Character: r}
	}
	return nil
}

// UsernameIsValid returns true if the given text represents a valid @username
func UsernameIsValid(username string) bool {
	if username == "" {
		return false
	}

	extracted := extract.MentionedScreenNames(username)
	return len(extracted) == 1 && extracted[0].Text == username
}

// HashtagIsValid returns true if the given text represents a valid #hashtag
func HashtagIsValid(hashtag string) bool {
	if hashtag == "" {
		return false
	}

	extracted := extract.Hashtags(hashtag)
	return len(extracted) == 1 && extracted[0].Text == hashtag
}

// URLIsValid returns true if the given text represents a valid URL
func URLIsValid(url string, requireProtocol bool, allowUnicode bool) bool {
	if url == "" {
		return false
	}

	match := validateURLUnencodedRe.FindStringSubmatchIndex(url)
	if match == nil || url[match[0]:match[1]] != url {
		return false
	}

	if requireProtocol {
		schemeStart := match[validateURLUnencodedGroupScheme*2]
		schemeEnd := match[validateURLUnencodedGroupScheme*2+1]
		if !protocolRe.MatchString(url[schemeStart:schemeEnd]) {
			return false
		}
	}

	pathStart := match[validateURLUnencodedGroupPath*2]
	pathEnd := match[validateURLUnencodedGroupPath*2+1]
	if !validateURLPathRe.MatchString(url[pathStart:pathEnd]) {
		return false
	}

	queryStart := match[validateURLUnencodedGroupQuery*2]
	queryEnd := match[validateURLUnencodedGroupQuery*2+1]
	if queryStart > 0 && !validateURLQueryRe.MatchString(
		url[queryStart:queryEnd],
	) {
		return false
	}

	fragmentStart := match[validateURLUnencodedGroupFragment*2]
	fragmentEnd := match[validateURLUnencodedGroupFragment*2+1]
	if fragmentStart > 0 && !validateURLFragmentRe.MatchString(
		url[fragmentStart:fragmentEnd],
	) {
		return false
	}

	authorityStart := match[validateURLUnencodedGroupAuthority*2]
	authorityEnd := match[validateURLUnencodedGroupAuthority*2+1]
	authority := url[authorityStart:authorityEnd]

	if allowUnicode {
		return validateURLUnicodeAuthorityRe.MatchString(authority)
	}
	return validateURLAuthorityRe.MatchString(authority)
}
