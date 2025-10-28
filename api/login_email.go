package api

import (
	"net/http"

	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/pkg/random"
	log "github.com/sirupsen/logrus"
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

func resendEmailOtp(w http.ResponseWriter, r *http.Request) {

	session, ok := getSession(r)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := helpers.Store(r).GetUser(session.UserID)

	if err != nil {
		helpers.WriteErrorStatus(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = newEmailOtp(session.UserID, user.Email, helpers.Store(r))
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": session.UserID,
			"context": "resend_email_otp",
		}).Error("Failed to add email otp verification")
		helpers.WriteErrorStatus(w, "Failed to create email OTP verification", http.StatusInternalServerError)
		return
	}
}
