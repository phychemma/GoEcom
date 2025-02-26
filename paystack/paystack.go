package paystack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var paystackSecretKey = "sk_test_xxxx" // Replace with your actual Paystack secret key

// VerifyPayment checks if a payment was successful
func VerifyPayment(c *gin.Context) {
	reference := c.Param("reference") // Get transaction reference from URL

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+paystackSecretKey).
		SetHeader("Content-Type", "application/json").
		Get("https://api.paystack.co/transaction/verify/" + reference)

	if err != nil {
		log.Println("Error verifying payment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed"})
		return
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body(), &result)

	// Check if the transaction was successful
	if status, ok := result["status"].(bool); ok && status {
		data := result["data"].(map[string]interface{})
		if data["status"] == "success" {
			// Mark order as paid in your database
			log.Printf("Payment verified: %s", reference)
			c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "data": data})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not successful"})
}

// PaystackWebhook listens for Paystack payment notifications

type PaystackWebhook struct {
	Event string `json:"event"`
	Data  struct {
		Reference string `json:"reference"`
		Status    string `json:"status"`
		Amount    int    `json:"amount"`
	} `json:"data"`
}

func PaystackWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON request
	var webhookData PaystackWebhook
	err = json.Unmarshal(body, &webhookData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Log webhook event
	fmt.Printf("Received Paystack Webhook: %+v\n", webhookData)

	// Check event type and process accordingly
	if webhookData.Event == "charge.success" && webhookData.Data.Status == "success" {
		fmt.Println("✅ Payment successful! Reference:", webhookData.Data.Reference)
		// TODO: Update database, mark order as paid, etc.
	} else {
		fmt.Println("❌ Payment failed or not a charge.success event")
	}

	// Respond to Paystack to confirm receipt of webhook
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}
