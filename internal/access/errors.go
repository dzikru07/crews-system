package access

import "errors"

var (
	ErrPermissionNotFound = errors.New("permission not found")
	ErrEndpointNotFound   = errors.New("endpoint not found")
	ErrAccessDenied       = errors.New("access denied")
	ErrPermissionExists   = errors.New("permission already assigned")
)
