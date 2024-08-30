package config

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	BASE_DIR   = filepath.Dir(filepath.Join(b, "..", ".."))
)
