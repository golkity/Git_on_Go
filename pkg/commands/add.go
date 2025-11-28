package commands

import (
	"gogit/pkg/ignore"
	"gogit/pkg/objects"
	"gogit/pkg/ui"
	"os"
	"path/filepath"
	"sync"
)

type AddResult struct {
	Path string
	Hash string
	Err  error
}

func Add(targetDir string) {
	if _, err := os.Stat(objects.RootDir); os.IsNotExist(err) {
		ui.Error("Репозиторий не инициализирован")
		return
	}

	ui.Step("Индексация файлов...")

	jobs := make(chan string, 100)
	results := make(chan AddResult, 100)
	var wg sync.WaitGroup

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				data, err := os.ReadFile(path)
				if err != nil {
					results <- AddResult{Path: path, Err: err}
					continue
				}
				hash, err := objects.SaveObject("blob", data)
				results <- AddResult{Path: path, Hash: hash, Err: err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if ignore.ShouldIgnore(path) {
				return nil
			}
			jobs <- path
			return nil
		})
		close(jobs)
	}()

	count := 0
	for res := range results {
		if res.Err != nil {
			ui.Error("%s: %v", res.Path, res.Err)
		} else {
			count++
		}
	}
	ui.Success("Добавлено объектов: %d", count)
}
