package math

import (
	"jvmgo/jvmgo/instructions/base"
	"jvmgo/jvmgo/rtda"
)

// cast type int to x
type I2B struct{ base.NoOperandsInstruction }
type I2C struct{ base.NoOperandsInstruction }
type I2S struct{ base.NoOperandsInstruction }
type I2F struct{ base.NoOperandsInstruction }
type I2D struct{ base.NoOperandsInstruction }
type I2L struct{ base.NoOperandsInstruction }

func (self *I2B) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopInt()
	i := int32(int8(d))
	stack.PushInt(i)
}

func (self *I2C) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopInt()
	i := int32(uint16(d))
	stack.PushInt(i)
}

func (self *I2S) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopInt()
	i := int32(int16(d))
	stack.PushInt(i)
}

func (self *I2D) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopInt()
	i := float64(d)
	stack.PushDouble(i)
}

func (self *I2L) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopInt()
	i := int64(d)
	stack.PushLong(i)
}
func (self *I2F) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopInt()
	i := float32(d)
	stack.PushFloat(i)
}
