package testdata

import (
	"fmt"
	"log"
)

func ErrStr() string {
	return "error occurred '%s'"
}

func main() {
	err := fmt.Errorf("this is an error")
	foo := "foo"
	log.Print(fmt.Errorf(ErrStr()+": %w", foo, err))
}
