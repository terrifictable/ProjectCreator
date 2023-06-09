package common

import "fmt"

const (
	ANSI_RESET  = "\033[0m"
	ANSI_GREY   = "\033[38;5;242m"
	ANSI_HIGH   = "\033[38;5;159m"
	ANSI_RED    = "\033[38;5;162m"
	ANSI_GREEN  = "\033[38;5;49m"
	ANSI_BLUE   = "\033[38;5;81m"
	ANSI_ORANGE = "\033[38;5;215m"
	ANSI_YELLOW = "\033[38;5;226m"
)

var dbg = ""

func MSG(format string, a ...interface{}) {
	fmt.Printf(" "+format+"\n", a...)
}

func OK(format string, a ...interface{}) {
	fmt.Printf(" "+ANSI_GREY+"["+ANSI_GREEN+"+"+ANSI_GREY+"] "+ANSI_RESET+format+"\n", a...)
}

func INFO(format string, a ...interface{}) {
	fmt.Printf(" "+ANSI_GREY+"["+ANSI_BLUE+"*"+ANSI_GREY+"] "+ANSI_RESET+format+"\n", a...)
}

func WARN(format string, a ...interface{}) {
	fmt.Printf(" "+ANSI_GREY+"["+ANSI_ORANGE+"#"+ANSI_GREY+"] "+ANSI_RESET+format+"\n", a...)
}

func ERR(format string, a ...interface{}) {
	fmt.Printf(" "+ANSI_GREY+"["+ANSI_RED+"-"+ANSI_GREY+"] "+ANSI_RESET+format+"\n", a...)
}

func DBG(format string, a ...interface{}) {
	if dbg != "" {
		fmt.Printf(" "+ANSI_GREY+"["+ANSI_YELLOW+"~"+ANSI_GREY+"] "+ANSI_RESET+format+"\n", a...)
	}
}
