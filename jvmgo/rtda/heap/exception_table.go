package heap

import (
	"jvmgo/jvmgo/classfile"
)

type ExceptionTable []*ExceptionHandler

type ExceptionHandler struct {
	startPc   int
	endPc     int
	handlerPc int
	catchType *ClassRef
}

func newExceptionTable(entries []*classfile.ExceptionTableEntry, cp *ConstantPool) ExceptionTable {
	table := make([]*ExceptionHandler, len(entries))
	for i, entry := range entries {
		table[i] = &ExceptionHandler{
			startPc:   int(entry.StartPc()),
			endPc:     int(entry.EndPc()),
			handlerPc: int(entry.HandlePc()),
			catchType: getCatchType(uint(entry.CatchType()), cp),
		}
	}
	return table
}

func getCatchType(index uint, cp *ConstantPool) *ClassRef {
	if index == 0 { // cp index 0 means catch all
		return nil
	}
	return cp.GetConstant(index).(*ClassRef)
}

func (self ExceptionTable) findExceptionHandler(exClass *Class, pc int) *ExceptionHandler {
	for _, handler := range self {
		if pc >= handler.startPc && pc < handler.endPc {
			// startPc is the next pc of try block, endPc is the next pc out of try block
			if handler.catchType == nil { // catch all type
				return handler
			}
			catchClass := handler.catchType.ResolvedClass()
			if catchClass == exClass || catchClass.IsSuperClassOf(exClass) {
				// if the type of exception Class is the same or sub of catch class
				return handler
			}
		}
	}
	return nil
}
