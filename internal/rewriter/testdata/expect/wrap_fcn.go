package testdata

import (
	"fmt"
	"log"
)

func ErrStr() string {
	return "error occurred"
}

func main() {
	err := fmt.Errorf("this is an error")
	log.Print(fmt.Errorf(ErrStr()+": %w", err))
}
