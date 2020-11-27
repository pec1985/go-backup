package cmd

import (
	"fmt"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use: "stop",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("name")
		backupBasePath, _ := cmd.Flags().GetString("path")
		var info backupInfo
		var err error
		if info, err = projectInfo(backupBasePath, projectName); err != nil {
			panic(err)
		}
		if info.Pid != 0 {
			fmt.Println("killing", projectName)
			if err := syscall.Kill(info.Pid, syscall.SIGKILL); err != nil {
				panic(err)
			}
			info.Pid = 0
			info.LastUpdated = time.Now()
			saveProjectInfo(backupBasePath, args[0], info)
		}
	},
}

func init() {
	dir, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	stopCmd.Flags().String("name", path.Base(dir), "the name of the backup project")
	stopCmd.Flags().String("path", path.Join(home, ".backups"), "the path to the backup project")
	rootCmd.AddCommand(stopCmd)
}
