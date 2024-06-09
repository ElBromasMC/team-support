package service

import (
	"alc/model/checkout"
	"alc/model/payment"
	"cmp"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
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

func (ps Payment) GetPaymentData(order checkout.Order, trans checkout.Transaction) []payment.FormData {
	formData := []payment.FormData{
		{Key: "vads_action_mode", Value: "IFRAME"},
		{Key: "vads_amount", Value: fmt.Sprintf("%d", trans.Amount)},
		{Key: "vads_ctx_mode", Value: string(ps.mode)},
		{Key: "vads_currency", Value: "840"},
		{Key: "vads_page_action", Value: "PAYMENT"},
		{Key: "vads_payment_config", Value: "SINGLE"},
		{Key: "vads_site_id", Value: ps.storeId},
		{Key: "vads_trans_date", Value: trans.CreatedAt.UTC().Format("20060102150405")},
		{Key: "vads_trans_id", Value: trans.TransId},
		{Key: "vads_version", Value: "V2"},
		{Key: "vads_order_id", Value: order.Id.String()},
		{Key: "vads_cust_email", Value: order.Email},
		{Key: "vads_cust_first_name", Value: order.Name},
		{Key: "vads_cust_cell_phone", Value: order.Phone},
		{Key: "vads_cust_country", Value: "PE"},
		{Key: "vads_ship_to_street", Value: order.Address},
		{Key: "vads_ship_to_zip", Value: order.PostalCode},
		{Key: "vads_ship_to_city", Value: order.City},
		{Key: "vads_ship_to_country", Value: "PE"},
		{Key: "vads_ship_to_first_name", Value: order.Name},
		{Key: "vads_ship_to_phone_num", Value: order.Phone},
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
