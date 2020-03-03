package calc

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) Unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) Scan() Token {
	ch := s.Read()
	if IsNumber(ch) {
		s.Unread()
		return s.ScanNumber()
	} else if unicode.IsLetter(ch) {
		s.Unread()
		return s.ScanWord()
	} else if IsOperator(ch) {
		return Token{OPERATOR, string(ch)}
	} else if IsWhitespace(ch) {
		s.Unread()
		return s.ScanWhitespace()
	}

	switch ch {
	case eof:
		return Token{EOF, ""}
	case '(':
		return Token{LPAREN, "("}
	case ')':
		return Token{RPAREN, ")"}
	}

	return Token{ERROR, string(ch)}
}

func (s *Scanner) ScanWord() Token {
	var buf bytes.Buffer
	buf.WriteRune(s.Read())
	value := buf.String()

	for v, tokenType := range elemNames {
		for {
			if value == v {
				return Token{tokenType, value}
			}

			if strings.HasPrefix(v, value) {
				ch := s.Read()
				if ch == eof {
					break
				}

				buf.WriteRune(ch)
				value = buf.String()
			} else {
				break
			}
		}
	}

	return Token{CONSTANT, value}
}

func (s *Scanner) ScanNumber() Token {
	var buf bytes.Buffer
	buf.WriteRune(s.Read())

	for {
		if ch := s.Read(); ch == eof {
			break
		} else if !IsNumber(ch) && !IsDegree(ch) {
			s.Unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return Token{NUMBER, degToDecString(buf.String())}
}

func (s *Scanner) ScanWhitespace() Token {
	var buf bytes.Buffer
	buf.WriteRune(s.Read())

	for {
		if ch := s.Read(); ch == eof {
			break
		} else if !IsWhitespace(ch) {
			s.Unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Token{WHITESPACE, buf.String()}
}

func IsNumber(r rune) bool {
	return unicode.IsDigit(r) || r == '.'
}

func IsDegree(r rune) bool {
	return r == 'd' || r == '\'' || r == 'm' || r == '"' || r == 's'
}

func IsOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '^'
}

func IsWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}
