package commands

import (
	"gogit/pkg/objects"
	"gogit/pkg/refs"
	"gogit/pkg/ui"
	"os"
	"path/filepath"
	"strings"
)

func Checkout(target string) {
	commitHash, err := refs.GetCommitHash("refs/heads/" + target)
	if commitHash == "" {
		commitHash = target
	}

	_, data, err := objects.ReadObject(commitHash)
	if err != nil {
		ui.Error("Коммит/Ветка '%s' не найдены", target)
		return
	}

	lines := strings.Split(string(data), "\n")
	var treeHash string
	if strings.HasPrefix(lines[0], "tree ") {
		treeHash = strings.TrimPrefix(lines[0], "tree ")
	}

	if treeHash == "" {
		ui.Error("Невалидный объект коммита")
		return
	}

	ui.Step("Переключение на " + target + "...")

	if err := restoreTree(treeHash, "."); err != nil {
		ui.Error("Ошибка восстановления: %v", err)
	} else {
		headPath := filepath.Join(objects.RootDir, "HEAD")
		isBranch := false
		if _, err := os.Stat(filepath.Join(objects.RootDir, "refs", "heads", target)); err == nil {
			os.WriteFile(headPath, []byte("ref: refs/heads/"+target), 0644)
			isBranch = true
		} else {
			os.WriteFile(headPath, []byte(commitHash), 0644)
		}

		state := "Detached HEAD"
		if isBranch {
			state = "Branch " + target
		}
		ui.Success("Рабочая директория обновлена. HEAD: %s", state)
	}
}

func restoreTree(treeHash, dirPath string) error {
	_, data, err := objects.ReadObject(treeHash)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		meta := strings.Fields(parts[0]) // [mode, type, hash]
		name := parts[1]

		if len(meta) < 3 {
			continue
		}
		objType := meta[1]
		objHash := meta[2]

		fullPath := filepath.Join(dirPath, name)

		if objType == "blob" {
			_, content, err := objects.ReadObject(objHash)
			if err != nil {
				return err
			}
			os.WriteFile(fullPath, content, 0644)
		} else if objType == "tree" {
			os.MkdirAll(fullPath, 0755)
			restoreTree(objHash, fullPath)
		}
	}
	return nil
}
