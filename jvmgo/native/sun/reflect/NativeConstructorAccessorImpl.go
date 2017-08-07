package reflect

import "jvmgo/jvmgo/instructions/base"
import "jvmgo/jvmgo/native"
import "jvmgo/jvmgo/rtda"

func init() {
	native.Register("sun/reflect/NativeConstructorAccessorImpl", "newInstance0", "(Ljava/lang/reflect/Constructor;[Ljava/lang/Object;)Ljava/lang/Object;", newInstance0)
}

// private static native Object newInstance0(Constructor<?> c, Object[] os)
// throws InstantiationException, IllegalArgumentException, InvocationTargetException;
// (Ljava/lang/reflect/Constructor;[Ljava/lang/Object;)Ljava/lang/Object;
func newInstance0(frame *rtda.Frame) {
	vars := frame.LocalVars()
	constructorObj := vars.GetRef(0)
	argArrObj := vars.GetRef(1)

	goConstructor := getGoConstructor(constructorObj)
	goClass := goConstructor.Class()
	if !goClass.InitStarted() {
		frame.RevertNextPC()
		base.InitClass(frame.Thread(), goClass)
		return
	}

	obj := goClass.NewObject()
	stack := frame.OperandStack()
	stack.PushRef(obj)

	// call <init>
	ops := convertArgs(obj, argArrObj, goConstructor)
	shimFrame := rtda.NewShimFrame(frame.Thread(), ops)
	frame.Thread().PushFrame(shimFrame)

	base.InvokeMethod(shimFrame, goConstructor)
}

func getGoMethod(methodObj *rtda.Object) *rtda.Method {
	return _getGoMethod(methodObj, false)
}
func getGoConstructor(constructorObj *rtda.Object) *rtda.Method {
	return _getGoMethod(constructorObj, true)
}
func _getGoMethod(methodObj *rtda.Object, isConstructor bool) *rtda.Method {
	extra := methodObj.Extra()
	if extra != nil {
		return extra.(*rtda.Method)
	}

	if isConstructor {
		root := methodObj.GetRefVar("root", "Ljava/lang/reflect/Constructor;")
		return root.Extra().(*rtda.Method)
	} else {
		root := methodObj.GetRefVar("root", "Ljava/lang/reflect/Method;")
		return root.Extra().(*rtda.Method)
	}
}

// Object[] -> []interface{}
func convertArgs(this, argArr *rtda.Object, method *rtda.Method) *rtda.OperandStack {
	if method.ArgSlotCount() == 0 {
		return nil
	}

	//	argObjs := argArr.Refs()
	//	argTypes := method.ParsedDescriptor().ParameterTypes()

	ops := rtda.NewOperandStack(method.ArgSlotCount())
	if !method.IsStatic() {
		ops.PushRef(this)
	}
	if method.ArgSlotCount() == 1 && !method.IsStatic() {
		return ops
	}

	//	for i, argType := range argTypes {
	//		argObj := argObjs[i]
	//
	//		if len(argType) == 1 {
	//			// base type
	//			// todo
	//			unboxed := box.Unbox(argObj, argType)
	//			args[i+j] = unboxed
	//			if argType.isLongOrDouble() {
	//				j++
	//			}
	//		} else {
	//			args[i+j] = argObj
	//		}
	//	}

	return ops
}
