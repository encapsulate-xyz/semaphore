package ssh

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// From upstream 20ebe7c05: Listen() creates the socket's folder when it doesn't exist (the project's folder under
// tmp_path, which a task run used to be the only thing to create).
func TestAgent_Listen_CreatesSocketDir(t *testing.T) {
	// Not t.TempDir(): its path exceeds the ~104-byte unix socket limit on macOS.
	tmp, err := os.MkdirTemp("/tmp", "ssh")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmp) }) //nolint:errcheck

	agent := Agent{
		SocketFile: filepath.Join(tmp, "project_2", "agent.sock"),
	}

	err = agent.Listen()
	require.NoError(t, err)
	require.NoError(t, agent.Close())
}
