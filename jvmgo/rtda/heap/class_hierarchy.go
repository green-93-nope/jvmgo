package heap

func (self *Class) isAssignableFrom(other *Class) bool {
	// in the procedure of assignment t is in the left, s is in the right
	s, t := other, self
	if s == t {
		return true
	}
	if !s.IsArray() {
		if !s.IsInterface() {
			if !t.IsInterface() {
				return s.IsSubClassOf(t)
			} else {
				return s.IsImplements(t)
			}
		} else {
			if !t.IsInterface() {
				return t.IsJlObject()
			} else {
				return t.IsSuperInterfaceOf(s)
			}
		}
	} else {
		if !t.IsArray() {
			if !t.IsInterface() {
				return s.IsJlObject()
			} else {
				return t.IsJlCloneable() || t.IsJioSerializable()
			}
		} else {
			sc := s.ComponentClass()
			tc := t.ComponentClass()
			return sc == tc || tc.isAssignableFrom(sc)
		}
	}
}

func (self *Class) IsSubClassOf(other *Class) bool {
	for c := self.superClass; c != nil; c = c.superClass {
		if c == other {
			return true
		}
	}
	return false
}

func (self *Class) IsSuperClassOf(other *Class) bool {
	return other.IsSubClassOf(self)
}

// whether the self class or its super class implements the iface
func (self *Class) IsImplements(iface *Class) bool {
	for c := self; c != nil; c = c.superClass {
		for _, i := range c.interfaces {
			if i == iface || i.IsSubInterfaceOf(iface) {
				return true
			}
		}
	}
	return false
}

func (self *Class) IsSubInterfaceOf(iface *Class) bool {
	for _, superInterface := range self.interfaces {
		if superInterface == iface || superInterface.IsSubInterfaceOf(iface) {
			return true
		}
	}
	return false
}

func (self *Class) IsSuperInterfaceOf(iface *Class) bool {
	return iface.IsSubInterfaceOf(self)
}
