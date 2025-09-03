package domain

type PaymentType string

const (
	PaymentTypeContrib PaymentType = "contrib"
	PaymentTypeAdvance PaymentType = "advance"
)

type TaxScheme string

const (
	TaxSchemeUSN6 TaxScheme = "usn_6"
)

const (
	BpDen int64 = 10_000
)
