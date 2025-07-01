package querymod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacePlaceholders(t *testing.T) {
	originalQuery := `
        SELECT loan_package_offer.id AS "loan_package_offer.id"
        FROM public.loan_package_offer
        WHERE loan_package_offer.loan_package_request_id IN ($1::bigint, $2::bigint, $3::bigint);
`
	expect := `
        SELECT loan_package_offer.id AS "loan_package_offer.id"
        FROM public.loan_package_offer
        WHERE loan_package_offer.loan_package_request_id IN (1::bigint, 2::bigint, 3::bigint);
`
	values := []interface{}{int64(1), int64(2), int64(3)}
	res := ReplacePlaceholders(originalQuery, values)
	assert.Equal(t, expect, res)
}
