package http

type ContentFile struct {
	FieldName   string
	ContentType string
	Content     []byte
}

type ListContentFile []*ContentFile

func (L *ListContentFile) Add(fieldName string, contentType string, content []byte) {
	*L = append(*L, &ContentFile{
		FieldName:   fieldName,
		Content:     content,
		ContentType: contentType,
	})
}
func (L *ListContentFile) Clear() {
	*L = []*ContentFile{}
}
func NewListContentFile() ListContentFile {
	return make(ListContentFile, 0)
}
