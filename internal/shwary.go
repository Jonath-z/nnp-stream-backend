package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nnp-stream-backend/models"
)

const (
	shwaryDefaultBaseURL = "https://api.shwary.com"
	shwaryCountryDRC     = "DRC"
	shwaryCurrencyUSD    = "USD"
)

func shwaryBaseURL() string {
	if v := os.Getenv("SHWARY_BASE_URL"); v != "" {
		return v
	}
	return shwaryDefaultBaseURL
}

func shwaryPaymentPath() string {
	if os.Getenv("SHWARY_SANDBOX") == "true" {
		return fmt.Sprintf("/api/v1/merchants/payment/sandbox/%s", shwaryCountryDRC)
	}
	return fmt.Sprintf("/api/v1/merchants/payment/%s", shwaryCountryDRC)
}

func shwaryHTTPClient() *http.Client {
	return &http.Client{Timeout: 15 * time.Second}
}

func shwaryAuthHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-merchant-id", os.Getenv("SHWARY_MERCHANT_ID"))
	req.Header.Set("x-merchant-key", os.Getenv("SHWARY_MERCHANT_KEY"))
}

// InitiateShwaryPayment kicks off a DRC mobile-money payment in USD.
func InitiateShwaryPayment(amount float64, phoneNumber, referenceID string) (models.ShwaryTransaction, error) {
	var tx models.ShwaryTransaction

	body := models.ShwaryPaymentRequest{
		Amount:            amount,
		ClientPhoneNumber: phoneNumber,
		CallbackUrl:       os.Getenv("SHWARY_CALLBACK_URL"),
		Currency:          shwaryCurrencyUSD,
		ReferenceID:       referenceID,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return tx, fmt.Errorf("marshal shwary payload: %w", err)
	}

	url := shwaryBaseURL() + shwaryPaymentPath()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return tx, fmt.Errorf("build shwary request: %w", err)
	}
	shwaryAuthHeaders(req)

	resp, err := shwaryHTTPClient().Do(req)
	if err != nil {
		return tx, fmt.Errorf("call shwary: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tx, fmt.Errorf("read shwary response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return tx, fmt.Errorf("shwary returned %d: %s", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, &tx); err != nil {
		return tx, fmt.Errorf("decode shwary response: %w", err)
	}
	return tx, nil
}

// GetShwaryTransaction fetches a transaction by its Shwary UUID.
func GetShwaryTransaction(id string) (models.ShwaryTransaction, error) {
	var tx models.ShwaryTransaction

	url := fmt.Sprintf("%s/api/v1/merchants/transactions/%s", shwaryBaseURL(), id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return tx, fmt.Errorf("build shwary request: %w", err)
	}
	shwaryAuthHeaders(req)

	resp, err := shwaryHTTPClient().Do(req)
	if err != nil {
		return tx, fmt.Errorf("call shwary: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tx, fmt.Errorf("read shwary response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return tx, fmt.Errorf("shwary returned %d: %s", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, &tx); err != nil {
		return tx, fmt.Errorf("decode shwary response: %w", err)
	}
	return tx, nil
}
