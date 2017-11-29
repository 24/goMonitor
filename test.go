package main

import (
	"fmt"
)

type DividerError struct {
	dividee int
	divider int
}

func (de *DividerError) Error() string {
	strFormat := `
	Cannot procceed, the divider is zero.
	dividee %d
	divider:0
	`
	return fmt.Sprintf(strFormat, de.dividee)
}
func Divide(varDividee int, varDivider int) (result int, errMsg string) {
	if varDivider == 0 {
		dData := DividerError{
			dividee: varDividee,
			divider: varDivider,
		}
		errMsg = dData.Error()
		return
	} else {
		return varDividee / varDivider, ""
	}
}

func main() {
	if result, errorMsg := Divide(100, 20); errorMsg == "" {
		fmt.Println(result)
	}

	if _, errorMsg := Divide(100, 0); errorMsg != "" {
		fmt.Println(errorMsg)
	}
}
