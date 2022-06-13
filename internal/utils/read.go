package utils

import (
	"io/ioutil"
	"path"
)

func ReadDir(dpath string) []string {
	return readDir(dpath, "")
}

// readDir will read everything in the dpath
// if it is a directory it will take the name of the directory
// and prepend it to the filepath until it reaches a file
// e.g.
// - dir1
//  - dir2
//   - f3.txt
//  - f1.txt
// 	- f2.txt
// will return [ f1.txt, f2.txt, dir2/f3.txt ]
func readDir(dpath, name string) []string {
	fis, err := ioutil.ReadDir(dpath)
	if err != nil {
		return nil
	}

	filepaths := make([]string, 0)
	for _, fi := range fis {
		if fi.IsDir() {
			newdpath := path.Join(dpath, fi.Name())
			filepaths = append(filepaths, readDir(newdpath, fi.Name())...)
			continue
		}

		fname := fi.Name()
		if name != "" {
			fname = path.Join(name, fname)
		}
		filepaths = append(filepaths, fname)
	}
	return filepaths
}
