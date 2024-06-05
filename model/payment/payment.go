package payment

type Mode string

const (
	TEST       Mode = "TEST"
	PRODUCTION Mode = "PRODUCTION"
)

type FormData struct {
	Key   string
	Value string
}
