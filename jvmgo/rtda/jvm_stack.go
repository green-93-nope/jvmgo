package rtda

// This stack is used in jvmgo, it stores the frame with the method of single linked list,
// and the pop frame can be collected by go gc

type Stack struct {
	maxSize uint
	size    uint
	_top    *Frame
}

func newStack(maxSize uint) *Stack {
	return &Stack{
		maxSize: maxSize,
	}
}

func (self *Stack) push(frame *Frame) {
	if self.size >= self.maxSize {
		panic("java.lang.StackOverflowError")
	}
	if self._top != nil {
		frame.lower = self._top
	}
	self._top = frame
	self.size++
}

func (self *Stack) top() *Frame {
	if self._top == nil {
		panic("jvm stack is empty!")
	}
	return self._top
}

func (self *Stack) pop() *Frame {
	top := self.top()
	self._top = top.lower
	top.lower = nil
	self.size--

	return top
}
