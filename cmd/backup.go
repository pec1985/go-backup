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

var backupCmd = &cobra.Command{
	Use:  "backup <backup_path, project_path>",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupPath, _ := filepath.Abs(args[0])
		projectPath, _ := filepath.Abs(args[1])
		ignorePaths, _ := cmd.Flags().GetStringArray("ignore_path")
		lastUpdatedString, _ := cmd.Flags().GetString("last_updated")
		var lastUpdated time.Time
		if lastUpdatedString != "" {
			lastUpdated, _ = time.Parse(time.RFC3339, lastUpdatedString)
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
			backedUpFiles := map[string]os.FileInfo{}
			projectFiles := map[string]os.FileInfo{}

			filepath.Walk(backupPath, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() && info.Name() == ".git" {
					return filepath.SkipDir
				}
				relative, _ := filepath.Rel(backupPath, path)
				backedUpFiles[relative] = info
				return nil
			})
			filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() && contains(info.Name(), ignorePaths) {
					return filepath.SkipDir
				}
				relative, _ := filepath.Rel(projectPath, path)
				projectFiles[relative] = info
				return nil
			})

			for relative, info := range projectFiles {
				dest := filepath.Join(backupPath, relative)
				if info.IsDir() {
					if !fileExists(dest) {
						if err := os.MkdirAll(dest, 0755); err != nil {
							panic(err)
						}
					}
					continue
				}
				if info.ModTime().After(lastUpdated) {
					from := filepath.Join(projectPath, relative)
					to := filepath.Join(backupPath, relative)
					if err := copy(from, to); err != nil {
						panic(err)
					}
				}
			}

			for relative := range backedUpFiles {
				if _, ok := projectFiles[relative]; !ok {
					if err := os.RemoveAll(filepath.Join(backupPath, relative)); err != nil {
						panic(err)
					}
				}
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
	backupCmd.Flags().StringArray("ignore_path", []string{}, "")
	backupCmd.Flags().String("last_updated", "", "")
	backupCmd.Hidden = true
	rootCmd.AddCommand(backupCmd)
}
