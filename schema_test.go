package wap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test_Schema(t *testing.T) {
	//printLayout(bar{})
	printLayout(r0{})
	//printLayout(r1{})
	//printLayout(r2{})
	//printLayout(r3{})
	//printLayout(r4{})

	s := &Schema{
		types: []Type{
			{
				kind: KindRecord,
				record: &Record{
					name:   "Bar",
					layout: RecordLayoutAligned,
					fields: []Field{
						NewField("a", "i", NumberType(KindInt64)),
						NewField("b", "s", NumberType(KindInt16)),
						NewField("c", "s", NumberType(KindInt32)),
						NewField("d", "s", NumberType(KindInt8)),
						NewField("e", "n", StringType()),
						NewField("f", "e", ListType(StringType())),
					},
				},
			},
		},
	}

	s.Layout()
	//fmt.Println(s)
	//_ = s

	printLayout(bar{})
	printRecordLayout(s.types[0].record)

	//fmt.Println(s)
	//fmt.Println(toJson(s))
}

type bar struct {
	a int64
	b int16
	c int32
	d int8
	e VPointer
	f VPointer
}

type r0 struct {
	a int32
}

type r1 struct {
	a byte
	b int16
}

type r2 struct {
	a byte
	b int32
	c int16
}

type r3 struct {
	a byte
	b int16
	c int32
}

type r4 struct {
	a byte
	b int16
	g int16
	d int32
	c int64
	e [3]byte
	f int16
}

func printRecordLayout(r *Record) {
	fmt.Println("WAP", " name:", r.Name(), " Size:", int(r.Size()), " Align:", int(r.Align()))
	for i := 0; i < len(r.fields); i++ {
		f := &r.fields[i]
		fmt.Println("\t", f.Name(), " -> ", int(f.Offset()), " Size:", int(f.Kind().Size()), " Align:", int(f.Align()))
	}
}

func printLayout(v interface{}) {
	t := reflect.TypeOf(v)
	fmt.Println(t.Name(), " ", "Size:", int(t.Size()), " Align:", int(t.Align()))
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fmt.Println("\t", field.Name, " -> ", int(field.Offset), " Size:", int(field.Type.Size()), " Align:", int(field.Type.Align()))
	}

}

func toJson(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(b)
}
