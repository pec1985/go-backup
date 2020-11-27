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
	Use:  "stop [name]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupBasePath, _ := cmd.Flags().GetString("path")
		var info backupInfo
		var err error
		if info, err = projectInfo(backupBasePath, args[0]); err != nil {
			panic(err)
		}
		if info.Pid != 0 {
			fmt.Println("killing", args[0])
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
	home, _ := os.UserHomeDir()
	stopCmd.Flags().String("path", path.Join(home, ".backups"), "the path to the backup project")
	rootCmd.AddCommand(stopCmd)
}
