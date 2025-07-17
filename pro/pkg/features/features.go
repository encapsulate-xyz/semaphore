package features

import "github.com/semaphoreui/semaphore/db"

type Features struct {
	ProjectRunners   bool `json:"project_runners"`
	TerraformBackend bool `json:"terraform_backend"`
	TaskSummary      bool `json:"task_summary"`
	SecretStorages   bool `json:"secret_storages"`
}

func GetFeatures(user *db.User) map[string]bool {
	return map[string]bool{
		"project_runners":   false,
		"terraform_backend": false,
		"task_summary":      false,
		"secret_storages":   false,
	}
}
