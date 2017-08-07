package rtda

func (self *Class) IsArray() bool {
	return self.name[0] == '['
}

func (self *Class) ComponentClass() *Class {
	componentClassName := getComponentClassName(self.name)
	return self.loader.LoadClass(componentClassName)
}

func (self *Class) NewArray(count uint) *Object {
	if !self.IsArray() {
		panic("Not array class: " + self.name)
	}
	switch self.Name() {
	case "[Z":
		return self.NewArrayObject(make([]int8, count))
	case "[B":
		return self.NewArrayObject(make([]int8, count))
	case "[C":
		return self.NewArrayObject(make([]uint16, count))
	case "[S":
		return self.NewArrayObject(make([]int16, count))
	case "[I":
		return self.NewArrayObject(make([]int32, count))
	case "[J":
		return self.NewArrayObject(make([]int64, count))
	case "[F":
		return self.NewArrayObject(make([]float32, count))
	case "[D":
		return self.NewArrayObject(make([]float64, count))
	default:
		return self.NewArrayObject(make([]*Object, count))
	}
}

func NewByteArray(loader *ClassLoader, bytes []int8) *Object {
	return loader.LoadClass("[B").NewArrayObject(bytes)
}
