package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/anonyindian/logger"
)

func DoUpdate(version string, chatId int64, msgId int) error {
	// err := gitPull()
	// if err != nil {
	// 	buildWithClone(".")
	// 	restart("giga", []string{}, 5, chatId, msgId, "Updated Successfully.")
	// }
	// err = buildBinary()
	// if err != nil {
	// 	return err
	// }
	if err := refreshChangelog(); err != nil {
		return err
	}
	if err := downloadUpdate(version); err != nil {
		return err
	}
	return restart("giga", []string{}, 5, chatId, msgId, "Updated Successfully.")
}

func refreshChangelog() error {
	origin := "https://raw.githubusercontent.com/GigaUserbot/GIGA/dev/changelog.json"
	resp, err := http.Get(origin)
	if err != nil {
		return err
	}
	b, err := os.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	const fileName = "changelog.json"
	_ = os.Remove(fileName)
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write(b)
	return err
}

func downloadUpdate(version string) error {
	os, err := GetSupportedOS()
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	arch, err := GetSupportedARCH()
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	url := fmt.Sprintf(
		"https://github.com/GigaUserbot/GIGA/releases/download/v%s/giga_%s_%s_%s",
		version, version, os, arch,
	)
	err = DownloadFile("giga", url)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	return nil
}

var CurrentUpdate = &Update{}
var currentVersion *version

type Update struct {
	Version string
	Changes []string
}

func InitUpdate(l *logger.Logger) {
	b, err := os.ReadFile("changelog.json")
	if err != nil {
		l.ChangeLevel(logger.LevelCritical).Printlnf("Failed to open changelog.json: %s\n", err.Error())
		return
	}
	err = json.Unmarshal(b, CurrentUpdate)
	if err != nil {
		l.ChangeLevel(logger.LevelCritical).Printlnf("Failed to parse changelog.json: %s\n", err.Error())
		return
	}
	currentVersion, err = parseVersion(CurrentUpdate.Version)
	if err != nil {
		l.ChangeLevel(logger.LevelCritical).Printlnf("Failed to parse current version: %s\n", err.Error())
		return
	}
}

func CheckChanges() (*Update, bool) {
	var u Update
	origin := "https://raw.githubusercontent.com/GigaUserbot/GIGA/dev/changelog.json"
	resp, err := http.Get(origin)
	if err != nil {
		return nil, false
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false
	}
	json.Unmarshal(b, &u)
	return &u, CompareVersion(u.Version)
}

// CompareVersion returns true if input version is greater than current one.
func CompareVersion(version string) bool {
	parsedVersion, err := parseVersion(version)
	if err != nil {
		return false
	}
	return currentVersion.compare(parsedVersion)
}

type version struct {
	minor int
	major int
	patch int
}

func parseVersion(s string) (*version, error) {
	v := version{}
	for index, velem := range strings.Split(s, ".") {
		vint, err := strconv.Atoi(velem)
		if err != nil {
			return nil, errors.New("failed to parse version")
		}
		switch index {
		case 0:
			v.major = vint
		case 1:
			v.minor = vint
		case 2:
			v.patch = vint
		}
	}
	return &v, nil
}

func (v *version) compare(v1 *version) bool {
	return v1.major > v.major || v1.minor > v.minor || v1.patch > v.patch
}
