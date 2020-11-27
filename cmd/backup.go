package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func contains(what string, in []string) bool {
	for _, i := range in {
		if what == i {
			return true
		}
	}
	return false
}

var backupCmd = &cobra.Command{
	Use:  "backup <backup_path, project_path>",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupPath, _ := filepath.Abs(args[0])
		projectPath, _ := filepath.Abs(args[1])
		ignorePaths, _ := cmd.Flags().GetStringArray("ignore_paths")
		lastUpdtedString, _ := cmd.Flags().GetString("last_updated")
		var lastUpdated time.Time
		if lastUpdtedString != "" {
			lastUpdated, _ = time.Parse(time.RFC3339, lastUpdtedString)
		}

		if !fileExists(backupPath) {
			if err := os.MkdirAll(backupPath, 0755); err != nil {
				panic(err)
			}
		}
		if !fileExists(path.Join(backupPath, ".git")) {
			if err := exec.Command("git", "init").Run(); err != nil {
				panic(err)
			}
		}

		for {
			if err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() && contains(info.Name(), ignorePaths) {
					return filepath.SkipDir
				}
				relative, _ := filepath.Rel(projectPath, path)
				dest := filepath.Join(backupPath, relative)
				if info.IsDir() && !fileExists(dest) {
					if err := os.MkdirAll(dest, 0755); err != nil {
						panic(err)
					}
					return nil
				}
				if info.ModTime().After(lastUpdated) {
					if err := copy(path, dest); err != nil {
						panic(err)
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
			lastUpdated = time.Now()
			c := exec.Command("git", "add", "--all")
			c.Dir = backupPath
			c.Run()
			c = exec.Command("git", "commit", "-m", lastUpdated.Format(time.RFC3339))
			c.Dir = backupPath
			c.Run()
			time.Sleep(time.Minute * 2)
		}
	},
}

func ticker() {

}

func copy(sourceFile, destinationFile string) error {

	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(destinationFile, input, 0644)
}

func init() {
	backupCmd.Flags().StringArray("ignore_paths", []string{}, "")
	backupCmd.Flags().String("last_updated", "", "")
	rootCmd.AddCommand(backupCmd)
}
