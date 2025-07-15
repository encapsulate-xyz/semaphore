package subscriptions

import (
	"errors"
	"github.com/semaphoreui/semaphore/api/helpers"
	"net/http"
)

func GetSubscription(w http.ResponseWriter, r *http.Request) {
	helpers.WriteError(w, errors.New("not implemented"))

}

func Activate(w http.ResponseWriter, r *http.Request) {
	helpers.WriteError(w, errors.New("not implemented"))
}
