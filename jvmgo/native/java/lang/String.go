package lang

import "jvmgo/jvmgo/native"
import "jvmgo/jvmgo/rtda"

const jlString = "java/lang/String"

func init() {
	native.Register(jlString, "intern", "()Ljava/lang/String;", intern)
}

// public native String intern();
// ()Ljava/lang/String;
func intern(frame *rtda.Frame) {
	this := frame.LocalVars().GetThis()
	interned := rtda.InternString(this)
	frame.OperandStack().PushRef(interned)
}
