package utils

import (
	"os/exec"

	git "github.com/go-git/go-git/v5"
)

func buildBinary() error {
	return exec.Command("go", "build", ".").Run()
}

func buildWithClone(dir string) error {
	if err := gitClone(dir); err != nil {
		return err
	}
	if err := exec.Command("go", "build", dir, "-o=giga").Run(); err != nil {
		return err
	}
	return nil
}

func gitPull() error {
	return exec.Command("git", "pull").Run()
}

func gitClone(dir string) error {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               "https://github.com/GigaUserbot/GIGA",
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	return err
}
