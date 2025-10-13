package dbdataset

import (
	"fmt"
	"strings"
)

type Fields struct {
	Owner *DataSet
	List  []*Field
}

func NewFields() *Fields {
	fields := &Fields{
		List: []*Field{},
	}
	return fields
}

func (fd *Fields) FindFieldByName(fieldName string) *Field {
	for i := 0; i < len(fd.List); i++ {
		if strings.EqualFold(fd.List[i].Name, fieldName) {
			return fd.List[i]
		}
	}
	return nil
}

func (fd *Fields) FieldByName(fieldName string) *Field {
	var field *Field

	for i := 0; i < len(fd.List); i++ {
		if strings.EqualFold(fd.List[i].Name, fieldName) {
			field = fd.List[i]
			return field
		}
	}

	if field == nil {
		field = &Field{Owner: fd}
		fmt.Println("Field " + fieldName + " doesn't exists")
	}

	return field
}

func (fd *Fields) Add(fieldName string) *Field {
	field := fd.FindFieldByName(fieldName)

	if field != nil {
		return field
	} else {
		field = NewField(fieldName)
		field.Owner = fd
		fd.List = append(fd.List, field)
		return field
	}
}

func (fd *Fields) Clear() *Fields {
	fd.List = nil
	return fd
}

func (fd *Fields) Count() int {
	return len(fd.List)
}
