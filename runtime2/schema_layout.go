package runtime2

func (l RecordLayout) Layout(record *Record) {
	switch l {
	case RecordLayoutCompact:
		LayoutCompact(record)
	case RecordLayoutCacheAligned:
		LayoutAligned(record)
	}
}

func LayoutCompact(record *Record) {

}

func LayoutAligned(record *Record) {

}

func (r *Record) Validate() {
	offset := int32(0)
	for _, field := range r.Fields {
		kind := field.Kind
		if field.Enum != nil {
			kind = field.Enum.Kind
		}

		field.Offset = offset
		switch kind {
		case KindUnknown:
		case KindBool, KindByte, KindInt8:
			if field.Optional {
				field.Align = 2
			} else {
				field.Align = 1
			}
			field.Pointer = false
		case KindInt16, KindUInt16:
			if field.Optional {
				field.Align = 3
			} else {
				field.Align = 2
			}
			field.Pointer = false
		case KindInt32, KindUInt32, KindFloat32:
			if field.Optional {
				field.Align = 5
			} else {
				field.Align = 4
			}
			field.Pointer = false
		case KindInt64, KindUInt64, KindFloat64:
			if field.Optional {
				field.Align = 9
			} else {
				field.Align = 8
			}
			field.Pointer = false

		case KindString, KindBytes:
			if field.Optional {
				field.Align = VPointerSize + 1
			} else {
				field.Align = VPointerSize
			}
			field.Pointer = true
			offset += field.Size

		case KindStringFixed, KindFixed:
			if field.Optional {
				field.Align = field.Size + 1
			} else {
				field.Align = field.Size
			}
			field.Pointer = false

		case KindEnum:
			// Ignore since the kind is unwrapped

		case KindStruct, KindRecord:
			record := field.Record
			if record != nil {
				record.Validate()

				if record.IsFlex() {
					if field.Optional {
						field.Align = VPointerSize + 1
					} else {
						field.Align = VPointerSize
					}
					field.Pointer = true
				} else {
					if field.Optional {
						field.Align = record.Size + 1
					} else {
						field.Align = record.Size
					}
					field.Pointer = false
				}
			}

		case KindList:

		case KindMap:

		case KindUnion:

		case KindPad:
			field.Align = field.Size
		}

		offset += field.Align
	}
}

func (r *Record) IsFlex() bool {
	for _, f := range r.Fields {
		if f.Pointer {
			return true
		}
	}
	return false
}
