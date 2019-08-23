package providers

type Sender interface {
	Send(msg string, receiver Receiver) error
}

type Creator func(opt map[string]string) (Sender, error)

type Receiver struct {
	Provider string
	To       []string
}
