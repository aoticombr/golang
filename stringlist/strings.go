package stringlist

import "strings"

type Strings struct {
	Delimiter string
	Items     []string
}

func (s *Strings) Append(value string) *Strings {
	s.Items = append(s.Items, value)
	return s
}
func (s *Strings) Clear() *Strings {
	s.Items = nil
	return s
}
func (s *Strings) Add(value string) *Strings {
	s.Items = append(s.Items, value)
	return s
}
func (s *Strings) AddStrings(value *Strings) {
	for i := 0; i < value.Count(); i++ {
		s.Add(value.Items[i])
	}
}
func (s *Strings) Count() int {
	return len(s.Items)
}
func (s *Strings) Text() string {
	return strings.Join(s.Items, s.Delimiter)
}

func (s *Strings) Byte() []byte {
	return []byte(s.Text())
}

func NewStrings(opt ...Options) *Strings {
	s := &Strings{
		Delimiter: " ",
	}
	for _, o := range opt {
		o(s)
	}
	return s
}
