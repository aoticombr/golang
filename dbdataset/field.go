package dbdataset

import (
	"database/sql"
	"strings"
	"time"

	"github.com/aoticombr/golang/variant"
)

type Field struct {
	Owner      *Fields
	Name       string
	Caption    string
	DataType   *sql.ColumnType
	IDataType  *DataType
	Precision  int64
	Scale      int64
	DataMask   string
	BoolValue  bool
	TrueValue  string
	FalseValue string
	Visible    bool
	AcceptNull bool
	QuoteNull  bool
	OmitNull   bool
	StrNull    string
	Order      int
	Index      int
}

func NewField(name string) *Field {
	field := &Field{
		Name:       name,
		Caption:    name,
		Precision:  0,
		Scale:      0,
		DataMask:   "",
		BoolValue:  false,
		TrueValue:  "",
		FalseValue: "",
		Visible:    true,
		AcceptNull: true,
		QuoteNull:  true,
		OmitNull:   false,
		StrNull:    "null",
		Order:      1,
		Index:      0,
	}

	return field
}

func (field *Field) AsValue() any {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsValue()
}

func (field *Field) AsString() string {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsString()
}
func (field *Field) AsStringNil() *string {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsStringNil()
}
func (field *Field) AsInt() int {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt()
}
func (field *Field) AsIntNil() *int {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsIntNil()
}
func (field *Field) AsInt8() int8 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt8()
}
func (field *Field) AsInt8Nil() *int8 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt8Nil()
}
func (field *Field) AsInt16() int16 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt16()
}
func (field *Field) AsInt16Nil() *int16 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt16Nil()
}
func (field *Field) AsInt32() int32 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt32()
}
func (field *Field) AsInt32Nil() *int32 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt32Nil()
}
func (field *Field) AsInt64() int64 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt64()
}
func (field *Field) AsInt64Nil() *int64 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsInt64Nil()
}
func (field *Field) AsFloat() float32 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsFloat()
}
func (field *Field) AsFloatNil() *float32 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsFloatNil()
}
func (field *Field) AsFloat64() float64 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsFloat64()
}
func (field *Field) AsFloat64Nil() *float64 {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsFloat64Nil()
}
func (field *Field) AsBool() bool {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsBool()
}
func (field *Field) AsBoolNil() *bool {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsBoolNil()
}
func (field *Field) AsDateTime() time.Time {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsDateTime()
}

func (field *Field) AsDateTimeNil() *time.Time {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsDateTimeNil()
}

func (field *Field) AsByte() []byte {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsByte()
}
func (field *Field) AsByteNil() *[]byte {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).AsByteNil()
}
func (field *Field) IsNull() bool {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).IsNull()
}

func (field *Field) IsNotNull() bool {
	return field.getVariant().SetSilent(field.Owner.Owner.Silent).IsNotNull()
}

func (field *Field) getVariant() *variant.Variant {
	var value *variant.Variant

	if field.Owner != nil {
		if field.Owner.Owner != nil {
			if len(field.Owner.Owner.Rows) > 0 {
				index := field.Owner.Owner.Index
				value = field.Owner.Owner.Rows[index].List[strings.ToUpper(field.Name)]
			}
		}
	}

	if value == nil {
		value = &variant.Variant{}
	}

	return value
}
