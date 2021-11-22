package storage

import "fmt"

type Sms struct {
	Id string
	To string
	ProviderCode string
	RequestId string
	CreatedAt int
	AuthorizedAt int
	CancelledAt int
}

type SmsStorage struct {
}

func (s *SmsStorage) findSmsById(id string) (*Sms, error) {
	return nil, fmt.Errorf("sms not found")
}