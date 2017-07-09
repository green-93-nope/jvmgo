package math

import (
	"jvmgo/jvmgo/instructions/base"
	"jvmgo/jvmgo/rtda"
)

// int left shift
type ISHL struct{ base.NoOperandsInstruction }

// int right shift
type ISHR struct{ base.NoOperandsInstruction }

// int right unsigned shift
type IUSHR struct{ base.NoOperandsInstruction }

// long left shift
type LSHL struct{ base.NoOperandsInstruction }

// long right shift
type LSHR struct{ base.NoOperandsInstruction }

// long right unsigned shift
type LUSHR struct{ base.NoOperandsInstruction }

func (self *ISHL) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	// the valid max shift value is 31
	s := uint32(v2) & 0x1f
	result := v1 << s
	stack.PushInt(result)
}

func (self *ISHR) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	s := uint32(v2) & 0x1f
	result := v1 >> s
	stack.PushInt(result)
}

func (self *IUSHR) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	s := uint32(v2) & 0x1f
	result := int32(uint32(v1) >> s)
	stack.PushInt(result)
}

func (self *LSHL) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopLong()
	// the valid max shift value for long is 63
	s := uint32(v2) & 0x3f
	result := v1 << s
	stack.PushLong(result)
}

func (self *LSHR) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopLong()
	s := uint32(v2) & 0x3f
	result := v1 >> s
	stack.PushLong(result)
}
func (self *LUSHR) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopLong()
	s := uint32(v2) & 0x3f
	result := int64(uint64(v1) >> s)
	stack.PushLong(result)
}
