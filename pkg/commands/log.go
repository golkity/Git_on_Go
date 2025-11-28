package commands

import (
	"fmt"
	"gogit/pkg/objects"
	"gogit/pkg/refs"
	"gogit/pkg/ui"
	"strings"
)

func Log() {
	headRef, err := refs.GetHead()
	if err != nil {
		ui.Error("Ошибка: %v", err)
		return
	}
	currentHash, _ := refs.GetCommitHash(headRef)

	ui.Info("История ветки: %s", headRef)
	fmt.Println()

	for currentHash != "" {
		_, data, err := objects.ReadObject(currentHash)
		if err != nil {
			break
		}
		content := string(data)
		lines := strings.Split(content, "\n")

		var parent string
		var msg string

		for i, line := range lines {
			if strings.HasPrefix(line, "parent ") {
				parent = strings.TrimPrefix(line, "parent ")
			}
			if line == "" && i < len(lines)-1 {
				msg = strings.Join(lines[i+1:], "\n")
				break
			}
		}

		fmt.Printf("%sCommit: %s%s\n", ui.Yellow, currentHash, ui.Reset)
		if msg != "" {
			fmt.Printf("    %s\n", strings.TrimSpace(msg))
		}
		fmt.Println(ui.Dim + "    |" + ui.Reset)

		currentHash = parent
	}
	fmt.Println(ui.Dim + "    (root)" + ui.Reset)
}
