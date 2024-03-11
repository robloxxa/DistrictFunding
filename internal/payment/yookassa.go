package payment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	yooKassaUrl = "https://api.yookassa.ru/v3/"
)

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type AuthorizationDetails struct {
	Rrn          string `json:"rrn"`
	AuthCode     string `json:"auth_code"`
	ThreeDSecure struct {
		Applied bool `json:"applied"`
	} `json:"three_d_secure"`
}

type PaymentMethod struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Saved bool   `json:"saved"`
	Title string `json:"title"`
	Card  *Card  `json:"card"`
}

type Card struct {
	First6        string `json:"first6"`
	Last4         string `json:"last4"`
	ExpiryMonth   string `json:"expiry_month"`
	ExpiryYear    string `json:"expiry_year"`
	CardType      string `json:"card_type"`
	IssuerCountry string `json:"issuer_country"`
	IssuerName    string `json:"issuer_name"`
}

type Payment struct {
	ID                   string                 `json:"id"`
	Status               string                 `json:"status"`
	Paid                 bool                   `json:"paid"`
	Amount               Amount                 `json:"amount"`
	AuthorizationDetails *AuthorizationDetails  `json:"authorization_details"`
	CreatedAt            time.Time              `json:"created_at"`
	Description          string                 `json:"description,omitempty"`
	ExpiresAt            time.Time              `json:"expires_at"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	PaymentMethod        *PaymentMethod         `json:"payment_method,omitempty"`
	Recipient            struct {
		AccountID string `json:"account_id"`
		GatewayID string `json:"gateway_id"`
	} `json:"recipient"`
	Refundable   bool `json:"refundable"`
	Test         bool `json:"test"`
	IncomeAmount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"income_amount"`
}

type (
	Me struct {
		AccountID            string         `json:"account_id"`
		Test                 bool           `json:"test"`
		Fiscalization        *Fiscalization `json:"fiscalization"`
		FiscalizationEnabled bool           `json:"fiscalization_enabled"`
		PaymentMethods       []string       `json:"payment_methods"`
		Status               string         `json:"status"`
	}

	Fiscalization struct {
		Provider string `json:"provider"`
		Enabled  bool   `json:"enabled"`
	}
)

type Yookassa struct {
	shopId int
	token  string
	c      *http.Client
}

func NewYookassa(shopId int, token string) *Yookassa {
	return &Yookassa{shopId, token, &http.Client{}}
}

// newRequest makes a request with Authorization Header and appends endpoint (like /me) to the yookassa url string
func (y *Yookassa) newRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	urlString, err := url.JoinPath(yooKassaUrl, endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %d:%s", y.shopId, y.token))
	return req, nil
}

func (y *Yookassa) Me() (*Me, error) {
	var me Me
	req, err := y.newRequest("GET", "/me", nil)
	if err != nil {
		return nil, err
	}

	res, err := y.c.Do(req)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&me); err != nil {
		return nil, err
	}

	return &me, nil
}

// CreatePayment creates new payment in yookassa via POST /payments
func (y *Yookassa) CreatePayment(idempotentcyKey string, payment *Payment) {

}
