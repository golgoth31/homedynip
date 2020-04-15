package dns

// DNS provider interface
type DNS interface {
	Write() error
}

type dummy struct{}
