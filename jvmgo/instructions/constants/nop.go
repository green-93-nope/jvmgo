package constants

import (
	"jvmgo/jvmgo/instructions/base"
	"jvmgo/jvmgo/rtda"
)

type NOP struct{ base.NoOperandsInstruction }

func (self *NOP) Execute(frame *rtda.Frame) {

}
