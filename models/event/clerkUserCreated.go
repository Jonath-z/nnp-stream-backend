package event

type ClerkUserEvent struct {
	Data            User            `json:"data"`
	EventAttributes EventAttributes `json:"event_attributes"`
	Object          string          `json:"object"`
	Timestamp       int64           `json:"timestamp"`
	Type            string          `json:"type"`
}

// User represents a user in the system
type User struct {
	Birthday              string                 `json:"birthday"`
	CreatedAt             int64                  `json:"created_at"`
	EmailAddresses        []EmailAddress         `json:"email_addresses"`
	ExternalAccounts      []interface{}          `json:"external_accounts"`
	ExternalID            string                 `json:"external_id"`
	FirstName             string                 `json:"first_name"`
	Gender                string                 `json:"gender"`
	ID                    string                 `json:"id"`
	ImageURL              string                 `json:"image_url"`
	LastName              string                 `json:"last_name"`
	LastSignInAt          int64                  `json:"last_sign_in_at"`
	Object                string                 `json:"object"`
	PasswordEnabled       bool                   `json:"password_enabled"`
	PhoneNumbers          []interface{}          `json:"phone_numbers"`
	PrimaryEmailAddressID string                 `json:"primary_email_address_id"`
	PrimaryPhoneNumberID  interface{}            `json:"primary_phone_number_id"`
	PrimaryWeb3WalletID   interface{}            `json:"primary_web3_wallet_id"`
	PrivateMetadata       map[string]interface{} `json:"private_metadata"`
	ProfileImageURL       string                 `json:"profile_image_url"`
	PublicMetadata        map[string]interface{} `json:"public_metadata"`
	TwoFactorEnabled      bool                   `json:"two_factor_enabled"`
	UnsafeMetadata        map[string]interface{} `json:"unsafe_metadata"`
	UpdatedAt             int64                  `json:"updated_at"`
	Username              interface{}            `json:"username"`
	Web3Wallets           []interface{}          `json:"web3_wallets"`
}

type EmailAddress struct {
	EmailAddress string            `json:"email_address"`
	ID           string            `json:"id"`
	LinkedTo     []interface{}     `json:"linked_to"`
	Object       string            `json:"object"`
	Verification EmailVerification `json:"verification"`
}

type EmailVerification struct {
	Status   string `json:"status"`
	Strategy string `json:"strategy"`
}

type EventAttributes struct {
	HTTPRequest HTTPRequest `json:"http_request"`
}

type HTTPRequest struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}
