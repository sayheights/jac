// Package auth provides common Auth method implementations
// that satisfy the jac.Authorizer interface
package auth

type In int

const (
	InHeader In = iota + 1
	InQuery
	InBody
)
