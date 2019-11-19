package testdata

import (
	"fmt"
	"log"
)

func main() {
	err := fmt.Errorf("this is an error")
	log.Print(fmt.Errorf("error occurred: %w", err))
}
