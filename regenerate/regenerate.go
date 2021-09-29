// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package regenerate regenerates some generated code.
package regenerate

import (
	"fmt"
	"os"
	"path/filepath"
)

func Regenerate(pattern string, opts ...Option) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("file path glob: %w", err)
	}

	for _, match := range matches {
		f, err := os.OpenFile(match, os.O_RDWR, os.ModePerm)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "File %q not exist.\n", match)

			continue

		} else if err != nil {
			return fmt.Errorf("os open file: %w", err)
		}

		err = Pipe(f, f, opts...)
		if err != nil {
			return fmt.Errorf("regenerate pipe: %s", err)
		}
	}

	return nil
}
