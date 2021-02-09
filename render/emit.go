package render

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

// Write emits the given artifact into the destination.
func Write(dir string, artifact Artifact) error {
	switch t := artifact.(type) {
	case *File:
		dst := filepath.Join(dir, t.FileName)
		if err := ioutil.WriteFile(dst, t.Buf, 0600); err != nil {
			return fmt.Errorf("unable to emit file: %w", err)
		}
	case *Dir:
		dst := filepath.Join(dir, t.DirName)
		if err := os.MkdirAll(dst, 0600); err != nil {
			return fmt.Errorf("unable to create directory: %s: %w", dst, err)
		}

		for _, file := range t.Files {
			if err := Write(dst, file); err != nil {
				return fmt.Errorf("unable to emit file: %s: %w", dst, err)
			}
		}

		for _, d := range t.Dirs {
			if err := Write(dst, d); err != nil {
				return fmt.Errorf("unable to emit dir: %w", err)
			}
		}
	default:
		return fmt.Errorf("invalid artifact type: %s", reflect.TypeOf(t).String())
	}

	return nil
}
