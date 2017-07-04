package classpath

import (
	"io/ioutil"
	"path/filepath"
)

// represent the classpath of dictionary
type DirEntry struct {
	absDir string
}

// newDirEntry absolute dir entry
func newDirEntry(path string) *DirEntry {
	absDir, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &DirEntry{absDir}
}

func (self *DirEntry) readClass(className string) ([]byte, Entry, error) {
	fileName := filepath.Join(self.absDir, className)
	data, error := ioutil.ReadFile(fileName)
	return data, self, error
}

func (self *DirEntry) String() string {
	return self.absDir
}
