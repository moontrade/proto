package runtime

import (
	"fmt"
	"testing"
)

func Test_Schema(t *testing.T) {
	s := &Schema{
		Records: []Record{
			{
				Name: "Bar",
				Fields: []Field{
					Int64Field("id", 0),
					Int64Field("start", 8),
					StringField("name", 16, 4),
					ListField("errors", 24, ListElement(StringElement())),
				},
				FieldsMap: map[string]int{
					"id":    0,
					"start": 1,
				},
			},
		},
		RecordsMap: map[string]int{
			"Bar": 0,
		},
	}
	_ = s
	fmt.Println(s)
}
