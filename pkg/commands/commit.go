package commands

import (
	"fmt"
	"gogit/pkg/ignore"
	"gogit/pkg/objects"
	"gogit/pkg/refs"
	"gogit/pkg/ui"
	"os"
	"path/filepath"
	"time"
)

func Commit(message string) {
	rootTreeHash, err := createTree(".")
	if err != nil {
		ui.Error("Ошибка построения дерева: %v", err)
		return
	}

	headRef, _ := refs.GetHead()
	parentHash, _ := refs.GetCommitHash(headRef)

	content := fmt.Sprintf("tree %s\n", rootTreeHash)
	if parentHash != "" {
		content += fmt.Sprintf("parent %s\n", parentHash)
	}
	ts := time.Now().Unix()
	content += fmt.Sprintf("author User <%d>\ncommitter User <%d>\n\n%s\n", ts, ts, message)

	commitHash, err := objects.SaveObject("commit", []byte(content))
	if err != nil {
		ui.Error("Ошибка записи: %v", err)
		return
	}

	refs.UpdateRef(headRef, commitHash)
	ui.Success("[%s] %s", headRef, commitHash[:7])
}

func createTree(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var buffer string
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if ignore.ShouldIgnore(path) {
			continue
		}

		var hash, mode, typeStr string
		if entry.IsDir() {
			hash, err = createTree(path)
			if hash == "" {
				continue
			}
			mode = "40000"
			typeStr = "tree"
		} else {
			data, _ := os.ReadFile(path)
			hash, err = objects.SaveObject("blob", data)
			mode = "100644"
			typeStr = "blob"
		}

		if err != nil {
			return "", err
		}
		buffer += fmt.Sprintf("%s %s %s\t%s\n", mode, typeStr, hash, entry.Name())
	}

	if buffer == "" {
		return "", nil
	}

	return objects.SaveObject("tree", []byte(buffer))
}
