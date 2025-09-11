package esolver

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	r         *bufio.Reader
	elemNames map[string]TokenType
}

func NewScanner(r io.Reader, elemNames map[string]TokenType) *Scanner {
	return &Scanner{
		r:         bufio.NewReader(r),
		elemNames: elemNames,
	}
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

func (s *Scanner) Peek(nRunes int) []rune {
	if nRunes <= 0 {
		return nil
	}
	var runes []rune
	for i := 0; i < nRunes; i++ {
		b, err := s.r.Peek(i + 1)
		if err != nil {
			return runes
		}
		firstByte := b[i]

		var size int
		// Determina o tamanho do rune a partir do primeiro byte.
		if firstByte < utf8.RuneSelf {
			size = 1
		} else if firstByte&0xE0 == 0xC0 {
			size = 2
		} else if firstByte&0xF0 == 0xE0 {
			size = 3
		} else if firstByte&0xF8 == 0xF0 {
			size = 4
		} else {
			return runes
		}
		b, err = s.r.Peek(i + size)
		if err != nil {
			return runes
		}
		r, _ := utf8.DecodeRune(b[i:])
		runes = append(runes, r)
	}
	return runes
}

func (s *Scanner) Scan() Token {
	ch := s.Read()
	if isNumber(ch) {
		s.Unread()
		return s.ScanNumber()
	} else if unicode.IsLetter(ch) {
		s.Unread()
		return s.ScanWord()
	} else if isOperator(ch) {
		return Token{OPERATOR, string(ch)}
	} else if isWhitespace(ch) {
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

	for {
		if ch := s.Read(); ch == eof {
			break
		} else if !unicode.IsLetter(ch) {
			s.Unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	value := buf.String()
	if tt, ok := s.elemNames[value]; ok {
		return Token{tt, value}
	}

	return Token{ERROR, value}
}

func (s *Scanner) ScanNumber() Token {
	var buf bytes.Buffer
	for {
		runes := s.Peek(1)
		if len(runes) == 0 || runes[0] == eof {
			s.Read()
			break
		}
		ch := runes[0]
		if isNumber(ch) {
			s.Read()
			_, _ = buf.WriteRune(ch)
			continue
		}
		if isDegree(ch) {
			nextRunes := s.Peek(2)
			if len(nextRunes) > 1 && unicode.IsLetter(nextRunes[1]) {
				break
			}
			s.Read()
			_, _ = buf.WriteRune(ch)
			continue
		}
		break
	}

	return Token{NUMBER, degToDecString(buf.String())}
}

func (s *Scanner) ScanWhitespace() Token {
	var buf bytes.Buffer
	buf.WriteRune(s.Read())

	for {
		if ch := s.Read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.Unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Token{WHITESPACE, buf.String()}
}

func isNumber(r rune) bool {
	return unicode.IsDigit(r) || r == '.'
}

func isDegree(r rune) bool {
	return r == 'd' || r == '\'' || r == 'm' || r == '"' || r == 's'
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '^'
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}
