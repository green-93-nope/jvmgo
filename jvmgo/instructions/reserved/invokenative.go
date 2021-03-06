package reserved

import (
	"jvmgo/jvmgo/instructions/base"
	"jvmgo/jvmgo/native"
	_ "jvmgo/jvmgo/native/java/io"
	_ "jvmgo/jvmgo/native/java/lang"
	_ "jvmgo/jvmgo/native/java/security"
	_ "jvmgo/jvmgo/native/java/util/concurrent/atomic"
	_ "jvmgo/jvmgo/native/sun/io"
	_ "jvmgo/jvmgo/native/sun/misc"
	_ "jvmgo/jvmgo/native/sun/reflect"
	"jvmgo/jvmgo/rtda"
)

type INVOKE_NATIVE struct{ base.NoOperandsInstruction }

func (self *INVOKE_NATIVE) Execute(frame *rtda.Frame) {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	methodDescriptor := method.Descriptor()

	nativeMethod := native.FindNativeMethod(className, methodName, methodDescriptor)
	if nativeMethod == nil {
		methodInfo := className + "." + methodName + methodDescriptor
		panic("java.lang.UnsatisfiedLinkError" + methodInfo)
	}

	nativeMethod(frame)
}
