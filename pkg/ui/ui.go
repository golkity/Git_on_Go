package ui

import "fmt"

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
)

func Header() {
	fmt.Print("\033[2J\033[H") //clr term
	logo := `
   ________  ________  ________  ___  _________   
  |\   ____\|\   __  \|\   ____\|\  \|\___   ___\ 
  \ \  \___|\ \  \|\  \ \  \___|\ \  \|___ \  \_| 
   \ \  \  __\ \  \\\  \ \  \  __\ \  \   \ \  \  
    \ \  \|\  \ \  \\\  \ \  \|\  \ \  \   \ \  \ 
     \ \_______\ \_______\ \_______\ \__\   \ \__\
      \|_______|\|_______|\|_______|\|__|    \|__| v1.0
`
	fmt.Println(Cyan + logo + Reset)
	fmt.Printf("   %sRunning on Go 1.25.3%s\n\n", Dim, Reset)
}

func Step(text string) {
	fmt.Printf("%s==>%s %s%s%s\n", Blue, Reset, Bold, text, Reset)
}

func Success(format string, a ...any) {
	fmt.Printf(" %s✔%s %s\n", Green, Reset, fmt.Sprintf(format, a...))
}

func Error(format string, a ...any) {
	fmt.Printf(" %s✘%s %s\n", Red, Reset, fmt.Sprintf(format, a...))
}

func Info(format string, a ...any) {
	fmt.Printf(" %sℹ%s %s\n", Magenta, Reset, fmt.Sprintf(format, a...))
}

func TableRow(cmd, args, desc string) {
	fmt.Printf("  %s%-10s%s %-20s %s%s%s\n",
		Yellow, cmd, Reset,
		Dim+args+Reset,
		Reset, desc, Reset)
}
