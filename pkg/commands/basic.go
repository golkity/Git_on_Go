package commands

import (
	"fmt"
	"gogit/pkg/objects"
	"gogit/pkg/refs"
	"gogit/pkg/ui"
	"os"
	"path/filepath"
)

func Init() {
	dirs := []string{
		filepath.Join(objects.RootDir, "objects"),
		filepath.Join(objects.RootDir, "refs", "heads"),
	}
	for _, d := range dirs {
		os.MkdirAll(d, 0755)
	}
	headPath := filepath.Join(objects.RootDir, "HEAD")
	if _, err := os.Stat(headPath); os.IsNotExist(err) {
		os.WriteFile(headPath, []byte("ref: refs/heads/master\n"), 0644)
	}
}

func CreateBranch(name string) {
	headRef, err := refs.GetHead()
	if err != nil {
		ui.Error("Ошибка чтения HEAD")
		return
	}
	currentHash, _ := refs.GetCommitHash(headRef)
	if currentHash == "" {
		ui.Error("Нельзя создать ветку: нет коммитов")
		return
	}

	newRef := filepath.Join("refs", "heads", name)
	if err := refs.UpdateRef(newRef, currentHash); err != nil {
		ui.Error("Ошибка: %v", err)
	} else {
		ui.Success("Ветка '%s' создана на %s", name, currentHash[:7])
	}
}

func CatFile(hash string) {
	objType, data, err := objects.ReadObject(hash)
	if err != nil {
		ui.Error("Ошибка чтения объекта: %v", err)
		return
	}
	ui.Info("Тип: %s | Размер: %d байт", objType, len(data))
	fmt.Println("------------------------------------------------")
	if objType == "tree" {
		fmt.Print(string(data))
	} else {
		fmt.Print(string(data))
	}
	fmt.Println("\n------------------------------------------------")
}
