package repository

import (
	"time"
)

type FinancingRepository interface {
	GetDateAfter(date time.Time, workingDays int) (time.Time, error)
}
