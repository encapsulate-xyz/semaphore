package server

import (
	"github.com/semaphoreui/semaphore/services/interfaces"
)

func NewSubscriptionService() interfaces.SubscriptionService {
	return nil
}

type SubscriptionServiceImpl struct {
}

func (s *SubscriptionServiceImpl) HasActiveSubscription() bool {
	return false
}

func (s *SubscriptionServiceImpl) CanAddProUser() (ok bool, err error) {
	return false, nil
}
