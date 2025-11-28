package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func ShouldIgnore(path string) bool {
	if path == ".gogit" || strings.HasPrefix(path, ".gogit/") {
		return true
	}
	if path == ".git" || strings.HasPrefix(path, ".git/") {
		return true
	}

	file, err := os.Open(".gitignore")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matched, _ := filepath.Match(line, filepath.Base(path))
		if matched {
			return true
		}
		if strings.HasSuffix(line, "/") && strings.Contains(path, line) {
			return true
		}
	}
	return false
}
