package nats

type Observant struct {
	ID string
}

func (s *Observant) Subscribe(consumerName string, subject string) string {
	return "subscribed"
}
