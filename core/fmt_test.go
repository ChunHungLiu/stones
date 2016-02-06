package core

import (
	"testing"
)

type vals []interface{}

func TestFmt(t *testing.T) {
	cases := []struct {
		s        string
		args     vals
		expected string
	}{
		// Basic SVO
		{"%s %v %o", vals{"you", "hit", "dog"}, "You hit the dog."},
		{"%s %v %o", vals{"mammoth", "hit", "dog"}, "The mammoth hits the dog."},
		{"%s %v %o", vals{"you", "hit", "Ugh"}, "You hit Ugh."},
		{"%s %v %o", vals{"Ugh", "hit", "you"}, "Ugh hits you."},

		// Embedded verb
		{"%s <hit> %o", vals{"you", "dog"}, "You hit the dog."},
		{"%s <hit> %o", vals{"mammoth", "dog"}, "The mammoth hits the dog."},
		{"%s <hit> %o", vals{"you", "Ugh"}, "You hit Ugh."},
		{"%s <hit> %o", vals{"Ugh", "you"}, "Ugh hits you."},

		// Verb phrases
		{"%s <scream loudly>", vals{"you"}, "You scream loudly."},
		{"%s <scream loudly>", vals{"dog"}, "The dog screams loudly."},

		// Irregular verbs
		{"%s <be cold>", vals{"you"}, "You are cold."},
		{"%s <be cold>", vals{"dog"}, "The dog is cold."},
		{"%s <can eat>", vals{"you"}, "You can eat."},
		{"%s <can cold>", vals{"dog"}, "The dog can cold."},
		{"%s <have> %o", vals{"you", "stick"}, "You have the stick."},
		{"%s <have> %o", vals{"dog", "stick"}, "The dog has the stick."},

		// Reflexive
		{"%s <hit> %o", vals{"you", "you"}, "You hit yourself."},
		{"%s <hit> %o", vals{"dog", "dog"}, "The dog hits itself."},
		{"%s <hit> %o", vals{"Ugh", "Ugh"}, "Ugh hits itself."}, // gender?

		// End punctuation
		{"%s <hit> %o!", vals{"you", "dog"}, "You hit the dog!"},
		{"%s <hit> %o!", vals{"mammoth", "dog"}, "The mammoth hits the dog!"},
		{"%s <hit> %o?", vals{"you", "Ugh"}, "You hit Ugh?"},
		{"%s <hit> %o?", vals{"Ugh", "you"}, "Ugh hits you?"},

		// Literals
		{"%s <hit> %o for %x", vals{"you", "dog", 3}, "You hit the dog for 3."},
		{"%s <hit> %o for %x", vals{"cat", "dog", 3}, "The cat hits the dog for 3."},
		{"%s <hit> %o for %x", vals{"you", "Ugh", 3}, "You hit Ugh for 3."},
		{"%s <hit> %o for %x", vals{"Ugh", "you", 3}, "Ugh hits you for 3."},
	}
	for _, c := range cases {
		if actual := Fmt(c.s, c.args...); actual != c.expected {
			t.Errorf("Fmt(%s, %v) = %s != %s", c.s, c.args, actual, c.expected)
		}
	}
}
