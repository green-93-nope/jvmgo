package lang

import "unsafe"
import "jvmgo/jvmgo/rtda"

// []*Class => Class[]
func toClassArr(loader *rtda.ClassLoader, classes []*rtda.Class) *rtda.Object {
	arrLen := len(classes)

	classArrClass := loader.LoadClass("java/lang/Class").ArrayClass()
	classArr := classArrClass.NewArray(uint(arrLen))

	if arrLen > 0 {
		classObjs := classArr.Refs()
		for i, class := range classes {
			classObjs[i] = class.JClass()
		}
	}

	return classArr
}

// []byte => byte[]
func toByteArr(loader *rtda.ClassLoader, goBytes []byte) *rtda.Object {
	if goBytes != nil {
		jBytes := castUint8sToInt8s(goBytes)
		return rtda.NewByteArray(loader, jBytes)
	}
	return nil
}
func castUint8sToInt8s(goBytes []byte) (jBytes []int8) {
	ptr := unsafe.Pointer(&goBytes)
	jBytes = *((*[]int8)(ptr))
	return
}

func getSignatureStr(loader *rtda.ClassLoader, signature string) *rtda.Object {
	if signature != "" {
		return rtda.JString(loader, signature)
	}
	return nil
}
