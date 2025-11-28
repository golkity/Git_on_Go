package main

import (
	"fmt"
	"gogit/pkg/commands"
	"gogit/pkg/ui"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "init":
		ui.Header()
		commands.Init()
		ui.Success("Готово! Используйте 'gogit add .' для старта.")

	case "add":
		path := "."
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		commands.Add(path)

	case "commit":
		if len(os.Args) < 3 {
			ui.Error("Где сообщение? Используйте: commit \"msg\"")
			return
		}
		commands.Commit(os.Args[2])

	case "branch":
		if len(os.Args) < 3 {
			ui.Error("Укажите имя ветки")
			return
		}
		commands.CreateBranch(os.Args[2])

	case "log":
		commands.Log()

	case "cat-file":
		if len(os.Args) < 3 {
			ui.Error("Укажите хеш")
			return
		}
		commands.CatFile(os.Args[2])

	case "checkout":
		if len(os.Args) < 3 {
			ui.Error("Укажите ветку или хеш")
			return
		}
		commands.Checkout(os.Args[2])

	case "diff":
		commands.Diff()

	case "help":
		showHelp()

	default:
		ui.Error("Команда '%s' не найдена", os.Args[1])
		showHelp()
	}
}

func showHelp() {
	ui.Header()
	fmt.Println(ui.Bold + " CORE COMMANDS" + ui.Reset)
	ui.TableRow("init", "", "Создать новый репозиторий")
	ui.TableRow("add", "<path>", "Индексировать файлы")
	ui.TableRow("commit", "<msg>", "Сделать commit")

	fmt.Println("\n" + ui.Bold + " HISTORY & STATE" + ui.Reset)
	ui.TableRow("log", "", "Показать историю коммитов")
	ui.TableRow("diff", "", "Сравнить workdir с HEAD")

	fmt.Println("\n" + ui.Bold + " BRANCHING" + ui.Reset)
	ui.TableRow("branch", "<name>", "Создать ветку")
	ui.TableRow("checkout", "<ref>", "Переключить ветку/коммит")

	fmt.Println("\n" + ui.Bold + " DEBUG" + ui.Reset)
	ui.TableRow("cat-file", "<hash>", "Показать содержимое объекта")
}
