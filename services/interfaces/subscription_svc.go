package interfaces

import "github.com/semaphoreui/semaphore/db"

type SubscriptionService interface {
	HasActiveSubscription() bool
	CanAddProUser(store db.Store) (ok bool, err error)
}
