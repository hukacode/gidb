package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	zglob "github.com/mattn/go-zglob"
	"github.com/mkideal/cli"
)

type void struct{}

var member void

type argT struct {
	cli.Helper
	Path     string `cli:"p,path" usage:"path to your Dropbox or a folder which contains .gitignore"`
	DryRun   bool   `cli:"d,dry-run" usage:"dry run"`
	UnIgnore bool   `cli:"u,unignore" usage:"unignore"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		if argv.Path == "" {
			argv.Path, _ = os.Getwd()
		}

		gitIgnoreFiles := getGitIgnoreFiles(argv.Path)
		allItems := getItems(gitIgnoreFiles)
		filteredItems := filter(allItems)

		for _, item := range filteredItems {
			if argv.DryRun {
				log.Println("Will be processed:", item)
			} else {
				var err error

				if argv.UnIgnore {
					err = include(item)
				} else {
					err = exclude(item)
				}

				if err != nil {
					log.Println(err)
				} else {
					log.Println("Processed:", item)
				}
			}
		}

		log.Println("Done")

		return nil
	}))

}

func getGitIgnoreFiles(root string) []string {
	gitIgnoreFiles := []string{}

	filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".gitignore" {
			gitIgnoreFiles = append(gitIgnoreFiles, filepath.ToSlash(path))
		}

		return nil
	})

	return gitIgnoreFiles
}

func getItems(gitIgnoreFiles []string) []string {
	set := make(map[string]void)

	for _, gitignoreFile := range gitIgnoreFiles {
		patterns := readGitIgnoreFileContent(gitignoreFile)
		gitIgnoreRoot := strings.TrimSuffix(gitignoreFile, ".gitignore")

		for _, pattern := range patterns {
			matches, _ := zglob.Glob(filepath.Join(gitIgnoreRoot, strings.Replace(pattern, "\\ ", " ", -1)))

			for _, match := range matches {
				set[match] = member
			}
		}
	}

	sortedItem := make([]string, 0)
	for item := range set {
		sortedItem = append(sortedItem, filepath.ToSlash(item))
	}
	sort.Strings(sortedItem)

	return sortedItem
}

func filter(allItems []string) []string {
	filteredItems := make([]string, 0)

	for _, item := range allItems {
		isExisted := false

		for _, existed := range filteredItems {
			if strings.HasPrefix(item, existed+"/") {
				isExisted = true
				break
			}
		}

		if !isExisted && len(item) > 0 {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func readGitIgnoreFileContent(path string) []string {
	f, _ := os.Open(path)
	defer f.Close()

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	var s []string
	for scan.Scan() {
		line := scan.Text()

		if len(line) == 0 || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}

		s = append(s, line)
	}
	return s
}

func exclude(path string) error {
	if runtime.GOOS == "windows" {
		return excludeInWindows(path)
	} else if runtime.GOOS == "linux" {
		return excludeInLinux(path)
	} else if runtime.GOOS == "darwin" {
		return excludeInMacOS(path)
	} else {
		return fmt.Errorf("Could not ignore: %v", path)
	}
}

func excludeInWindows(path string) error {
	command := fmt.Sprintf("Set-Content -Path '%s' -Stream com.dropbox.ignored -Value 1", path)
	_, err := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell", "-NoProfile", command).CombinedOutput()

	return err
}

func excludeInLinux(path string) error {
	command := fmt.Sprintf("attr -s com.dropbox.ignored -V 1 '%s'", path)
	_, err := exec.Command("bash", "-c", command).CombinedOutput()

	return err
}

func excludeInMacOS(path string) error {
	command := fmt.Sprintf("xattr -w com.dropbox.ignored 1 '%s'", path)
	_, err := exec.Command("bash", "-c", command).CombinedOutput()

	return err
}

func include(path string) error {
	if runtime.GOOS == "windows" {
		includeInWindows(path)
	} else if runtime.GOOS == "linux" {
		includeInLinux(path)
	} else if runtime.GOOS == "darwin" {
		includeInMacOS(path)
	} else {
		return fmt.Errorf("Could not ignore: %v", path)
	}

	return nil
}

func includeInWindows(path string) error {
	command := fmt.Sprintf("Clear-Content -Path '%s' -Stream com.dropbox.ignored", path)
	_, err := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell", "-NoProfile", command).CombinedOutput()

	return err
}

func includeInLinux(path string) error {
	command := fmt.Sprintf("attr -r com.dropbox.ignored '%s'", path)
	_, err := exec.Command("bash", "-c", command).CombinedOutput()

	return err
}

func includeInMacOS(path string) error {
	command := fmt.Sprintf("xattr -d com.dropbox.ignored '%s'", path)
	_, err := exec.Command("bash", "-c", command).CombinedOutput()

	return err
}
