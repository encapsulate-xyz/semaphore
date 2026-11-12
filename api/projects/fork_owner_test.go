package projects

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"github.com/stretchr/testify/assert"
)

// 95898999: only project owners (and admins) change variable groups, CLI arguments or survey variables. A manager
// can still edit templates otherwise. These requests are refused before the handler reads the database.

func forkRequest(method, body string, role db.ProjectUserRole, admin bool, ctx map[string]any) *http.Request {
	r := httptest.NewRequest(method, "/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = helpers.SetContextValue(r, "user", &db.User{ID: 2, Admin: admin})
	r = helpers.SetContextValue(r, "projectUserRole", role)
	r = helpers.SetContextValue(r, "project", db.Project{ID: 1})
	for k, v := range ctx {
		r = helpers.SetContextValue(r, k, v)
	}
	return r
}

func TestIsOwnerOrAdmin(t *testing.T) {
	assert.True(t, isOwnerOrAdmin(forkRequest("GET", "", db.ProjectOwner, false, nil)))
	assert.True(t, isOwnerOrAdmin(forkRequest("GET", "", db.ProjectGuest, true, nil)), "an admin")
	for _, role := range []db.ProjectUserRole{db.ProjectManager, db.ProjectTaskRunner, db.ProjectGuest} {
		assert.False(t, isOwnerOrAdmin(forkRequest("GET", "", role, false, nil)), "role %s", role)
	}
}

func TestEnvironment_OnlyOwners(t *testing.T) {
	c := &EnvironmentController{}
	env := map[string]any{"environment": db.Environment{ID: 3, ProjectID: 1}}
	for name, call := range map[string]func(http.ResponseWriter, *http.Request){
		"add": c.AddEnvironment, "update": c.UpdateEnvironment, "remove": c.RemoveEnvironment,
	} {
		w := httptest.NewRecorder()
		call(w, forkRequest("POST", `{"name": "x", "json": "{}"}`, db.ProjectManager, false, env))
		assert.Equal(t, http.StatusForbidden, w.Code, name)
	}
}

func TestTemplate_ManagerCantSetConfiguration(t *testing.T) {
	for name, body := range map[string]string{
		"cli arguments":     `{"name": "t", "playbook": "p.yml", "arguments": "[\"-e\", \"x=1\"]"}`,
		"argument override": `{"name": "t", "playbook": "p.yml", "allow_override_args_in_task": true}`,
		"survey variables":  `{"name": "t", "playbook": "p.yml", "survey_vars": [{"name": "x"}]}`,
	} {
		w := httptest.NewRecorder()
		AddTemplate(w, forkRequest("POST", body, db.ProjectManager, false, nil))
		assert.Equal(t, http.StatusForbidden, w.Code, "add, "+name)

		w = httptest.NewRecorder()
		old := map[string]any{"template": db.Template{ID: 4, ProjectID: 1, Name: "t", Playbook: "p.yml"}}
		UpdateTemplate(w, forkRequest("PUT", body, db.ProjectManager, false, old))
		assert.Equal(t, http.StatusForbidden, w.Code, "update, "+name)
	}
}

// e4efa8be: an inventory can only be a file in a repository, not a static list or a terraform workspace.
func TestInventory_OnlyFileType(t *testing.T) {
	for _, typ := range []string{"static", "static-yaml", "terraform-workspace", ""} {
		body := `{"project_id": 1, "name": "i", "type": "` + typ + `", "inventory": "hosts"}`
		w := httptest.NewRecorder()
		AddInventory(w, forkRequest("POST", body, db.ProjectOwner, false, nil))
		assert.Equal(t, http.StatusBadRequest, w.Code, "add %q", typ)
		assert.Contains(t, w.Body.String(), "must be type file")

		body = `{"id": 5, "project_id": 1, "name": "i", "type": "` + typ + `", "inventory": "hosts"}`
		w = httptest.NewRecorder()
		UpdateInventory(w, forkRequest("PUT", body, db.ProjectOwner, false, map[string]any{"inventory": db.Inventory{ID: 5, ProjectID: 1}}))
		assert.Equal(t, http.StatusBadRequest, w.Code, "update %q", typ)
		assert.Contains(t, w.Body.String(), "must be type file")
	}
}
