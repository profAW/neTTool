package helper

import (
	"os"
	"path/filepath"
)

// https://stackoverflow.com/questions/33450980/how-to-remove-all-contents-of-a-directory-using-golang/52685448
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
