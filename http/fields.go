package http

type Fields map[string][]string

func (f *Fields) Clear() {
	for k := range *f {
		delete(*f, k)
	}
}
func (f *Fields) Add(fieldName string, fieldValue string) {
	if *f == nil {
		*f = make(Fields)
	}
	(*f)[fieldName] = append((*f)[fieldName], fieldValue)
}
func NewFields() Fields {
	return make(Fields)
}
