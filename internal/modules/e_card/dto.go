package ecard


type CardResponse struct {
	CardNumber string `json:"card_number"`
	Expiry     string `json:"expiry"`
	Status     string `json:"status"`
}