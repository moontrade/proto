package wap

func (union *Union) DoLayout() {

}

func (record *Record) DoLayout() {
	compact := false
	if record.Layout() == RecordLayoutCompact {
		record.SetAlign(1)
		compact = true
	} else {
		record.SetAlign(0)
	}
	record.SetOffset(0)
	record.SetOptionals(0)
	record.SetFixed(true)
	fields := record.Fields()
	for i := 0; i < len(fields); i++ {
		field := &fields[i]
		kind := field.Kind().Kind()
		if field.Kind().Enum() != nil {
			kind = field.Kind().Enum().Kind()
		}

		field.SetAlign(0)

		switch kind {
		case KindUnknown:
		case KindBool, KindByte, KindInt8:
			field.Kind().SetSize(1)
			field.Kind().SetAlign(1)
			field.SetPointer(false)
		case KindInt16, KindUInt16:
			field.Kind().SetSize(2)
			field.Kind().SetAlign(2)
			field.SetPointer(false)
		case KindInt32, KindUInt32, KindFloat32:
			field.Kind().SetSize(4)
			field.Kind().SetAlign(4)
			field.SetPointer(false)
		case KindInt64, KindUInt64, KindFloat64:
			field.Kind().SetSize(8)
			field.Kind().SetAlign(8)
			field.SetPointer(false)

		case KindEnum:
			// Ignore since the kind is unwrapped

		case KindStruct, KindRecord:
			r := field.Kind().Record()
			if r != nil {
				field.Kind().SetAlign(r.Align())
				if !r.Fixed() {
					record.SetFixed(false)
					field.Kind().SetAlign(VPointerSize)
					field.Kind().SetSize(VPointerSize)
					field.SetPointer(true)
				} else {
					field.Kind().SetAlign(r.Align())
					field.Kind().SetSize(r.Size())
					field.SetPointer(false)
				}
			}

		case KindStringInline, KindBytesInline:
			field.Kind().SetAlign(1)
			field.SetPointer(false)

		case KindString, KindBytes:
			record.SetFixed(false)
			field.Kind().SetAlign(VPointerSize)
			field.Kind().SetSize(VPointerSize)
			field.SetPointer(true)

		case KindList:
			record.SetFixed(false)
			field.Kind().SetAlign(VPointerSize)
			field.Kind().SetSize(VPointerSize)
			field.SetPointer(true)

		case KindMap:
			record.SetFixed(false)
			field.Kind().SetAlign(VPointerSize)
			field.Kind().SetSize(VPointerSize)
			field.SetPointer(true)

		case KindUnion:
			// TODO: Does union have any pointers?
			field.Kind().SetAlign(field.Kind().Union().Align())
			field.Kind().SetSize(field.Kind().Union().Size())
			field.SetPointer(false)

		case KindPad:
		}

		if field.Optional() && !field.Pointer() {
			record.SetOptionals(record.Optionals() + 1)
		}

		if compact {
			field.SetAlign(1)
		} else {
			field.SetAlign(field.Kind().Align())
			if field.Align() > record.Align() {
				record.SetAlign(field.Align())
			}
		}
	}

	if record.Optionals() > 0 {
		record.SetOptionalOffset(record.Offset())
		record.SetOptionalSize(record.Optionals() / 8)
		if record.Optionals()%8 == 0 {
			record.SetOptionalSize(record.OptionalSize() + 1)
		}
		record.SetOffset(record.OptionalOffset())
	} else {
		record.SetOptionalOffset(0)
		record.SetOptionalSize(0)
	}

	if record.Align() == 0 {
		record.SetAlign(1)
	}

	if record.Offset() > 0 {
		record.SetOffset(align(record.Offset(), record.Align()))
	}

	offset := record.Offset()
	for i := 0; i < len(fields); i++ {
		field := &fields[i]
		offset = align(offset, field.Align())
		field.SetOffset(offset)
		offset += field.Kind().Size()
	}
	record.SetSize(align(offset, record.Align()))
}

func align(offset, align int32) int32 {
	if offset == 0 || align == 0 {
		return 0
	}
	extras := offset % align
	if extras == 0 {
		return offset
	}
	return ((offset / align) + 1) * align
}
