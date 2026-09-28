package identity_role

import "errors"

var (
	ErrIdentityNotFound  = errors.New("identity not found")
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleExists        = errors.New("role already assigned")
	ErrInvalidAssignment = errors.New("invalid role assignment")
)
