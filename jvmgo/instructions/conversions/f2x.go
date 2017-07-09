package math

import (
	"jvmgo/jvmgo/instructions/base"
	"jvmgo/jvmgo/rtda"
)

// cast type float to x
type F2D struct{ base.NoOperandsInstruction }
type F2I struct{ base.NoOperandsInstruction }
type F2L struct{ base.NoOperandsInstruction }

func (self *F2I) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopFloat()
	i := int32(d)
	stack.PushInt(i)
}

func (self *F2L) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopFloat()
	i := int64(d)
	stack.PushLong(i)
}
func (self *F2D) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopFloat()
	i := float64(d)
	stack.PushDouble(i)
}
