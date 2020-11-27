package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
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

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {

		ignorePaths, _ := cmd.Flags().GetStringArray("ignore_paths")
		projectPath, _ := os.Getwd()
		projectName, _ := cmd.Flags().GetString("name")
		backupBasePath, _ := cmd.Flags().GetString("path")
		backupFullPath := path.Join(backupBasePath, projectName)

		info, err := projectInfo(backupBasePath, projectName)
		if err != nil {
			panic(err)
		}
		if info.Pid > 0 {
			p, _ := os.FindProcess(info.Pid)
			w, _ := p.Wait()
			if w != nil {
				fmt.Println("theere already a backup running")
				os.Exit(1)
			}
		}

		backupArgs := []string{"backup"}
		backupArgs = append(backupArgs, backupFullPath, projectPath)
		for _, ignore := range ignorePaths {
			backupArgs = append(backupArgs, "--ignore_paths", ignore)
		}
		backupArgs = append(backupArgs, "--last_updated", info.LastUpdated.Format(time.RFC3339))
		fmt.Println(backupArgs)
		backup := exec.Command(os.Args[0], backupArgs...)
		if err := backup.Start(); err != nil {
			panic(err)
		}
		info.Pid = backup.Process.Pid
		if err := saveProjectInfo(backupBasePath, projectName, info); err != nil {
			panic(err)
		}
	},
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func init() {
	dir, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	startCmd.Flags().String("name", path.Base(dir), "the name of the backup project")
	startCmd.Flags().String("path", path.Join(home, ".backups"), "the path to the backup project")
	startCmd.Flags().StringArray("ignore_paths", []string{"node_modules", "dist", "build"}, "folders to ignore")
	rootCmd.AddCommand(startCmd)
}
