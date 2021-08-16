package schema

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfigFromFS(t *testing.T) {
	config, err := loadFromFS("testdata")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(config)
}

func TestParser(t *testing.T) {
	f, err := os.Open("schema.bu")
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	file, err := ParseFile(f.Name(), f.Name(), bytes)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(file)
}

func TestConst(t *testing.T) {
	file, err := ParseFile("", "", []byte(`
package model

import (
	bu
)

//const HIGH 		i64 = 1000
//const B 		string8 		= "HELLO" 	// comment
//const C 		Code 			= Close
//const B 		[12] string12 			// comment

struct Candle {
	open  		?i64 			
	high  		?Code 			= Open
	low   		?i64 			= nil
	close 		[8] i64
	closeList 	[8] string16

	// Comments above the map field
	//closeMap 	[8] string16 -> string32
}

union Value16 {
	i 	i64
	f 	f64
	d 	epoch
	s 	string16
}

struct Spread {
	low 	i64
	mid 	i64
	high 	i64
}

struct Bar {
	bid 	Candle
	ask 	Candle
	spread 	Spread
	val 	Value16
}

enum Code : byte {
	Open = 0
	Close = 1
}
`))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(file)
}

//func BenchmarkAccess(b *testing.B) {
//	buffer := &BarMut{}
//	rawStruct := (*BarStruct)(unsafe.Pointer(&buffer.Bar[0]))
//
//	rawStruct.bid = 99
//	fmt.Println(buffer.Bid())
//
//	b.Run("Buffer", func(b *testing.B) {
//		b.ReportAllocs()
//		b.ResetTimer()
//
//		for i := 0; i < b.N; i++ {
//			//buffer.Bid()
//			buffer.SetBid(buffer.Bid())
//		}
//	})
//
//	b.Run("Field", func(b *testing.B) {
//		b.ReportAllocs()
//		b.ResetTimer()
//
//		for i := 0; i < b.N; i++ {
//			//rawStruct.Bid()
//			rawStruct.bid = rawStruct.Bid()
//		}
//	})
//}
