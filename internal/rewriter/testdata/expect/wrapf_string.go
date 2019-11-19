package testdata

import (
	"fmt"
	"log"
)

func main() {
	err := fmt.Errorf("this is an error")
	foo := "foo"
	log.Print(fmt.Errorf("error occurred '%s': %w", foo, err))
}
