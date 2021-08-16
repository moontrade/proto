package schema

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type Builder struct {
	*strings.Builder
}

func NewBuilder() *Builder {
	return &Builder{
		Builder: &strings.Builder{},
	}
}

func (b *Builder) WriteLine(v string) {
	_, _ = b.WriteString(v)
	_, _ = b.WriteString("\n")
}

func (b Builder) W(v string, params ...interface{}) {
	_, _ = b.WriteString(fmt.Sprintf(v, params...))
	_, _ = b.WriteString("\n")
}

// AlignUp rounds n up to a multiple of a. a must be a power of 2.
func AlignUp(n, a uintptr) uintptr {
	return (n + a - 1) &^ (a - 1)
}

// AlignDown rounds n down to a multiple of a. a must be a power of 2.
func AlignDown(n, a uintptr) uintptr {
	return n &^ (a - 1)
}

// DivRoundUp returns ceil(n / a).
func DivRoundUp(n, a uintptr) uintptr {
	// a is generally a power of two. This will get inlined and
	// the compiler will optimize the division.
	return (n + a - 1) / a
}

func PackageName(path string) string {
	return filepath.Base(filepath.Dir(path))
}

func Align(t *Type) int {
	switch t.Kind {
	case KindString:
		return t.Size
	}
	size := t.Size
	if size <= 0 {
		return 1
	}
	switch size {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 4
	case 4:
		return 4
	case 5, 6, 7, 8:
		return 8
	default:
		diff := size % 8
		if diff == 0 {
			return size
		}
		return size + (8 - diff)
	}
}

func FieldAlign(size int) int {
	if size <= 0 {
		return 1
	}
	switch size {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 4
	case 4:
		return 4
	default:
		return 8
	}
}

func SimpleName(str string) string {
	for i := len(str) - 1; i > -1; i-- {
		if str[i] == '.' {
			return str[i+1:]
		}
	}
	return str
}

func StartsWith(val string, s string) bool {
	if len(val) < len(s) {
		return false
	}
	return val[0:len(s)] == s
}

func Join(left, right string) string {
	if len(left) == 0 {
		return right
	}
	if len(right) == 0 {
		return left
	}
	return filepath.Join(left, right)
}

func RelativePath(base, relative string) string {
	dir := base
	if len(filepath.Ext(base)) > 0 {
		dir = filepath.Dir(dir)
	}
	if len(relative) == 0 {
		return base
	}
	// Is absolute path?
	if relative[0] == '/' {
		return relative
	}
	if StartsWith(relative, "./") {
		return Join(dir, relative[2:])
	}
	for StartsWith(relative, "../") {
		if len(dir) < 2 {
			return ""
		}
		dir = filepath.Dir(dir)
		relative = relative[3:]
	}
	return Join(dir, relative)
}

func IsValidName(n string) bool {
	if len(n) == 0 {
		return false
	}
	if !IsValidFirst(int32(n[0])) {
		return false
	}
	for _, c := range n {
		if !IsValidNameCharacter(c) {
			return false
		}
	}
	return true
}

func IsValidFirst(c rune) bool {
	switch c {
	case '_', 'a', 'A', 'b', 'B', 'c', 'C', 'd', 'D', 'e', 'E', 'f', 'F', 'g', 'G', 'h', 'H',
		'i', 'I', 'j', 'J', 'k', 'K', 'l', 'L', 'm', 'M', 'n', 'N', 'o', 'O', 'p', 'P', 'q', 'Q',
		'r', 'R', 's', 'S', 't', 'T', 'u', 'U', 'v', 'V', 'w', 'W', 'x', 'X', 'y', 'Y', 'z', 'Z':
		return true
	default:
		return false
	}
}

func IsValidNameCharacter(c rune) bool {
	switch c {
	case '_', 'a', 'A', 'b', 'B', 'c', 'C', 'd', 'D', 'e', 'E', 'f', 'F', 'g', 'G', 'h', 'H',
		'i', 'I', 'j', 'J', 'k', 'K', 'l', 'L', 'm', 'M', 'n', 'N', 'o', 'O', 'p', 'P', 'q', 'Q',
		'r', 'R', 's', 'S', 't', 'T', 'u', 'U', 'v', 'V', 'w', 'W', 'x', 'X', 'y', 'Y', 'z', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func IsNumeral(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func IsLetter(c byte) bool {
	switch c {
	case 'a', 'A', 'b', 'B', 'c', 'C', 'd', 'D', 'e', 'E', 'f', 'F', 'g', 'G', 'h', 'H',
		'i', 'I', 'j', 'J', 'k', 'K', 'l', 'L', 'm', 'M', 'n', 'N', 'o', 'O', 'p', 'P', 'q', 'Q',
		'r', 'R', 's', 'S', 't', 'T', 'u', 'U', 'v', 'V', 'w', 'W', 'x', 'X', 'y', 'Y', 'z', 'Z':
		return true
	default:
		return false
	}
}

func IsUpper(c byte) bool {
	switch c {
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
		'W', 'X', 'Y', 'Z':
		return true
	default:
		return false
	}
}

func IsLower(c byte) bool {
	switch c {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
		'w', 'x', 'y', 'z':
		return true
	default:
		return false
	}
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if IsLetter(c) {
			if IsUpper(c) {
				return s
			}
			if i > 0 {
				return s[0:i] + strings.ToUpper(s[i:i+1]) + s[i+1:]
			} else {
				return strings.ToUpper(s[i:i+1]) + s[i+1:]
			}
		}
	}
	return s
}

func Uncapitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if IsLetter(c) {
			if IsLower(c) {
				return s
			}
			if i > 0 {
				return s[0:i] + strings.ToLower(s[i:i+1]) + s[i+1:]
			} else {
				return strings.ToLower(s[i:i+1]) + s[i+1:]
			}
		}
	}
	return s
}

func PadEnd(s string, l int) string {
	if len(s) >= l {
		return s
	}
	for i := 0; i < l-len(s); i++ {
		s = s + " "
	}
	return s
}

func IsWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r'
}

func ParseInt(kind Kind, v string) (interface{}, error) {
	switch kind {
	case KindByte, KindUInt16, KindUInt32, KindUInt64:
		return strconv.ParseUint(v, 10, 64)

	case KindInt8, KindInt16, KindInt32, KindInt64:
		return strconv.ParseInt(v, 10, 64)

	default:
		return nil, errors.New("not integer type")
	}
}

//// fnv1 incorporates the list of bytes into the hash x using the FNV-1 hash function.
//func fnv1(x uint32, list ...byte) uint32 {
//	for _, b := range list {
//		x = x*16777619 ^ uint32(b)
//	}
//	return x
//}
//
//func (s *Type) hash() uint64 {
//	d := xxhash.New()
//	buf := make([]byte, 8)
//	buf[0] = byte(s.Kind)
//	_, _ = d.Write(buf[0:1])
//	if s.Optional {
//		buf[0] = 1
//	} else {
//		buf[0] = 0
//	}
//	_, _ = d.Write(buf[0:1])
//	_, _ = d.WriteString(s.Name)
//
//	switch s.Kind {
//	case KindEnum:
//
//	case KindStruct:
//	}
//	return d.Sum64()
//}
