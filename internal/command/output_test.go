package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shoenig/ssh-key-sync/internal/ssh"

	"github.com/stretchr/testify/require"
)

func Test_writeToFile(t *testing.T) {
	tmpFilename := filepath.Join(os.TempDir(), "test1")

	e := &execer{fakeChown: true}
	err := e.writeToFile(tmpFilename, "somebody", "abc123")
	require.NoError(t, err)

	bs, err := ioutil.ReadFile(tmpFilename)
	require.NoError(t, err)
	require.Equal(t, "abc123", string(bs))
}

func Test_generateFileContent(t *testing.T) {

	keys := []ssh.Key{
		{Managed: false, Value: "aaaaaaa"},
		{Managed: false, Value: "jjjjjjj"},
		{Managed: false, Value: "bbbbbbb", User: "bob", Host: "b1"},
		{Managed: true, Value: "ccccccc", User: "alice", Host: "a1"},
		{Managed: true, Value: "ddddddd", User: "alice"},
	}

	now := time.Date(2017, 12, 17, 12, 44, 0, 0, time.UTC)

	output := generateFileContent(keys, now)

	require.Equal(t, exp1, output)
}

const exp1 = `# Autogenerated by ssh-key-sync on Sun, 17 Dec 2017 12:44:00 UTC

aaaaaaa

jjjjjjj

bbbbbbb bob@b1

# managed by ssh-key-sync
ccccccc alice@a1

# managed by ssh-key-sync
ddddddd

`
