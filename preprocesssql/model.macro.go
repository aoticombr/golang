package preprocesssql

type TMacro struct {
	Value any
	Name  string
	SQL   string
}

func NewMacro() *TMacro {
	return &TMacro{}
}

type TMacros struct {
	Items []*TMacro
}

func (ms *TMacros) Clear() {
	ms.Items = []*TMacro{}
}
func (ms *TMacros) Add(value *TMacro) {
	ms.Items = append(ms.Items, value)
}
func (ms *TMacros) Count() int {
	return len(ms.Items)
}
func (ms *TMacros) FindMacro(name string) *TMacro {
	for _, m := range ms.Items {
		if m.Name == name {
			return m
		}
	}
	return nil
}
func (ms *TMacros) NewMacro() *TMacro {
	m := NewMacro()
	ms.Add(m)
	return m
}

func NewMacros() *TMacros {
	return &TMacros{}
}
