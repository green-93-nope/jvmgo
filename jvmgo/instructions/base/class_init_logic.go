package base

import (
	"jvmgo/jvmgo/rtda"
)

// class init is invoked in the following situations:
// new
// putstatic getstatic
// invokestatic
// when the class is initialized and its super class is not initialized
// some reflection operation
func InitClass(thread *rtda.Thread, class *rtda.Class) {
	class.StartInit()
	scheduleClinit(thread, class)
	initSuperClass(thread, class)
}

func scheduleClinit(thread *rtda.Thread, class *rtda.Class) {
	clinit := class.GetClinitMethod()
	if clinit != nil {
		newFrame := thread.NewFrame(clinit)
		thread.PushFrame(newFrame)
	}
}

func initSuperClass(thread *rtda.Thread, class *rtda.Class) {
	if !class.IsInterface() {
		superClass := class.SuperClass()
		if superClass != nil && !superClass.InitStarted() {
			InitClass(thread, superClass)
		}
	}
}
