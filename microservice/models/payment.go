package models

// PaymentForm represents the structure of the payment form data.
type PaymentForm struct {
	CardNumber     string `json:"cardNumber"`
	ExpirationDate string `json:"expirationDate"`
	CVV            string `json:"cvv"`
	Name           string `json:"name"`
	Address        string `json:"address"`
}

type PaymentDetails struct {
	CardNumber     string `json:"cardNumber"`
	ExpirationDate string `json:"expirationDate"`
	Name           string `json:"name"`
	Address        string `json:"address"`
}
