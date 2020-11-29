package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func branchName(backupPath string) (string, error) {
	c := exec.Command("git", "branch")
	c.Dir = backupPath
	b, err := c.Output()
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(b), " ")
	return strings.TrimSpace(parts[len(parts)-1]), nil
}

var restoreCmd = &cobra.Command{
	Use:  "restore",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		backupBasePath, _ := cmd.Flags().GetString("path")
		outBasePath, _ := cmd.Flags().GetString("out")
		backupPath := filepath.Join(backupBasePath, name)
		outPath := filepath.Join(outBasePath, name)
		if !fileExists(backupPath) {
			fmt.Println(fmt.Sprintf("project %v does not exist in path %s", name, backupBasePath))
			os.Exit(0)
		}
		currentBranch, err := branchName(backupPath)
		if err != nil {
			panic(err)
		}

		c := exec.Command("git", "log", "--format=oneline")
		c.Dir = backupPath
		b, err := c.Output()
		if err != nil {
			panic(err)
		}

		commits := strings.Split(string(b), "\n")
		messages := []string{}
		hashes := map[string]string{}
		for _, commit := range commits {
			parts := strings.Split(commit, " ")
			if len(parts) == 2 {
				hash := parts[0]
				message := parts[1]
				messages = append(messages, message)
				hashes[message] = hash
			}
		}
		prompt := promptui.Select{
			Label: "Select time",
			Items: messages,
		}
		_, result, err := prompt.Run()
		checkout := hashes[result]

		c = exec.Command("git", "checkout", checkout)
		c.Dir = backupPath
		b, err = c.Output()
		if err != nil {
			panic(err)
		}
		filepath.Walk(backupPath, func(path string, info os.FileInfo, err error) error {
			relative, _ := filepath.Rel(backupPath, path)
			if info.IsDir() {
				if info.Name() == ".git" {
					return filepath.SkipDir
				}
				os.MkdirAll(filepath.Join(outPath, relative), 0755)
				return nil
			}
			copy(path, filepath.Join(outPath, relative))
			return nil
		})
		c = exec.Command("git", "checkout", currentBranch)
		c.Dir = backupPath
		b, err = c.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println("done! cd to ", outPath)

	},
}

func init() {
	dir, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	restoreCmd.Flags().String("path", path.Join(home, ".backups"), "the path to the backup project")
	restoreCmd.Flags().String("out", dir, "the path to the backup project")
	rootCmd.AddCommand(restoreCmd)
}
