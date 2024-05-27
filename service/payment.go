package service

type Payment struct {
	apikey string
}

func NewPaymentService(apikey string) Payment {
	return Payment{
		apikey: apikey,
	}
}
