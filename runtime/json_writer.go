// Package jwriter contains a JSON writer.
package runtime

import (
	"strconv"
	"unicode/utf8"
	"unsafe"
)

// JsonWriterFlags describe various encoding options. The behavior may be actually implemented in the encoder, but
// JsonWriterFlags field in JsonWriter is used to set and pass them around.
type JsonWriterFlags int

const (
	NilMapAsEmpty   JsonWriterFlags = 1 << iota // Encode nil map as '{}' rather than 'null'.
	NilSliceAsEmpty                             // Encode nil slice as '[]' rather than 'null'.
)

// JsonWriter is a JSON writer.
type JsonWriter struct {
	W            Buffer
	Flags        JsonWriterFlags
	Error        error
	NoEscapeHTML bool
}

// Size returns the size of the data that was written out.
func (w *JsonWriter) Size() int {
	return w.W.i
}

// RawByte appends raw binary data to the buffer.
func (w *JsonWriter) RawByte(c byte) {
	w.W.AppendByte(c)
}

// RawString appends raw binary data to the buffer.
func (w *JsonWriter) RawString(s string) {
	w.W.AppendString(s)
}

// Raw appends raw binary data to the buffer or sets the error if it is given. Useful for
// calling with results of MarshalJSON-like functions.
func (w *JsonWriter) Raw(data []byte, err error) {
	switch {
	case w.Error != nil:
		return
	case err != nil:
		w.Error = err
	case len(data) > 0:
		w.W.AppendBytes(data)
	default:
		w.RawString("null")
	}
}

// RawText encloses raw binary data in quotes and appends in to the buffer.
// Useful for calling with results of MarshalText-like functions.
func (w *JsonWriter) RawText(data []byte, err error) {
	switch {
	case w.Error != nil:
		return
	case err != nil:
		w.Error = err
	case len(data) > 0:
		w.String(*(*string)(unsafe.Pointer(&data)))
	default:
		w.RawString("null")
	}
}

// Base64Bytes appends data to the buffer after base64 encoding it
func (w *JsonWriter) Base64Bytes(data []byte) {
	if data == nil {
		w.W.AppendString("null")
		return
	}
	w.W.AppendByte('"')
	w.base64(data)
	w.W.AppendByte('"')
}

func (w *JsonWriter) Uint8(n uint8) {
	w.W.Ensure(3)
	b := strconv.AppendUint(w.W.b[w.W.i:w.W.i], uint64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Uint16(n uint16) {
	w.W.Ensure(5)
	b := strconv.AppendUint(w.W.b[w.W.i:w.W.i], uint64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Uint32(n uint32) {
	w.W.Ensure(10)
	b := strconv.AppendUint(w.W.b[w.W.i:w.W.i], uint64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Uint(n uint) {
	w.W.Ensure(20)
	b := strconv.AppendUint(w.W.b[w.W.i:w.W.i], uint64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Uint64(n uint64) {
	w.W.Ensure(20)
	b := strconv.AppendUint(w.W.b[w.W.i:w.W.i], uint64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Int8(n int8) {
	w.W.Ensure(4)
	b := strconv.AppendInt(w.W.b[w.W.i:w.W.i], int64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Int16(n int16) {
	w.W.Ensure(6)
	b := strconv.AppendInt(w.W.b[w.W.i:w.W.i], int64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Int32(n int32) {
	w.W.Ensure(11)
	b := strconv.AppendInt(w.W.b[w.W.i:w.W.i], int64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Int(n int) {
	w.W.Ensure(21)
	b := strconv.AppendInt(w.W.b[w.W.i:w.W.i], int64(n), 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Int64(n int64) {
	w.W.Ensure(21)
	b := strconv.AppendInt(w.W.b[w.W.i:w.W.i], n, 10)
	w.W.i += len(b)
}

func (w *JsonWriter) Uint8Str(n uint8) {
	w.W.Ensure(5)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendUint(w.W.b[w.W.i+1:w.W.i+1], uint64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Uint16Str(n uint16) {
	w.W.Ensure(7)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendUint(w.W.b[w.W.i+1:w.W.i+1], uint64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Uint32Str(n uint32) {
	w.W.Ensure(12)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendUint(w.W.b[w.W.i+1:w.W.i+1], uint64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) UintStr(n uint) {
	w.W.Ensure(22)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendUint(w.W.b[w.W.i+1:w.W.i+1], uint64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Uint64Str(n uint64) {
	w.W.Ensure(22)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendUint(w.W.b[w.W.i+1:w.W.i+1], uint64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) UintptrStr(n uintptr) {
	w.W.Ensure(22)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendUint(w.W.b[w.W.i+1:w.W.i+1], uint64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Int8Str(n int8) {
	w.W.Ensure(6)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendInt(w.W.b[w.W.i+1:w.W.i+1], int64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Int16Str(n int16) {
	w.W.Ensure(8)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendInt(w.W.b[w.W.i+1:w.W.i+1], int64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Int32Str(n int32) {
	w.W.Ensure(13)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendInt(w.W.b[w.W.i+1:w.W.i+1], int64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) IntStr(n int) {
	w.W.Ensure(23)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendInt(w.W.b[w.W.i+1:w.W.i+1], int64(n), 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Int64Str(n int64) {
	w.W.Ensure(23)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendInt(w.W.b[w.W.i+1:w.W.i+1], n, 10)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Float32(n float32) {
	w.W.Ensure(20)
	b := strconv.AppendFloat(w.W.b[w.W.i:w.W.i], float64(n), 'g', -1, 32)
	w.W.i += len(b)
}

func (w *JsonWriter) Float32Str(n float32) {
	w.W.Ensure(22)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendFloat(w.W.b[w.W.i+1:w.W.i+1], float64(n), 'g', -1, 32)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Float64(n float64) {
	w.W.Ensure(20)
	b := strconv.AppendFloat(w.W.b[w.W.i:w.W.i], n, 'g', -1, 64)
	w.W.i += len(b)
}

func (w *JsonWriter) Float64Str(n float64) {
	w.W.Ensure(22)
	w.W.b[w.W.i] = '"'
	b := strconv.AppendFloat(w.W.b[w.W.i+1:w.W.i+1], n, 'g', -1, 64)
	w.W.b[w.W.i+len(b)+1] = '"'
	w.W.i += 2 + len(b)
}

func (w *JsonWriter) Bool(v bool) {
	if v {
		w.W.AppendString("true")
	} else {
		w.W.AppendString("false")
	}
}

const chars = "0123456789abcdef"

func getTable(falseValues ...int) [128]bool {
	table := [128]bool{}

	for i := 0; i < 128; i++ {
		table[i] = true
	}

	for _, v := range falseValues {
		table[v] = false
	}

	return table
}

var (
	htmlEscapeTable   = getTable(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '&', '<', '>', '\\')
	htmlNoEscapeTable = getTable(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '\\')
)

func (w *JsonWriter) String(s string) {
	w.W.AppendByte('"')

	// Portions of the string that contain no escapes are appended as
	// byte slices.

	p := 0 // last non-escape symbol

	escapeTable := &htmlEscapeTable
	if w.NoEscapeHTML {
		escapeTable = &htmlNoEscapeTable
	}

	for i := 0; i < len(s); {
		c := s[i]

		if c < utf8.RuneSelf {
			if escapeTable[c] {
				// single-width character, no escaping is required
				i++
				continue
			}

			w.W.AppendString(s[p:i])
			switch c {
			case '\t':
				w.W.AppendString(`\t`)
			case '\r':
				w.W.AppendString(`\r`)
			case '\n':
				w.W.AppendString(`\n`)
			case '\\':
				w.W.AppendString(`\\`)
			case '"':
				w.W.AppendString(`\"`)
			default:
				w.W.AppendString(`\u00`)
				w.W.AppendByte(chars[c>>4])
				w.W.AppendByte(chars[c&0xf])
			}

			i++
			p = i
			continue
		}

		// broken utf
		runeValue, runeWidth := utf8.DecodeRuneInString(s[i:])
		if runeValue == utf8.RuneError && runeWidth == 1 {
			w.W.AppendString(s[p:i])
			w.W.AppendString(`\ufffd`)
			i++
			p = i
			continue
		}

		// jsonp stuff - tab separator and line separator
		if runeValue == '\u2028' || runeValue == '\u2029' {
			w.W.AppendString(s[p:i])
			w.W.AppendString(`\u202`)
			w.W.AppendByte(chars[runeValue&0xf])
			i += runeWidth
			p = i
			continue
		}
		i += runeWidth
	}
	w.W.AppendString(s[p:])
	w.W.AppendByte('"')
}

const encode = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
const padChar = '='

func (w *JsonWriter) base64(in []byte) {
	if len(in) == 0 {
		return
	}

	w.W.Ensure(((len(in)-1)/3 + 1) * 4)

	si := 0
	n := (len(in) / 3) * 3

	for si < n {
		// Convert 3x 8bit source bytes into 4 bytes
		val := uint(in[si+0])<<16 | uint(in[si+1])<<8 | uint(in[si+2])

		w.W.b[w.W.i] = encode[val>>18&0x3F]
		w.W.b[w.W.i+1] = encode[val>>12&0x3F]
		w.W.b[w.W.i+2] = encode[val>>6&0x3F]
		w.W.b[w.W.i+3] = encode[val&0x3F]
		w.W.i += 4
		//w.Buffer.Buf = append(w.Buffer.Buf, encode[val>>18&0x3F], encode[val>>12&0x3F], encode[val>>6&0x3F], encode[val&0x3F])

		si += 3
	}

	remain := len(in) - si
	if remain == 0 {
		return
	}

	// Add the remaining small block
	val := uint(in[si+0]) << 16
	if remain == 2 {
		val |= uint(in[si+1]) << 8
	}

	w.W.b[w.W.i] = encode[val>>18&0x3F]
	w.W.b[w.W.i+1] = encode[val>>12&0x3F]
	w.W.i += 2
	//w.Buffer.Buf = append(w.Buffer.Buf, encode[val>>18&0x3F], encode[val>>12&0x3F])

	switch remain {
	case 2:
		w.W.b[w.W.i] = encode[val>>6&0x3F]
		w.W.b[w.W.i+1] = byte(padChar)
		w.W.i += 2
		//w.Buffer.Buf = append(w.Buffer.Buf, encode[val>>6&0x3F], byte(padChar))
	case 1:
		w.W.b[w.W.i] = byte(padChar)
		w.W.b[w.W.i+1] = byte(padChar)
		w.W.i += 2
		//w.Buffer.Buf = append(w.Buffer.Buf, byte(padChar), byte(padChar))
	}
}
