package testdata

import (
	"fmt"
	"log"
)

func main() {
	foo := "foo"
	log.Print(fmt.Errorf("error occurred '%s'", foo))
}
