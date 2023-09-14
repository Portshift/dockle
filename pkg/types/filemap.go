package types

import (
	"io"
	"os"
)

type FileMap map[string]FileData

type FileData struct {
	ContentReader io.ReadCloser
	FileMode      os.FileMode
}

type FilterFunc func(filePath string, fileMode os.FileMode) (bool, error)
