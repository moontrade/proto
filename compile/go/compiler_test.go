package _go

import (
	"fmt"
	. "github.com/moontrade/proto/schema"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	c, err := NewSchema(&Config{
		Path: "../../pricing",
	})

	if err != nil {
		t.Fatal(err)
	}

	goCompiler, err := NewCompiler(c, &GoConfig{
		BigEndianSafe: false,
		Fluent:        true,
		Mutable:       true,
		Package:       "github.com/moontrade/proto/model",
		Output:        "../../pricing",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = goCompiler.Compile()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c)
}
