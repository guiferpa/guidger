package webhook

type Webhook interface {
	ValidateSignature() error
	ValidateNonce() error
}

type whk struct{}

func (w *whk) ValidateSignature() error {
	return nil
}

func (w *whk) ValidateNonce() error {
	return nil
}

func New() Webhook {
	return &whk{}
}
