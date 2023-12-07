package file

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/huyungtang/go-lib/strings"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	NotExist statusEnum = iota + 1
	IsFile
	IsDir
)

const (
	Create FileOp = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Dir
// ****************************************************************************************************************************************
func Dir(path string) string {
	return filepath.Dir(path)
}

// IsExist
// ****************************************************************************************************************************************
func IsExist(path string) statusEnum {
	f, e := os.Stat(path)
	if e != nil && os.IsNotExist(e) {
		return NotExist
	} else if f.IsDir() {
		return IsDir
	}

	return IsFile
}

// MakeDir
// ****************************************************************************************************************************************
func MakeDir(path string) (err error) {
	switch IsExist(path) {
	case IsDir:
		return
	case IsFile:
		return errors.New("target is pointed to a file")
	}

	return os.MkdirAll(path, os.ModePerm)
}

// Path
// ****************************************************************************************************************************************
func Path(root string, dirs ...string) string {
	dirs = append([]string{filepath.Dir(root)}, dirs...)

	return filepath.Join(dirs...)
}

// PathCurrent
// ****************************************************************************************************************************************
func PathCurrent(dirs ...string) string {
	return Path(filepath.Dir(os.Args[0]), dirs...)
}

// GetFilename
// ****************************************************************************************************************************************
func GetFilename(path string) string {
	return filepath.Base(path)
}

// PathHome
// ****************************************************************************************************************************************
func PathHome(dirs ...string) string {
	home, _ := os.UserHomeDir()

	return Path(home, dirs...)
}

// PathSavename
// ****************************************************************************************************************************************
func PathSavename(root string, layer int) string {
	if layer > 5 {
		layer = 5
	}
	dl := layer * 2

	fn, dn := []rune(strings.ToLower(strings.Replace(uuid.New().String(), "-", ""))), make([]rune, dl)
	copy(dn, fn)

	dir := make([]string, 0, layer)
	for i := 0; i < dl; i += 2 {
		dir = append(dir, string(dn[i:i+2]))
	}

	root = Path(strings.Format("%s/", root), dir...)
	MakeDir(root)

	return filepath.Join(root, string(fn))
}

// GetSavename
// ****************************************************************************************************************************************
func GetSavename(root, fn string, layer int) string {
	if layer > 5 {
		layer = 5
	}
	dn := make([]rune, layer*2)
	copy(dn, []rune(fn))
	dir := make([]string, layer)
	for i := 0; i < (layer * 2); i += 2 {
		dir = append(dir, string(dn[i:i+2]))
	}
	root = Path(strings.Format("%s/", root), dir...)

	return filepath.Join(root, fn)
}

// PathTemp
// ****************************************************************************************************************************************
func PathTemp(dirs ...string) string {
	return Path(os.TempDir(), dirs...)
}

// PathWorking
// ****************************************************************************************************************************************
func PathWorking(dirs ...string) string {
	if path, err := os.Getwd(); err == nil {
		return Path(path, dirs...)
	}

	return PathTemp(dirs...)
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// FileOp
// ****************************************************************************************************************************************
type FileOp = uint32

// statusEnum *****************************************************************************************************************************
type statusEnum = uint8

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
