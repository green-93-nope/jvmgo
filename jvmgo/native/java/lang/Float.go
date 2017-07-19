package lang

import (
	"jvmgo/jvmgo/native"
	"jvmgo/jvmgo/rtda"
	"math"
)

func init() {
	native.Register("java/lang/Float", "floatToRawIntBits", "(F)I", floatToRawIntBits)
	native.Register("java/lang/Float", "intBitsToFloat", "(I)F", intBitsToFloat)
}

// public static native int floatToRawIntBits
func floatToRawIntBits(frame *rtda.Frame) {
	value := frame.LocalVars().GetFloat(0)
	bits := math.Float32bits(value)
	frame.OperandStack().PushInt(int32(bits))
}

// public static native float intBitsToFloat(int bits);
func intBitsToFloat(frame *rtda.Frame) {
	bits := frame.LocalVars().GetInt(0)
	value := math.Float32frombits(uint32(bits))
	frame.OperandStack().PushFloat(value)
}
