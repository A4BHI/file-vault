package errorcheck

import "fmt"

func PrintError(statement string, err error) {
	if err != nil {
		fmt.Println(statement, err)
	}
}
