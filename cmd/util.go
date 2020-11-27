package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type backupInfo struct {
	Pid         int       `json:"pid"`
	LastUpdated time.Time `json:"last_updated"`
}

func projectInfo(path string, name string) (backupInfo, error) {
	infoJson := filepath.Join(path, "backups_info.json")

	var info backupInfo
	if fileExists(infoJson) {
		var infos map[string]backupInfo
		b, _ := ioutil.ReadFile(infoJson)
		if err := json.Unmarshal(b, &infos); err != nil {
			return info, err
		}
		info = infos[name]
	}
	return info, nil
}

func saveProjectInfo(path string, name string, info backupInfo) error {
	infoJson := filepath.Join(path, "backups_info.json")

	var infos map[string]backupInfo
	if fileExists(infoJson) {
		b, _ := ioutil.ReadFile(infoJson)
		if err := json.Unmarshal(b, &infos); err != nil {
			return err
		}
	} else {
		infos = map[string]backupInfo{}
	}
	infos[name] = info
	b, _ := json.Marshal(infos)
	return ioutil.WriteFile(infoJson, b, 0644)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func contains(what string, in []string) bool {
	for _, i := range in {
		if what == i {
			return true
		}
	}
	return false
}
