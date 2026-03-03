package dbdataset

import "fmt"

type DataType int

const (
	Text     DataType = 0
	Integer  DataType = 1
	Float    DataType = 2
	DateTime DataType = 3
	Boolean  DataType = 4
)

func (dt *DataType) IntToDataType(value int) (DataType, error) {
	switch value {
	case 0:
		return Text, nil
	case 1:
		return Integer, nil
	case 2:
		return Float, nil
	case 3:
		return DateTime, nil
	case 4:
		return Boolean, nil
	default:
		return -1, fmt.Errorf("type not found")
	}
}
