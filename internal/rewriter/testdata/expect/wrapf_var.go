package testdata

import (
	"fmt"
	"log"
)

const ErrStr = "error occurred '%s'"

func main() {
	err := fmt.Errorf("this is an error")
	foo := "foo"
	log.Print(fmt.Errorf(ErrStr+": %w", foo, err))
}
