package http

type ContentFile struct {
	Key         string
	FileName    string
	ContentType string
	Content     []byte
}

type ListContentFile []*ContentFile

func (L *ListContentFile) Add(key string, fileName string, contentType string, content []byte) {
	*L = append(*L, &ContentFile{

		Key:         key,
		FileName:    fileName,
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
