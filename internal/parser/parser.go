package parser

import (
	"fmt"
	"os"
	"path/filepath"
)

func processYAML(outdir string, content []byte, data map[string]interface{}) error {
	unknownPath := filepath.Join(outdir, "other.yaml")
	for t, file := range typeMap {
		if file == "" {
			fmt.Printf("Skipping '%s' (%s)\n", filepath.Join(outdir, file), t)
			continue
		}
		if isType(data, t) {
			fmt.Printf("Adding '%s' (%s)\n", filepath.Join(outdir, file), t)
			if err := appendFile(content, filepath.Join(outdir, file)); err != nil {
				return err
			}

			return nil
		}
	}

	fmt.Printf("Adding(unknown) '%s' (%s)\n", unknownPath, getType(data))
	if err := appendFile(content, unknownPath); err != nil {
		return err
	}

	return nil
}

func appendFile(content []byte, file string) error {
	_, err := os.Stat(file)
	if err == nil {
		content = append([]byte("---\n"), content...)
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(content); err != nil {
		return err
	}

	return nil
}
