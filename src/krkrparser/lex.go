package krkrparser

import (
	"errors"
	"fmt"
	"strings"
)

const (
	INIT = iota
	TEXT
	COMMAND
	COMMA
	SPACE
	SCENE
)

type Parser struct {
	data []byte
	pos  int
}

func (p *Parser) Init(data []byte) {
	p.data = data
	p.pos = 0
}

func (p *Parser) Read(s byte) (string, error) {
	r := make([]byte, 0)
	for ; p.data[p.pos] != s; p.pos++ {
		if p.pos >= len(p.data) {
			return string(r), errors.New("over ranged")
		}
		r = append(r, p.data[p.pos])
	}
	r = append(r, p.data[p.pos])
	p.pos++
	return string(r), nil
}

func (p *Parser) format(s string) string {
	r := ""
	for i := 0; i < len(s); {
		switch s[i] {
		case '[':
			cc := make([]byte, 0)
			for ; s[i] != ']' || s[i-1] == '"'; i++ {
				cc = append(cc, s[i])
			}
			cc = append(cc,s[i])
			i++
			c := string(cc[1:len(cc)-1])
			switch {
			case strings.Contains(c, "'"):
				arr := strings.Split(c, "'")
				if len(arr) == 2 {
					r += fmt.Sprintf("%s(%s)", arr[0], arr[1])
				}
			case len(strings.Split(c, " ")) > 1:
				arr := strings.Split(c, " ")
				if strings.HasPrefix(strings.ToUpper(arr[0]), "CH") {
					rbs := ""
					fmt.Sscanf(arr[1], "text=%s", &rbs)
					r += strings.Replace(rbs, "\"", "", -1)
					for j := 2; j < len(arr); j++ {
						r += arr[j]
					}
				}
			default:
				r += c
			}
		default:
			r = string(append([]byte(r), s[i]))
			i++
		}
	}
	return r
}

func (p *Parser) Next() (int, string, error) {
	if p.pos >= len(p.data) {
		return SPACE, "", errors.New("eof")
	}
	switch p.data[p.pos] {
	case '\n', '\r':
		s, err := p.Read('\n')
		if len(s) == 0 && err != nil {
			return SPACE, s, err
		}
		return SPACE, s, nil
	case '[':
		s, err := p.Read('\n')
		if len(s) == 0 && err != nil {
			return COMMAND, s, err
		}
		return COMMAND, s, nil
	case ';':
		s, err := p.Read('\n')
		if len(s) == 0 && err != nil {
			return COMMA, s, err
		}
		return COMMA, s, nil
	case '*':
		s, err := p.Read('\n')
		if len(s) == 0 && err != nil {
			return SCENE, s, err
		}
		return SCENE, s, nil
	default:
		s, err := p.Read('\n')
		if len(s) == 0 && err != nil {
			return COMMA, s, err
		}
		return TEXT, p.format(s), nil
	}
	return INIT, "", nil
}

func NewParser() *Parser {
	return new(Parser)
}
