package Interface

import (
	"fmt"
	"github.com/raphamorim/go-rainbow"
)

func NormalMsg(text string, variables []any) {
	fmt.Printf(fmt.Sprintf(text, variables...))
	fmt.Println()
}

func ErrorMsg(text string, variables []any) {
	fmt.Println()
	//fmt.Printf(rainbow.BgRed(rainbow.White(rainbow.Bold(fmt.Sprintf(text, variables...)))))
	fmt.Printf(rainbow.Red(rainbow.Bold(fmt.Sprintf(text, variables...))))
	fmt.Println()
}

func SuccessMsg(text string, variables []any) {
	fmt.Println()
	//fmt.Printf(rainbow.BgGreen(rainbow.Black(fmt.Sprintf(text, variables...))))
	fmt.Printf(rainbow.Green(fmt.Sprintf(text, variables...)))
	fmt.Println()
}

func WarningMsg(text string, variables []any) {
	fmt.Println()
	//fmt.Printf(rainbow.BgYellow(rainbow.Black(fmt.Sprintf(text, variables...))))
	fmt.Printf(rainbow.Yellow(fmt.Sprintf(text, variables...)))
	fmt.Println()
}

func InfoMsg(text string, variables []any) {
	fmt.Println()
	//fmt.Printf(rainbow.BgBlack(rainbow.Bold(fmt.Sprintf(text, variables...))))
	fmt.Printf(rainbow.Bold(fmt.Sprintf(text, variables...)))
	fmt.Println()
}

func StartMsgLine() {
	fmt.Println()
}

func EndMsgLine() {
	fmt.Println()
}

func PrintSimpleText(text string) {
	fmt.Print(text)
}

func PrintVariablesText(text string, variables []any) {
	fmt.Printf(fmt.Sprintf(text, variables...))
}

func PrintInfoText(text string) {
	fmt.Print(rainbow.Bold(" " + text))
}

func PrintSuccessText(text string) {
	fmt.Print(rainbow.Bold(rainbow.Green(" " + text)))
}

func PrintFailText(text string) {
	fmt.Print(rainbow.Bold(rainbow.Red(" " + text)))
}
