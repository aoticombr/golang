package http

type ContentFormField struct {
	FieldName  string
	FieldValue string
}

type ListContentFormField []*ContentFormField

func (L *ListContentFormField) Add(fieldName string, fieldValue string) {
	*L = append(*L, &ContentFormField{fieldName, fieldValue})
}
func (L *ListContentFormField) Clear() {
	*L = []*ContentFormField{}
}
func NewListContentFormField() ListContentFormField {
	return make(ListContentFormField, 0)
}
