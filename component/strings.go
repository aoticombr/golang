package component

import "strings"

type Strings struct {
	itens []string
}

func (s *Strings) Append(value string) *Strings {
	s.itens = append(s.itens, value)
	return s
}
func (s *Strings) Clear() *Strings {
	s.itens = nil
	return s
}
func (s *Strings) Add(value string) *Strings {
	s.itens = append(s.itens, value)
	return s
}
func (s *Strings) Count() int {
	return len(s.itens)
}
func (s *Strings) Text() string {
	return strings.Join(s.itens, " ")
}
func NewStrings() *Strings {
	return &Strings{}
}
