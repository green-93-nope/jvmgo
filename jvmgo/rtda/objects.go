package rtda

const MAX_OBJECTS = 8

// this stores all objects
type Objects struct {
	firstObject *Object
	numObjects  int
	maxObjects  int
}

var objects *Objects

func InitObjects() {
	objects = &Objects{
		firstObject: nil,
		numObjects:  0,
		maxObjects:  MAX_OBJECTS,
	}
}

func GetNewObject(object *Object) *Object {
	if objects.numObjects == objects.maxObjects {
		objects.GC()
	}
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
	cur := objects.firstObject
	for cur != nil {
		markObjectSlots(cur)
		cur = cur.Next()
	}
}

func markObjectSlots(object *Object) {
	slots := object.Data().(Slots)
	for i := range slots {
		object := slots[i].ref
		markObject(object)
	}
}

func markObject(object *Object) {
	if object == nil {
		return
	}

	object.SetMark(object.Mark() + 1)
	markObjectSlots(object)
}

func (self *Objects) sweep() {
	head := &Object{}
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
