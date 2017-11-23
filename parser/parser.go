package parser

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

// Parser is a quick-n-dirty Makefile "parser", not
// really, just comments and a few directives, but
// you'll forgive me.
type Parser struct {
	i          int
	lines      []string
	commentBuf []string
	target     string
	nodes      []Node
}

// New parser.
func New() *Parser {
	return &Parser{}
}

// Parse the given input reader.
func (p *Parser) Parse(r io.Reader) ([]Node, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "reading")
	}

	p.lines = strings.Split(string(b), "\n")

	if err := p.parse(); err != nil {
		return nil, errors.Wrap(err, "parsing")
	}

	return p.nodes, nil
}

// Peek at the next line.
func (p *Parser) peek() string {
	return p.lines[p.i]
}

// Advance the next line.
func (p *Parser) advance() string {
	s := p.lines[p.i]
	p.i++
	return s
}

// Buffer comment.
func (p *Parser) bufferComment() {
	s := p.advance()[1:]

	if len(s) > 0 {
		if s[0] == '-' {
			return
		}

		// leading space
		if s[0] == ' ' {
			s = s[1:]
		}
	}

	p.commentBuf = append(p.commentBuf, s)
}

// Push comment node.
func (p *Parser) pushComment() {
	if len(p.commentBuf) == 0 {
		return
	}

	s := strings.Join(p.commentBuf, "\n")
	p.nodes = append(p.nodes, Comment{
		Target: p.target,
		Value:  strings.Trim(s, "\n"),
	})

	p.commentBuf = nil
	p.target = ""
}

// Push include node.
func (p *Parser) pushInclude() {
	s := strings.Trim(strings.Replace(p.advance(), "include ", "", 1), " ")
	p.nodes = append(p.nodes, Include{
		Value: s,
	})
}

// Parse the input.
func (p *Parser) parse() error {
	for {
		switch {
		case p.i == len(p.lines)-1:
			return nil
		case strings.HasPrefix(p.peek(), ".PHONY"):
			p.advance()
		case len(p.peek()) == 0:
			p.pushComment()
			p.advance()
		case p.peek()[0] == '#':
			p.bufferComment()
		case strings.HasPrefix(p.peek(), "include "):
			p.pushInclude()
		case strings.ContainsRune(p.peek(), ':'):
			p.target = strings.Split(p.advance(), ":")[0]
			p.pushComment()
		default:
			p.advance()
		}
	}
}
