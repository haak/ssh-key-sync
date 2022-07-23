package command

import (
	"bytes"
	"path/filepath"
	"strings"
	"time"

	"github.com/shoenig/ssh-key-sync/internal/ssh"
	"gophers.dev/pkgs/atomicfs"
)

func generateFileContent(keys []ssh.Key, now time.Time) string {
	var buf bytes.Buffer

	formattedTime := now.Format(time.RFC1123)
	buf.WriteString("# Autogenerated by ssh-key-sync on " + formattedTime + "\n\n")

	for _, key := range keys {
		if key.Managed {
			buf.WriteString("# managed by ssh-key-sync\n")
		}
		buf.WriteString(key.Value)
		if key.User != "" && key.Host != "" {
			buf.WriteString(" ")
			buf.WriteString(key.User)
			buf.WriteString("@")
			buf.WriteString(key.Host)
		}
		buf.WriteString("\n\n")
	}

	return buf.String()
}

// safely write to a tmp file and then do an atomic rename

func (e *execer) writeToFile(file, user, content string) error {
	fw := atomicfs.NewFileWriter(atomicfs.Options{
		TmpDirectory: filepath.Dir(file),
		TmpExtension: "tmp",
		Mode:         0600,
	})

	if err := fw.Write(strings.NewReader(content), file); err != nil {
		return err
	}

	return e.touch(file, user)
}
