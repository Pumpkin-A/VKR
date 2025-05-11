package models

type (
	PaymentStatus string
)

var (
	SuccessPaymentStatus      PaymentStatus = "success"
	FailedPaymentStatus       PaymentStatus = "failed"
	InProcessingPaymentStatus PaymentStatus = "inProcessing"
	ErrorStatus               PaymentStatus = "error"
	RefundedStatus            PaymentStatus = "refunded"
	CancelledStatus           PaymentStatus = "cancelled"
)

type EventExternalPaymentResult struct {
	UUID   string        `json:"id"`
	Status PaymentStatus `json:"status"`
	Error  string        `json:"errorText"`
}
