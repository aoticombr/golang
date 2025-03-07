package http

type ContentFormField struct {
	FieldName   string
	FieldValue  string
	ContentType string
}

type ListContentFormField []*ContentFormField

func (L *ListContentFormField) Add(fieldName string, fieldValue string) {
	*L = append(*L, &ContentFormField{FieldName: fieldName, FieldValue: fieldValue, ContentType: ""})
}
func (L *ListContentFormField) AddContentType(fieldName string, fieldValue string, contentType string) {
	*L = append(*L, &ContentFormField{FieldName: fieldName, FieldValue: fieldValue, ContentType: contentType})
}
func (L *ListContentFormField) Clear() {
	*L = []*ContentFormField{}
}
func NewListContentFormField() ListContentFormField {
	return make(ListContentFormField, 0)
}
