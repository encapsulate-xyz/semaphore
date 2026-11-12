package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for the rules the encapsulate-xyz fork adds on top of upstream Semaphore.

// 50116606: a repository must track main.
func TestRepository_ValidateMainBranch(t *testing.T) {
	repo := func(branch string) Repository {
		return Repository{Name: "repo", GitURL: "git@github.com:encapsulate-xyz/repo.git", GitBranch: branch}
	}
	assert.NoError(t, repo("main").Validate())
	for _, b := range []string{"master", "dev", "", "main2", " main"} {
		assert.Error(t, repo(b).Validate(), "branch %q", b)
	}
}

// e4efa8be: an inventory is always read from a repository. d95b8da2: a runner tag is optional, but not empty.
func TestInventory_Validate(t *testing.T) {
	repoID := 1
	empty, tag := "", "runner-1"
	assert.Error(t, Inventory{}.Validate(), "no repository")
	assert.NoError(t, Inventory{RepositoryID: &repoID}.Validate())
	assert.NoError(t, Inventory{RepositoryID: &repoID, RunnerTag: &tag}.Validate())
	assert.Error(t, Inventory{RepositoryID: &repoID, RunnerTag: &empty}.Validate())
}

// c53635e4: a template can't point at another branch.
func TestTemplate_ValidateMainBranch(t *testing.T) {
	invID := 1
	tpl := func(branch *string) *Template {
		return &Template{Name: "t", Playbook: "p.yml", App: AppAnsible, InventoryID: &invID, GitBranch: branch}
	}
	main, empty, dev := "main", "", "dev"
	assert.NoError(t, tpl(nil).Validate())
	assert.NoError(t, tpl(&empty).Validate(), "empty means the repository's branch")
	assert.NoError(t, tpl(&main).Validate())
	require.Error(t, tpl(&dev).Validate())
}

// 95898999: only owners manage a project's configuration (variable groups, CLI arguments, survey variables).
func TestProjectUser_OnlyOwnersManageConfiguration(t *testing.T) {
	assert.True(t, ProjectOwner.Can(CanManageProjectConfiguration))
	for _, role := range []ProjectUserRole{ProjectManager, ProjectTaskRunner, ProjectGuest} {
		assert.False(t, role.Can(CanManageProjectConfiguration), "role %s", role)
	}
}
