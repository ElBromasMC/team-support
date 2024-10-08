package admin

import (
	"alc/service"
)

type Handler struct {
	AdminService    service.Admin
	AuthService     service.Auth
	DeviceService   service.Device
	CurrencyService service.Currency
	SurveyService   service.Survey
}
