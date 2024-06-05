package service

import (
	"alc/model/cart"
	"alc/model/checkout"
	"alc/model/payment"
	"cmp"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
	"time"
)

type Payment struct {
	mode    payment.Mode
	storeId string
	apiKey  string
}

func NewPaymentService(mode payment.Mode, storeId string, apiKey string) Payment {
	return Payment{
		mode:    mode,
		storeId: storeId,
		apiKey:  apiKey,
	}
}

func (ps Payment) GetPaymentData(order checkout.Order, cartItems []cart.Item) []payment.FormData {
	// Get final amount
	amount := 0
	for _, cartItem := range cartItems {
		amount += cartItem.Quantity * cartItem.Product.Price
	}

	formData := []payment.FormData{
		{Key: "vads_action_mode", Value: "IFRAME"},
		{Key: "vads_amount", Value: fmt.Sprintf("%d", amount)},
		{Key: "vads_ctx_mode", Value: string(ps.mode)},
		{Key: "vads_currency", Value: "840"},
		{Key: "vads_page_action", Value: "PAYMENT"},
		{Key: "vads_payment_config", Value: "SINGLE"},
		{Key: "vads_site_id", Value: ps.storeId},
		{Key: "vads_trans_date", Value: time.Now().UTC().Format("20060102150405")},
		{Key: "vads_trans_id", Value: "xrT15p"},
		{Key: "vads_version", Value: "V2"},
	}

	// Sort the form data alphabetically by key
	slices.SortFunc(formData, func(a, b payment.FormData) int {
		return cmp.Compare(a.Key, b.Key)
	})

	// Get the values and append the apiKey
	values := make([]string, 0, len(formData)+1)
	for _, field := range formData {
		values = append(values, field.Value)
	}
	values = append(values, ps.apiKey)
	valuesJoin := strings.Join(values, "+")

	// Encoding the values to get the signature
	h := hmac.New(sha256.New, []byte(ps.apiKey))
	h.Write([]byte(valuesJoin))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Append the signature
	formData = append(formData, payment.FormData{Key: "signature", Value: signature})

	return formData
}
