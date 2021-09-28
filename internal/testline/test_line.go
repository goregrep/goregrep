// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testline

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// New reports file and line number information about function invocations.
func New() string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	return "It was not possible to recover file and line number information about function invocations!"
}
