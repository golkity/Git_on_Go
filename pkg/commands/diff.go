package commands

import (
	"fmt"
	"gogit/pkg/ignore"
	"gogit/pkg/objects"
	"gogit/pkg/refs"
	"gogit/pkg/ui"
	"os"
	"path/filepath"
	"strings"
)

func Diff() {
	headRef, _ := refs.GetHead()
	commitHash, _ := refs.GetCommitHash(headRef)
	if commitHash == "" {
		ui.Info("Нет коммитов для сравнения")
		return
	}

	_, data, _ := objects.ReadObject(commitHash)
	treeHash := strings.TrimPrefix(strings.Split(string(data), "\n")[0], "tree ")

	headIndex := make(map[string]string)
	loadTreeIndex(treeHash, ".", headIndex)

	ui.Info("Изменения (Working Dir vs HEAD):")
	hasChanges := false

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || ignore.ShouldIgnore(path) {
			return nil
		}

		currentHash, _ := objects.HashFile(path)

		relPath := "./" + path
		if strings.HasPrefix(path, ".") {
			relPath = path
		} else {
			relPath = "./" + path
		}
		relPath = filepath.Clean(relPath)

		storedHash, exists := headIndex[relPath]

		if !exists {
			fmt.Printf("%s[NEW] %s%s\n", ui.Green, path, ui.Reset)
			hasChanges = true
		} else if storedHash != currentHash {
			fmt.Printf("%s[MOD] %s%s\n", ui.Yellow, path, ui.Reset)
			hasChanges = true
		}

		delete(headIndex, relPath)
		return nil
	})

	for path := range headIndex {
		fmt.Printf("%s[DEL] %s%s\n", ui.Red, path, ui.Reset)
		hasChanges = true
	}

	if !hasChanges {
		ui.Success("Изменений нет, дерево чистое.")
	}
}

func loadTreeIndex(treeHash, dirCtx string, index map[string]string) {
	_, data, _ := objects.ReadObject(treeHash)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}
		meta := strings.Fields(parts[0])
		name := parts[1]

		objType := meta[1]
		hash := meta[2]

		fullPath := filepath.Join(dirCtx, name)

		if objType == "blob" {
			index[fullPath] = hash
		} else if objType == "tree" {
			loadTreeIndex(hash, fullPath, index)
		}
	}
}
