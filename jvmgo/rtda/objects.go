package rtda

import (
	"jvmgo/jvmgo/rtda/heap"
)

const MAX_OBJECTS = 8

// this stores all objects
type Objects struct {
	mainThread  *Thread
	firstObject *heap.Object
	numObjects  int
	maxObjects  int
}

var objects *Objects

func InitObjects(mainThread *Thread) {
	objects = &Objects{
		mainThread:  mainThread,
		firstObject: nil,
		numObjects:  0,
		maxObjects:  MAX_OBJECTS,
	}
}

func NewObject(class *heap.Class) *heap.Object {
	if objects.numObjects == objects.maxObjects {
		objects.GC()
	}
	object := heap.OneNewObject(class)
	object.SetNext(objects.firstObject)
	object.SetMark(0)
	objects.firstObject = object
	objects.numObjects += 1
	return object
}

func (self *Objects) GC() {
	self.markAll()
	self.sweep()
	self.maxObjects = self.numObjects * 2
}

func (self *Objects) markAll() {
	currentFrame := self.mainThread.CurrentFrame()
	for currentFrame != nil {
		markStack(currentFrame.OperandStack())
		markVars(currentFrame.LocalVars())
		currentFrame = currentFrame.Lower()
	}
}

func markStack(stack *OperandStack) {
	size := stack.GetSize()
	for i := uint(0); i < size; i++ {
		object := stack.GetObject(i)
		markObject(object)
	}
}

func markVars(vars LocalVars) {
	for i := range vars {
		object := vars[i].ref
		markObject(object)
	}
}

func markObject(object *heap.Object) {
	if object == nil {
		return
	}

	object.SetMark(object.Mark() + 1)
	slots := object.Data().([]Slot)
	for i := range slots {
		object := slots[i].ref
		markObject(object)
	}
}

func (self *Objects) sweep() {
	head := &heap.Object{}
	head.SetNext(self.firstObject)
	prev, cur := head, self.firstObject
	for cur != nil {
		if cur.Mark() == 0 {
			cur = cur.Next()
			prev.SetNext(cur)
			self.numObjects -= 1
		} else {
			cur.SetMark(0)
			prev, cur = cur, cur.Next()
		}
	}
	self.firstObject = head.Next()
}
