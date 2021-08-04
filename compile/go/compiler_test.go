package _go

import (
	"fmt"
	. "github.com/moontrade/proto/schema"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	var (
		err      error
		p        *Schema
		compiler *Compiler
	)
	p, err = LoadFromFS("./testdata", true)
	if err != nil {
		t.Fatal(err)
	}

	if compiler, err = NewCompiler(p, &Config{
		BigEndian: false,
		Fluent:    true,
		Mutable:   true,
		Package:   "github.com/moontrade/proto/compile/go/testdata",
		Output:    "./testdata",
	}); err != nil {
		t.Fatal(err)
	}
	if err = compiler.Compile(); err != nil {
		t.Fatal(err)
	}

	fmt.Println(p)
}
