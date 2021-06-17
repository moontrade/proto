package _go

import (
	"fmt"
	. "github.com/moontrade/proto/schema"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	c, err := NewSchema(&Config{
		Path: "./testdata",
	})

	if err != nil {
		t.Fatal(err)
	}

	goCompiler, err := NewCompiler(c, &GoConfig{
		BigEndian: false,
		Fluent:    true,
		Mutable:   true,
		Package:   "github.com/moontrade/proto/compile/go/testdata/go",
		Output:    "",
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
