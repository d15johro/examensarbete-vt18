// Package fs contains filesystem util functions for the experiment.
package fs

import "io/ioutil"

func FileCount(dirpath string) (uint32, error) {
	i := 0
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return 0, err
	}
	for _, file := range files {
		if !file.IsDir() {
			i++
		}
	}
	return uint32(i), nil
}
