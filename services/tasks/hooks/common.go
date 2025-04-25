package hooks

import "github.com/semaphoreui/semaphore/db"

type Hook interface {
	End() error
}

func GetHook(app db.TemplateApp) Hook {
	return nil
}
