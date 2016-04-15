/*
See: https://golang.org/doc/codewalk/markov/
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLength int) *Chain {
	return &Chain{
		chain:     make(map[string][]string),
		prefixLen: prefixLength,
	}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	bufReader := bufio.NewReader(r)
	prefix := make(Prefix, c.prefixLen)

	for {
		var word string
		_, err := fmt.Fscan(bufReader, &word)
		if err != nil {
			break
		}

		key := prefix.String()
		c.chain[key] = append(c.chain[key], word)
		prefix.Shift(word)
	}
}

// Generate writes a string of at most n words, generated from Chain,
// to the standard output.
func (c *Chain) Generate(w io.Writer, n int) error {
	bufWriter := bufio.NewWriter(w)
	defer bufWriter.Flush()
	prefix := make(Prefix, c.prefixLen)

	for i := 0; i < n; i++ {
		key := prefix.String()
		words := c.chain[key]
		if len(words) == 0 {
			break
		}

		nextWord := words[rand.Intn(len(words))]
		_, err := bufWriter.WriteString(fmt.Sprint(nextWord, " "))
		if err != nil {
			return err
		}

		prefix.Shift(nextWord)
	}

	return nil
}

func main() {
	numWords := flag.Int("words", 100, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	// Parse flags and seed the random number generator
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// Build up a Markov Chain from the standard input
	chain := NewChain(*prefixLen)
	chain.Build(os.Stdin)

	// Write our generated text to the standard output
	err := chain.Generate(os.Stdout, *numWords)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
