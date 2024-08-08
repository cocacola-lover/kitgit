package gitinit

import (
	"errors"
	ge "example/git-clone/pkg/gitError"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Func assumes that path leads to directory
func isEmptyDir(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func createPathAnd(do func(path string) error, path ...string) error {
	writeTo, err := filePath(true, path...)
	if err != nil {
		return err
	}
	err = do(writeTo)
	return err
}

func filePath(mkDir bool, path ...string) (string, error) {
	pathToFile := filepath.Join(path...)
	_, err := dirPath(mkDir, filepath.Dir(pathToFile))
	if err != nil {
		return "", err
	}

	return pathToFile, nil
}

// Creates dir if dir is absent in path
func dirPath(mkDir bool, path ...string) (string, error) {
	pathToFile := filepath.Join(path...)

	fileInfo, err := os.Stat(pathToFile)
	// If does not exist
	if errors.Is(err, os.ErrNotExist) {
		if mkDir {
			err := os.MkdirAll(pathToFile, 0750)
			return pathToFile, err
		} else {
			return "", ge.NewGitError(
				fmt.Sprintf("Dir %s does not exist", pathToFile),
				os.ErrNotExist,
			)
		}
	}
	// If does exist
	if err == nil {
		if fileInfo.IsDir() {
			return pathToFile, nil
		} else {
			return "", ge.NewGitError(
				fmt.Sprintf("%s is not a directory", pathToFile),
				os.ErrInvalid,
			)
		}
	}

	return "", err
}

func createDefaultConfig() (*ini.File, error) {
	iniData := ini.Empty()

	coreSection, err := iniData.NewSection("core")
	if err != nil {
		return nil, err
	}
	_, err = coreSection.NewKey("repositoryformatversion", "0")
	if err != nil {
		return nil, err
	}
	_, err = coreSection.NewKey("filemode", "false")
	if err != nil {
		return nil, err
	}
	_, err = coreSection.NewKey("bare", "false")
	if err != nil {
		return nil, err
	}

	return iniData, nil
}
