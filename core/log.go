package core

import (
	"fmt"
	"regexp"
	"strings"
)

// Log applies the Stones log formating language to create a LogMessage.
//
// The format specifiers include the following:
// 	%s - subject
// 	%o - object
// 	%v - verb
// 	%l - literal
// Additionally, verb literals maybe include using the form <verb>.
//
// Each format specifier can be mapped to any arbitrary value, and is converted
// to a string by the fmt package. Consequently, format values should probably
// implement the fmt.Stringer interface to ensure that the log message is
// correctly formated.
//
// Example usage:
// 	Log("%s <hit> %o", hero, orc) yields "You hit the orc."
// 	Log("%s <hit> %o", goblin, hero) yields "The goblin hits you."
// 	Log("%s <hit> %o!", goblin, orc) yields "The goblin hits the orc!"
// 	Log("%s <hit> %o?", goblin, goblin) yields "The goblin hits itself?"
// 	Log("%s <laugh>", unique) yields "Morgoth laughs."
//
// Note that that the hero String conversion is "you" so that the formatter
// knows which grammatical-person to use. Named monsters should have string
// representations which are capitalized so the formatter knows not to add
// certain articles to the names.
//
// Also note that if no ending punctuation is given, then a period is added
// automatically. The sentence is also capitalized if was not already.
func Log(s string, args ...interface{}) string {
	objects := []interface{}{} // subject is always objects[0]

	replace := func(match string) string {
		var noun interface{}

		switch match {
		case "%s":
			noun, args = args[0], args[1:]
			objects = append(objects, noun)
			objects[0] = noun
			return getName(noun)
		case "%o":
			noun, args = args[0], args[1:]
			objects = append(objects, noun)
			if noun == objects[0] {
				return getReflexive(noun)
			}
			return getName(noun)
		case "%v":
			noun, args = args[0], args[1:]
			return getVerb(noun, objects[0])
		case "%l":
			noun, args = args[0], args[1:]
			return fmt.Sprintf("%v", noun)
		}

		return getVerb(match[1:len(match)-1], objects[0])
	}

	return makeSentence(formatRE.ReplaceAllStringFunc(s, replace))
}

var (
	formatRE             = regexp.MustCompile("%s|%o|%v|%l|<.+?>")
	articles             = []string{"the", "a"}
	irregularVerbsSecond = map[string]string{
		"be": "are"}
	irregularVerbsThird = map[string]string{
		"can":  "can",
		"be":   "is",
		"have": "has"}
	esEndings      = []string{"ch", "sh", "ss", "x", "o"}
	endPunctuation = []string{".", "!", "?"}
)

// includesArticle returns true if the given name starts with an article.
func includesArticle(name string) bool {
	for _, article := range articles {
		if strings.HasPrefix(name, article+" ") {
			return true
		}
	}
	return false
}

// getName returns the string name for a particular noun.
func getName(noun interface{}) string {
	name := fmt.Sprintf("%v", noun)
	if name == "you" {
		return "you"
	} else if includesArticle(name) || strings.Title(name) == name {
		return name
	}
	return "the " + name
}

// getReflexive turns a noun into a reflexive pronoun.
func getReflexive(noun interface{}) string {
	name := fmt.Sprintf("%v", noun)
	if name == "you" {
		return "yourself"
	}
	return "itself"
}

// conjuageSecond conjugates a verb in the second person tense.
func conjugateSecond(verb string) string {
	if conjugated, irregular := irregularVerbsSecond[verb]; irregular {
		return conjugated
	}
	return verb
}

// conjugateThird conjugates a verb in the third person tense.
func conjugateThird(verb string) string {
	if congugated, irregular := irregularVerbsThird[verb]; irregular {
		return congugated
	}
	for _, ending := range esEndings {
		if strings.HasSuffix(verb, ending) {
			return verb + "es"
		}
	}
	if strings.HasSuffix(verb, "y") {
		return verb[:len(verb)-1] + "ies"
	}
	return verb + "s"
}

// getVerb conjugates a verb given a particular subject.
func getVerb(verb, subject interface{}) string {
	phrase := strings.Fields(fmt.Sprintf("%v", verb))
	// TODO handle both plural and singular nouns (currently just singular)
	if fmt.Sprintf("%v", subject) == "you" {
		phrase[0] = conjugateSecond(phrase[0])
	} else {
		phrase[0] = conjugateThird(phrase[0])
	}
	return strings.Join(phrase, " ")
}

// makeSentence ensures proper capitolization and punctuation.
func makeSentence(s string) string {
	s = strings.ToUpper(s[:1]) + s[1:]
	for _, punctuation := range endPunctuation {
		if strings.HasSuffix(s, punctuation) {
			return s
		}
	}
	return s + "."
}
