package api

import (
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/pkg/random"
)

func newEmailOtp(userID int, userEmail string, store db.Store) error {
	code := random.Number(6)
	_, err := store.AddEmailOtpVerification(userID, code)

	if err != nil {
		return err
	}
	err = sendEmailVerificationCode(code, userEmail)
	return err
}
