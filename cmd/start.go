package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {

		ignorePaths, _ := cmd.Flags().GetStringArray("ignore_path")
		projectPath, _ := os.Getwd()
		projectName, _ := cmd.Flags().GetString("name")
		backupBasePath, _ := cmd.Flags().GetString("path")
		backupFullPath := path.Join(backupBasePath, projectName)

		info, err := projectInfo(backupBasePath, projectName)
		if err != nil {
			panic(err)
		}

		if info.Pid > 0 {
			b, _ := exec.Command("ps", fmt.Sprint(info.Pid)).Output()
			if strings.Contains(string(b), fmt.Sprint(info.Pid)) {
				fmt.Println("backup already running")
				return
			}
		}

		backupArgs := []string{"backup"}
		backupArgs = append(backupArgs, backupFullPath, projectPath)
		for _, ignore := range ignorePaths {
			backupArgs = append(backupArgs, "--ignore_path", ignore)
		}
		backupArgs = append(backupArgs, "--last_updated", info.LastUpdated.Format(time.RFC3339))
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

func init() {
	dir, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	startCmd.Flags().String("name", path.Base(dir), "the name of the backup project")
	startCmd.Flags().String("path", path.Join(home, ".backups"), "the path to the backup project")
	startCmd.Flags().StringArray("ignore_path", []string{"node_modules", "dist", "build", ".git"}, "folders to ignore")
	rootCmd.AddCommand(startCmd)
}
