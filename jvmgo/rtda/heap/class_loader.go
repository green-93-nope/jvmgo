package heap

import "fmt"
import "jvmgo/jvmgo/classfile"
import "jvmgo/jvmgo/classpath"

/*
class names:
    - primitive types: boolean, byte, int ...
    - primitive arrays: [Z, [B, [I ...
    - non-array classes: java/lang/Object ...
    - array classes: [Ljava/lang/Object; ...
*/
type ClassLoader struct {
	verboseFlag bool
	parent      *ClassLoader
	readFun     classpath.ReadFunction
	classMap    map[string]*Class // loaded classes
}

var bootstrap_classloader *ClassLoader
var extension_classloader *ClassLoader
var application_classloader *ClassLoader

func NewClassLoader(cp *classpath.Classpath, verboseFlag bool) *ClassLoader {
	bootstrap_classloader = &ClassLoader{
		verboseFlag: verboseFlag,
		parent:      nil,
		readFun:     cp.ReadBootClass,
		classMap:    make(map[string]*Class),
	}

	extension_classloader = &ClassLoader{
		verboseFlag: verboseFlag,
		parent:      bootstrap_classloader,
		readFun:     cp.ReadExtClass,
		classMap:    make(map[string]*Class),
	}

	application_classloader = &ClassLoader{
		verboseFlag: verboseFlag,
		parent:      extension_classloader,
		readFun:     cp.ReadUserClass,
		classMap:    make(map[string]*Class),
	}

	bootstrap_classloader.loadBasicClasses()
	bootstrap_classloader.loadPrimitiveClasses()
	return application_classloader
}

func (self *ClassLoader) loadBasicClasses() {
	jlClassClass := self.LoadClass("java/lang/Class")
	// set the jClass and its extra of the class when first load "java/lang/Class"
	for _, class := range self.classMap {
		if class.jClass == nil {
			// jClass is the instance of class struct
			class.jClass = jlClassClass.NewObject()
			// extra is the instance of which class
			class.jClass.extra = class
		}
	}
}

// load classes of basic type
func (self *ClassLoader) loadPrimitiveClasses() {
	for primitiveType, _ := range primitiveTypes {
		self.loadPrimitiveClass(primitiveType)
	}
}

func (self *ClassLoader) loadPrimitiveClass(className string) {
	class := &Class{
		accessFlags: ACC_PUBLIC, // todo
		name:        className,
		loader:      self,
		initStarted: true,
	}
	// jClass is the object of type "java/lang/Class"
	class.jClass = self.classMap["java/lang/Class"].NewObject()
	class.jClass.extra = class
	self.classMap[className] = class
}

func (self *ClassLoader) LoadClass(name string) *Class {
	if c := self.findLoadedClass(name); c != nil {
		return c
	}

	if self.parent != nil {
		if c := self.parent.LoadClass(name); c != nil {
			return c
		}
	}
	c := self.findClass(name)
	return c
}

func (self *ClassLoader) findLoadedClass(name string) *Class {
	if self.parent != nil {
		if c := self.parent.findLoadedClass(name); c != nil {
			return c
		}
	}

	for key, v := range self.classMap {
		if key == name {
			return v
		}
	}
	return nil
}

// to use user defined classloader, this class should be override
func (self *ClassLoader) findClass(name string) *Class {
	var class *Class
	if name[0] == '[' { // array class
		class = self.loadArrayClass(name)
	} else {
		class = self.loadNonArrayClass(name)
	}

	if class == nil {
		// can not load
		return nil
	}

	if jlClassClass := self.findLoadedClass("java/lang/Class"); jlClassClass != nil {
		class.jClass = jlClassClass.NewObject()
		class.jClass.extra = class
	}

	return class
}

func (self *ClassLoader) loadArrayClass(name string) *Class {
	class := &Class{
		accessFlags: ACC_PUBLIC, // todo
		name:        name,
		loader:      self,
		initStarted: true,
		superClass:  self.LoadClass("java/lang/Object"),
		interfaces: []*Class{
			self.LoadClass("java/lang/Cloneable"),
			self.LoadClass("java/io/Serializable"),
		},
	}
	self.classMap[name] = class
	return class
}

func (self *ClassLoader) loadNonArrayClass(name string) *Class {
	data, entry, err := self.readFun(name)

	if err != nil {
		return nil
	}

	class := self.defineClass(data)
	link(class)

	if self.verboseFlag {
		fmt.Printf("[Loaded %s from %s]\n", name, entry)
	}

	return class
}

// jvms 5.3.5
func (self *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)
	hackClass(class)
	class.loader = self
	resolveSuperClass(class)
	resolveInterfaces(class)
	self.classMap[class.name] = class
	return class
}

func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data)
	if err != nil {
		//panic("java.lang.ClassFormatError")
		panic(err)
	}
	return newClass(cf)
}

// jvms 5.4.3.1
func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}
func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}

func link(class *Class) {
	verify(class)
	prepare(class)
}

func verify(class *Class) {
	// todo
}

// jvms 5.4.2
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)
	calcStaticFieldSlotIds(class)
	allocAndInitStaticVars(class)
}

func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.instanceSlotCount = slotId
}

func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}

func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount)
	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {
			initStaticFinalVar(class, field)
		}
	}
}

func initStaticFinalVar(class *Class, field *Field) {
	vars := class.staticVars
	cp := class.constantPool
	cpIndex := field.ConstValueIndex()
	slotId := field.SlotId()

	if cpIndex > 0 {
		switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotId, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotId, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotId, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotId, val)
		case "Ljava/lang/String;":
			goStr := cp.GetConstant(cpIndex).(string)
			jStr := JString(class.Loader(), goStr)
			vars.SetRef(slotId, jStr)
		}
	}
}

// todo
func hackClass(class *Class) {
	if class.name == "java/lang/ClassLoader" {
		loadLibrary := class.GetStaticMethod("loadLibrary", "(Ljava/lang/Class;Ljava/lang/String;Z)V")
		loadLibrary.code = []byte{0xb1} // return void
	}
}
