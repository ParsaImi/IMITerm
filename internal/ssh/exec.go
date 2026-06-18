package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/ParsaImi/imiterm/internal/model"
)

// Exec replaces the current process with an SSH session to the given host.
// On success, this function never returns — the process becomes ssh.
func Exec(h model.Host) error {
	argv, err := buildArgs(h)
	if err != nil {
		return err
	}

	binary, err := exec.LookPath(argv[0])
	if err != nil {
		return fmt.Errorf("cannot find %s in PATH: %w", argv[0], err)
	}

	if os.Getenv("TMUX") != "" {
		exec.Command("tmux", "rename-window", h.Name).Run()
	}

	return syscall.Exec(binary, argv, os.Environ())
}

func buildArgs(h model.Host) ([]string, error) {
	sshArgs := []string{
		"ssh",
		"-p", fmt.Sprintf("%d", h.Port),
		fmt.Sprintf("%s@%s", h.User, h.Hostname),
	}

	switch h.AuthMethod() {
	case model.AuthKey:
		// Insert -i keypath before user@host
		expanded := expandTilde(h.KeyPath)
		sshArgs = []string{
			"ssh",
			"-p", fmt.Sprintf("%d", h.Port),
			"-i", expanded,
			fmt.Sprintf("%s@%s", h.User, h.Hostname),
		}

	case model.AuthPass:
		if _, err := exec.LookPath("sshpass"); err != nil {
			return nil, fmt.Errorf("password auth requires sshpass (apt install sshpass)")
		}
		// sshpass wraps the ssh command, passing the password via stdin
		return []string{
			"sshpass", "-p", h.Password,
			"ssh",
			"-p", fmt.Sprintf("%d", h.Port),
			fmt.Sprintf("%s@%s", h.User, h.Hostname),
		}, nil

	case model.AuthAgent:
		// No extra flags — ssh-agent handles it
	}

	return sshArgs, nil
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	return path
}
