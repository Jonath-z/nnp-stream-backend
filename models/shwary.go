package models

type ShwaryPaymentRequest struct {
	Amount            float64 `json:"amount"`
	ClientPhoneNumber string  `json:"clientPhoneNumber"`
	CallbackUrl       string  `json:"callbackUrl,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	ReferenceID       string  `json:"referenceId,omitempty"`
}

type ShwaryTransaction struct {
	ID                   string  `json:"id"`
	UserID               string  `json:"userId"`
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
	Type                 string  `json:"type"`
	Status               string  `json:"status"`
	RecipientPhoneNumber string  `json:"recipientPhoneNumber"`
	ReferenceID          string  `json:"referenceId"`
	TxHash               string  `json:"txHash,omitempty"`
	FailureReason        string  `json:"failureReason,omitempty"`
	CompletedAt          string  `json:"completedAt,omitempty"`
	IsSandbox            bool    `json:"isSandbox"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}
