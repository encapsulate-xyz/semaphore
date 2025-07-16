package interfaces

type SubscriptionService interface {
	HasActiveSubscription() bool
	CanAddProUser() (ok bool, err error)
}
