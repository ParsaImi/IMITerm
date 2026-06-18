package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// IsRepo returns true if dir is inside a git working tree.
func IsRepo(dir string) bool {
	return runGit(dir, "rev-parse", "--is-inside-work-tree") == nil
}

// Pull runs git pull --rebase --autostash in dir.
func Pull(dir string) error {
	return runGit(dir, "pull", "--rebase", "--autostash")
}

// CommitAndPush stages config.toml, commits with a timestamp, and pushes.
// Returns nil if there's nothing to commit.
func CommitAndPush(dir string) error {
	configFile := filepath.Join(dir, "config.toml")

	if err := runGit(dir, "add", configFile); err != nil {
		return fmt.Errorf("git add: %w", err)
	}

	// Check if there's anything to commit
	if err := runGit(dir, "diff", "--cached", "--quiet"); err == nil {
		return nil
	}

	msg := fmt.Sprintf("imiterm: update config %s", time.Now().Format("2006-01-02 15:04:05"))
	if err := runGit(dir, "commit", "-m", msg); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}

	if err := runGit(dir, "push"); err != nil {
		return fmt.Errorf("git push: %w", err)
	}

	return nil
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, out)
	}
	return nil
}
