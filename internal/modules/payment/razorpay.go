package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/razorpay/razorpay-go"
)

type razorpayService struct {
	key    string
	secret string
	client *razorpay.Client
}

func NewRazorpayService(key, secret string) Service {

	client := razorpay.NewClient(key, secret)

	return &razorpayService{
		key:    key,
		secret: secret,
		client: client,
	}
}

func (r *razorpayService) CreateOrder(amount int64, userID uint) (string, error) {

	data := map[string]interface{}{
		"amount":   amount, // amount in paise (IMPORTANT)
		"currency": "INR",
		"receipt":  fmt.Sprintf("cryptox_%d", time.Now().UnixNano()),
		"notes": map[string]interface{}{
			"user_id": userID,
			"source":  "cryptox",
		},
	}

	order, err := r.client.Order.Create(data, nil)
	if err != nil {
		return "", err
	}

	return order["id"].(string), nil
}

func (r *razorpayService) VerifySignature(orderID, paymentID, signature string) bool {

	payload := orderID + "|" + paymentID

	h := hmac.New(sha256.New, []byte(r.secret))
	h.Write([]byte(payload))

	expected := hex.EncodeToString(h.Sum(nil))

	return expected == signature
}

func (r *razorpayService) CreatePayout(userID uint, amount int64, name, ifsc, account string) (string, error) {

	//  MOCK MODE
	if os.Getenv("USE_MOCK_PAYOUT") == "true" {
		return fmt.Sprintf("mock_payout_%d", time.Now().UnixNano()), nil
	}

	// REAL MODE (keep for later)
	url := "https://api.razorpay.com/v1/payouts"

	payload := map[string]interface{}{
		"account_number": "YOUR_ACCOUNT_NUMBER",
		"fund_account": map[string]interface{}{
			"account_type": "bank_account",
			"bank_account": map[string]interface{}{
				"name":           name,
				"ifsc":           ifsc,
				"account_number": account,
			},
		},
		"amount":   amount,
		"currency": "INR",
		"mode":     "IMPS",
		"purpose":  "payout",
		"reference_id": fmt.Sprintf("withdraw_%d", time.Now().UnixNano()),
	}

	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.SetBasicAuth(r.key, r.secret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("payout failed: %s", string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return result["id"].(string), nil
}