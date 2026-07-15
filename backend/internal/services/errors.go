package services

import "errors"

// Erreurs sentinelles communes aux services métier. Les handlers les
// mappent vers les codes HTTP demandés par la spec :
//
//	ErrValidation           -> 400
//	ErrForbidden            -> 403
//	repository.ErrNotFound  -> 404
//	repository.ErrConflict  -> 409
var (
	ErrValidation = errors.New("services: données invalides")
	ErrForbidden  = errors.New("services: action interdite")
)
