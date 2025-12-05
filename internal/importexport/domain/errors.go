package domain

import "errors"

var (
	ErrJobNotFound        = errors.New("import/export job not found")
	ErrFilePathRequired   = errors.New("file path is required")
	ErrInvalidFileFormat  = errors.New("invalid file format")
	ErrJobAlreadyRunning  = errors.New("job is already running")
	ErrJobNotProcessing   = errors.New("job is not in processing state")
)
