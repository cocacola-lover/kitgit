package gitrepository

import (
	"fmt"
	"os"
	"path/filepath"

	ge "github.com/cocacola-lover/kitgit/pkg/gitError"

	"gopkg.in/ini.v1"
)

type GitRepository struct {
	workTree string
	gitDir   string
	config   *ini.File
}

func NewGitRepository(path string) (*GitRepository, error) {
	workTree := path
	gitDir := filepath.Join(path, ".git")
	configPath := filepath.Join(gitDir, "config")

	gitDirInfo, err := os.Stat(gitDir)
	if err != nil {
		return nil, err
	} else if !gitDirInfo.IsDir() {
		return nil, ge.NewGitError(
			fmt.Sprintf("Not a Git repository : %s", gitDir),
			os.ErrInvalid,
		)
	}

	config, err := ini.Load(configPath)
	if err != nil {
		return nil, ge.NewGitError("Configuration file is missing", err)
	}

	if ver, err := config.Section("core").GetKey("repositoryformatversion"); err != nil {
		return nil, ge.NewGitError("Broken config", err)
	} else if verInt, err := ver.Int(); err != nil {
		return nil, ge.NewGitError("Broken config", err)
	} else if verInt != 0 {
		return nil, ge.NewGitError(
			fmt.Sprintf("Unsupported repositoryformatversion %d", verInt),
			nil,
		)
	}

	return &GitRepository{
		workTree: workTree,
		gitDir:   gitDir,
		config:   config,
	}, nil
}
