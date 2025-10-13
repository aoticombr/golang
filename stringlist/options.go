package stringlist

type Options func(*Strings)

func WithDelimiter(delimiter string) Options {
	return func(s *Strings) {
		s.Delimiter = delimiter
	}
}
