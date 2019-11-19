package testdata

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func ErrStr() string {
	return "error occurred"
}

func main() {
	err := fmt.Errorf("this is an error")
	log.Print(errors.Wrap(err, ErrStr()))
}
