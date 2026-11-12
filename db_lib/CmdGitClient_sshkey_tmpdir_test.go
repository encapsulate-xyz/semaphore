package db_lib

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"os"
	"testing"

	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/pkg/ssh"
	"github.com/semaphoreui/semaphore/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cryptossh "golang.org/x/crypto/ssh"
)

// Every Platform repository in Semaphore uses an SSH key. Before git runs, the key installer starts an SSH agent
// whose socket lives in the project's folder under tmp_path, so with that folder deleted the agent failed to start
// ("listening on socket ...: no such file or directory") before the git client could create the folder, and the
// project's commit schedules stayed silent until someone ran a task by hand. The agent now creates the folder.
// TestCmdGitClient_CreatesMissingProjectFolder covers the same path without a key.
func TestCmdGitClient_SSHKeyRepoWithMissingProjectFolder(t *testing.T) {
	// Not t.TempDir(): its path exceeds the ~104-byte unix socket limit on macOS.
	tmp, err := os.MkdirTemp("/tmp", "sem")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmp) }) //nolint:errcheck

	original := util.Config
	t.Cleanup(func() { util.Config = original })
	util.Config = &util.ConfigType{TmpPath: tmp, Process: &util.ConfigProcess{}}
	dir := util.Config.GetProjectTmpDir(1)
	require.NoDirExists(t, dir)

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	block, err := cryptossh.MarshalPrivateKey(priv, "")
	require.NoError(t, err)

	projectID := 1
	upstream := t.TempDir()
	gitInit(t, upstream)
	repo := newTestGitRepo(t, upstream, "main")
	repo.Repository.SSHKey = db.AccessKey{
		ID:        7,
		ProjectID: &projectID,
		Type:      db.AccessKeySSH,
		SshKey:    db.SshKey{Login: "git", PrivateKey: string(pem.EncodeToMemory(block))},
	}
	client := CreateCmdGitClient(ssh.KeyInstaller{})

	hash, err := client.GetLastRemoteCommitHash(repo)
	require.NoError(t, err)
	assert.Len(t, hash, 40)
	assert.DirExists(t, dir)
}
