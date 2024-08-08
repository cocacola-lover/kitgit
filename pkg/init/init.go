package gitinit

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	ge "github.com/cocacola-lover/kitgit/pkg/gitError"
)

func InitCmd(cmdLineArgs ...string) error {
	initConfig := flag.NewFlagSet("init", flag.ExitOnError)

	err := initConfig.Parse(cmdLineArgs)
	if err != nil {
		return err
	}

	args := initConfig.Args()
	if len(args) == 0 {
		err = initGitRepo(".")
	} else {
		err = initGitRepo(filepath.Join(args...))
	}
	return err
}

func initGitRepo(path string) error {
	gitDirPath := filepath.Join(path, ".git")

	if workTree, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(gitDirPath, 0750); err != nil {
			return err
		}
	} else if err == nil {
		if !workTree.IsDir() {
			return ge.NewGitError(
				fmt.Sprintf("%s is not a git worktree", path),
				os.ErrInvalid,
			)
		} else if gitDir, err := os.Stat(gitDirPath); errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(gitDirPath, 0750); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if !gitDir.IsDir() {
			return ge.NewGitError(
				fmt.Sprintf("%s is not a directory", gitDirPath),
				os.ErrInvalid,
			)
		} else if isEmpty, err := isEmptyDir(gitDirPath); err != nil {
			return err
		} else if !isEmpty {
			return ge.NewGitError(
				fmt.Sprintf("%s is not empty", gitDirPath),
				os.ErrInvalid,
			)
		}
	}

	for _, dirToCreate := range []string{
		"branches",
		"objects",
		filepath.Join("refs", "tags"),
		filepath.Join("refs", "heads"),
	} {
		_, err := dirPath(true, gitDirPath, dirToCreate)
		if err != nil {
			return err
		}
	}

	// .git/description
	err := createPathAnd(func(path string) error {
		return os.WriteFile(path, []byte("Unnamed repository; edit this file 'description' to name the repository.\n"), 0644)
	}, gitDirPath, "description")
	if err != nil {
		return err
	}

	// .git/HEAD
	err = createPathAnd(func(path string) error {
		return os.WriteFile(path, []byte("ref: refs/heads/main\n"), 0644)
	}, gitDirPath, "HEAD")
	if err != nil {
		return err
	}

	// .git/config
	err = createPathAnd(func(path string) error {
		config, err := createDefaultConfig()
		if err != nil {
			return err
		}
		return config.SaveTo(path)
	}, gitDirPath, "config")
	return err
}
