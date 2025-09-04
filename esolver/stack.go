package esolver

// Stack is a LIFO data structure
type Stack struct {
	Values []Token
}

// Pop removes the token at the top of the stack and returns its value
func (s *Stack) Pop() Token {
	if len(s.Values) == 0 {
		return Token{}
	}
	token := s.Values[len(s.Values)-1]
	s.Values = s.Values[:len(s.Values)-1]
	return token
}

// Push adds tokens to the top of the stack
func (s *Stack) Push(i ...Token) {
	s.Values = append(s.Values, i...)
}

// Peek returns the token at the top of the stack
func (s *Stack) Peek() Token {
	if len(s.Values) == 0 {
		return Token{}
	}
	return s.Values[len(s.Values)-1]
}

// EmptyInto dumps all tokens from one stack to another
func (s *Stack) EmptyInto(stack *Stack) {
	if !s.IsEmpty() {
		for i := s.Length() - 1; i >= 0; i-- {
			stack.Push(s.Pop())
		}
	}
}

// IsEmpty checks if there are any tokens in the stack
func (s *Stack) IsEmpty() bool {
	return len(s.Values) == 0
}

// Length returns the amount of tokens in the stack
func (s *Stack) Length() int {
	return len(s.Values)
}
