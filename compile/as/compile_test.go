package as

import (
	"fmt"
	. "github.com/moontrade/proto/schema"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	schema, err := NewSchema(&Schema{
		Path: "../../../code/test/assembly",
	})

	if err != nil {
		t.Fatal(err)
	}

	compiler, err := NewCompiler(schema, &ASConfig{
		Mutable: true,
		Output:  "../../../code/test/assembly",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = compiler.Compile()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(schema)
}
