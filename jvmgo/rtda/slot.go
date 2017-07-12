package rtda

import (
	"jvmgo/jvmgo/rtda/heap"
)

// Slot is the element that stored in LocalVars

type Slot struct {
	num int32
	ref *heap.Object
}
