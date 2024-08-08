package gitrepository

import (
	"path/filepath"
)

func (git *GitRepository) Path(path ...string) string {
	return filepath.Join(git.gitDir, filepath.Join(path...))
}
