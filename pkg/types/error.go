package types

import "errors"

var (
	ErrSetImageOrFile = errors.New("image name or image file must be specified")
    ErrorCreateDockerExtractor = errors.New("error create docker extractor")
    ErrorAnalyze = errors.New("error analyze")
)
