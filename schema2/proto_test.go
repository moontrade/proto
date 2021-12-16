package schema2

import (
	"fmt"
	"github.com/moontrade/nogc"
	"github.com/moontrade/proto"
	"testing"
	"unsafe"
)

func TestProto(t *testing.T) {
	{
		builder := wap.AllocBuilder()
		//builder := wap.NewBuilder()
		c := AllocContact(builder, 128)
		buildContact(&c.ContactMut)
		cc := c.Finish()
		defer cc.Free()
		fmt.Println(cc)
	}
	{
		builder := wap.NewBuilder()
		c := NewContact(builder, 128)
		buildContact(&c.ContactMut)
		fmt.Println(c.Finish())
	}
}

func buildContact(c *ContactMut) {
	cc := c.Unsafe()
	c.SetID(101)
	c.SetDesc("hello")
	desc := c.Desc()
	fmt.Println(desc)
	//n := c.Name()
	c.Name().SetFirst("world")
	cc = c.Unsafe()
	fmt.Println(c.Name().First())
	c.Name().SetNumber(102)
	c.Alt().SetFirst("hi")
	fmt.Println(c.Alt().First())
	cc = c.Unsafe()
	fmt.Println(cc)
}

type ContactRoot struct {
	*Contact
	length_ int
}

type ContactRootPointer struct {
	ContactRoot
}

func (c *ContactRootPointer) Free() {
	if c.Contact != nil {
		nogc.Free(nogc.Pointer(unsafe.Pointer(c.Contact)))
		c.Contact = nil
	}
}

type Contact struct {
	id   int64
	desc wap.VPointer
	_    [4]byte
	name Name
	alt  wap.VPointer
	_    [4]byte
}

type ContactPointer struct {
	*Contact
}

func (m *ContactPointer) Free() {
	if m.Contact != nil {
		nogc.Free(nogc.Pointer(unsafe.Pointer(m.Contact)))
		m.Contact = nil
	}
}

type ContactMut struct {
	//*contactWrapper
	m wap.Mutable
}

type ContactMutRoot struct {
	ContactMut
}

func (c *ContactMutRoot) Finish() ContactRoot {
	return ContactRoot{(*Contact)(c.m.Finish()), c.m.Len()}
}

type ContactMutPointer struct {
	ContactMut
}

func (c *ContactMutPointer) Finish() ContactPointer {
	return ContactPointer{(*Contact)(c.m.FinishPointer())}
}

func AllocContact(b wap.BuilderProvider, flex int32) ContactMutPointer {
	return ContactMutPointer{ContactMut{initContactMut(b.Get().Alloc(int32(unsafe.Sizeof(Contact{})), int32(unsafe.Sizeof(Contact{}))+flex))}}
}

func NewContact(b *wap.Builder, flex int32) ContactMutRoot {
	return ContactMutRoot{ContactMut{initContactMut(b.New(int32(unsafe.Sizeof(Contact{})), int32(unsafe.Sizeof(Contact{}))+flex))}}
}

func initContactMut(m wap.Mutable) wap.Mutable {
	mm := (*Contact)(m.Unsafe())
	initNameMut(m.Slice(unsafe.Pointer(&mm.name)))
	return m
}

func (m *ContactMut) Unsafe() *Contact          { return (*Contact)(m.m.Unsafe()) }
func (self *Contact) ID() int64                 { return self.id }
func (m ContactMut) ID() int64                  { return m.Unsafe().id }
func (m *ContactMut) SetID(value int64)         { m.Unsafe().id = value }
func (self *Contact) Desc() string              { return wap.Str(&self.desc) }
func (m ContactMut) Desc() string               { return wap.Str(&m.Unsafe().desc) }
func (m *ContactMut) SetDesc(value string)      { m.m.WStr(&m.Unsafe().desc, value) }
func (self *Contact) Name() *Name               { return &self.name }
func (m ContactMut) Name() NameMut              { return NameMut{m.m.Slice(unsafe.Pointer(&m.Unsafe().name))} }
func (m ContactMut) SetName(n *Name)            { m.Name().Set(n) }
func (self *Contact) Alt() *Name                { return (*Name)(wap.Slice(&self.alt)) }
func (self *Contact) AltPointer() *wap.VPointer { return &self.alt }
func (self *Contact) IsAltNil() bool            { return self.alt == 0 }
func (m ContactMut) Alt() NameMut {
	return NameMut{m.m.SliceAlloc(&m.Unsafe().alt, int32(unsafe.Sizeof(Name{})))}
}
func (m ContactMut) SetAlt(v *Name) {
	if v == nil {
		m.m.Free(&m.Unsafe().alt)
	} else {
		m.Alt().Set(v)
	}
}
func (m ContactMut) RemoveAlt() {
	m.m.Free(&m.Unsafe().alt)
}

func (self *Contact) Copy(to ContactMut) ContactMut {
	t := to.Unsafe()
	*t = *self
	if self.desc > 0 {
		t.desc = 0
		to.SetDesc(self.Desc())
	}
	to.SetName(self.Name())
	if self.alt != 0 {
		to.SetAlt(self.Alt())
	}
	return to
}

type NameRoot struct {
	*Name
	Size int32
}

type NameRootPointer struct {
	NameRoot
}

func (n *NameRootPointer) Free() {
	if n.Name != nil {
		nogc.Free(nogc.Pointer(unsafe.Pointer(n.Name)))
		n.Name = nil
	}
}

type Name struct {
	first  wap.VPointer
	last   wap.VPointer
	number int64
}

type NameMut struct {
	m wap.Mutable
}

func initNameMut(m wap.Mutable) wap.Mutable {
	//s := (*Name)(unsafe.Pointer(m.Buffer))
	//return n
	return m
}

func (c NameMut) Set(n *Name) {
	if c.m.IsNil() || n == nil {
		return
	}
	c.SetNumber(n.Number())
	c.SetFirst(n.First())
	c.SetLast(n.Last())
}

func (c *Name) CopyTo(to NameMut) NameMut {
	t := to.Unsafe()
	*t = *c
	if c.first > 0 {
		t.first = 0
		to.SetFirst(c.First())
	}
	if c.last > 0 {
		t.last = 0
		to.SetLast(c.Last())
	}
	t.number = c.number
	return to
}

func (n *NameMut) Unsafe() *Name { return (*Name)(n.m.Unsafe()) }
func (c *Name) First() string    { return wap.Str(&c.first) }
func (c NameMut) First() string  { return wap.Str(&c.Unsafe().first) }
func (c NameMut) SetFirst(value string) {
	_ = value
	cast := c.Unsafe()
	c.m.WStr(&cast.first, value)
}
func (c *Name) Last() string            { return wap.Str(&c.last) }
func (c NameMut) Last() string          { return wap.Str(&c.Unsafe().last) }
func (c NameMut) SetLast(value string)  { c.m.WStr(&c.Unsafe().first, value) }
func (c *Name) Number() int64           { return c.number }
func (c NameMut) Number() int64         { return c.Unsafe().number }
func (c NameMut) SetNumber(value int64) { c.Unsafe().number = value }

type ListOfStrings struct {
	length int32
	cap    int32
}

func (l *ListOfStrings) Len() int32 {
	return l.length
}
func (l *ListOfStrings) Get(index int32) string {
	return ""
}

type ListOfStringsMut struct {
	m wap.Mutable
}

type UnionTypeValue struct {
	*UnionType
	selected int
}

type UnionType struct {
	value [16]byte
}

func (u UnionTypeValue) Value() interface{} {
	switch u.selected {
	case 1:
		return (*UnionTypeNumber)(unsafe.Pointer(&u.value[0]))
	default:
		return nil
	}
}

type UnionTypeNumber struct {
	number int64
}

func (u *UnionTypeNumber) Value() int64 {
	return u.number
}

type UnionTypeStruct struct {
}

type UnionTypeRecord struct {
	first wap.VPointer
}
