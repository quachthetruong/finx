package apperrors

import (
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func IsConstraintViolationError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pgerrcode.IsIntegrityConstraintViolation(string(pqErr.Code)) {
			return true
		}
	}
	return false
}

func IsNotFoundError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pgerrcode.IsCaseNotFound(string(pqErr.Code)) {
			return true
		}
	}
	if errors.Is(err, qrm.ErrNoRows) {
		return true
	}
	return false
}

func IsObjectNotInPrerequisiteStateError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pgerrcode.IsObjectNotInPrerequisiteState(string(pqErr.Code)) {
			return true
		}
	}
	return false
}
