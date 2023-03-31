package main

import (
	"bytes"
	"io"
	"os"

	"github.com/danilvpetrov/gobump/transformers"
)

// runTransformers runs transformers against the file. If transformers performed
// conflicting changes to the file, the last transformer always takes precedence.
func runTransformers(file string, tt ...transformers.Transformer) error {
	f, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	for _, t := range tt {
		ok, err := t(f, &buf)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		if err := f.Truncate(0); err != nil {
			return err
		}

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return err
		}
		if _, err := buf.WriteTo(f); err != nil {
			return err
		}

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return err
		}

		buf.Reset()
	}

	return nil
}
