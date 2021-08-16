package runtime

func (r *Record) InitialSize(hint int) int {
	if !r.Flex {
		return int(r.Size)
	}
	if hint > int(r.Size) {
		return hint
	}
	return int(r.Size)
}

type JsonMarshaller struct {
	w JsonWriter
}

func NewJsonMarshaller(size int) *JsonMarshaller {
	return &JsonMarshaller{w: JsonWriter{
		W: Buffer{
			b: GetBytes(size),
			i: 0,
		},
	}}
}

type JsonUnmarshaller struct {
	l JsonLexer
}

func (r *Record) NewJsonUnmarshaller(b []byte) *JsonUnmarshaller {
	return &JsonUnmarshaller{JsonLexer{
		Data: b,
	}}
}

func (r *Record) MarshalJSON(p Pointer, w *JsonMarshaller) error {
	for _, field := range r.Fields {
		_ = field
	}
	return nil
}

func (r *Record) UnmarshalJSON(rd *JsonUnmarshaller) error {
	return nil
}
