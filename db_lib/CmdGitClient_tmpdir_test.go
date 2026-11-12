package db_lib

import (
	"testing"

	"github.com/semaphoreui/semaphore/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// b3873566: ls-remote runs in the project's folder under tmp_path. Only a task run used to create that folder, so a
// deleted folder stopped the project's commit schedules (they call ls-remote) until someone ran a task by hand. The
// git client now creates the folder itself.
func TestCmdGitClient_CreatesMissingProjectFolder(t *testing.T) {
	original := util.Config
	t.Cleanup(func() { util.Config = original })
	util.Config = &util.ConfigType{TmpPath: t.TempDir(), Process: &util.ConfigProcess{}}
	dir := util.Config.GetProjectTmpDir(1)
	require.NoDirExists(t, dir)

	upstream := t.TempDir()
	gitInit(t, upstream)
	client := CreateCmdGitClient(nopKeyInstaller{})
	repo := newTestGitRepo(t, upstream, "main")

	hash, err := client.GetLastRemoteCommitHash(repo)
	require.NoError(t, err)
	assert.Len(t, hash, 40)
	assert.DirExists(t, dir)

	branches, err := client.GetRemoteBranches(repo)
	require.NoError(t, err)
	assert.Equal(t, []string{"main"}, branches)
}
