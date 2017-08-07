package main

import (
	"fmt"
	"jvmgo/jvmgo/classpath"
	"jvmgo/jvmgo/instructions/base"
	"jvmgo/jvmgo/rtda"
	"strings"
)

type JVM struct {
	cmd         *Cmd
	classLoader *rtda.ClassLoader
	mainThread  *rtda.Thread
}

func newJVM(cmd *Cmd) *JVM {
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)
	mainThread := rtda.NewThread()
	rtda.InitObjects()
	classLoader := rtda.NewClassLoader(cp, cmd.verboseClassFlag)
	return &JVM{
		cmd:         cmd,
		classLoader: classLoader,
		mainThread:  mainThread,
	}
}

func (self *JVM) start() {
	self.initVM()
	self.execMain()
}

func (self *JVM) initVM() {
	// manually init class of vm
	vmClass := self.classLoader.LoadClass("sun/misc/VM")
	base.InitClass(self.mainThread, vmClass)
	interpret(self.mainThread, self.cmd.verboseInstFlag)
}

func (self *JVM) execMain() {
	className := strings.Replace(self.cmd.class, ".", "/", -1)
	mainClass := self.classLoader.LoadClass(className)
	mainMethod := mainClass.GetMainMethod()
	if mainMethod == nil {
		fmt.Printf("Main method not found in class %s\n", self.cmd.class)
		return
	}

	argsArr := self.createArgsArray()
	frame := self.mainThread.NewFrame(mainMethod)
	frame.LocalVars().SetRef(0, argsArr) // pass arguments to main method
	self.mainThread.PushFrame(frame)
	interpret(self.mainThread, self.cmd.verboseInstFlag)
}

// store strings in the form of java string
func (self *JVM) createArgsArray() *rtda.Object {
	stringClass := self.classLoader.LoadClass("java/lang/String")
	argsLen := uint(len(self.cmd.args))
	argsArr := stringClass.ArrayClass().NewArray(argsLen)
	jArgs := argsArr.Refs()
	for i, arg := range self.cmd.args {
		// get java string in string pool
		jArgs[i] = rtda.JString(self.classLoader, arg)
	}
	return argsArr
}
