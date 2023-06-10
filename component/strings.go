package component

import "strings"

type stringslist struct {
	itens []string
}

func (s *stringslist) Append(value string) *stringslist {
	s.itens = append(s.itens, value)
	return s
}
func (s *stringslist) Clear() *stringslist {
	s.itens = nil
	return s
}
func (s *stringslist) Add(value string) *stringslist {
	s.itens = append(s.itens, value)
	return s
}
func (s *stringslist) Count() int {
	return len(s.itens)
}
func (s *stringslist) Text() string {
	return strings.Join(s.itens, " ")
}
