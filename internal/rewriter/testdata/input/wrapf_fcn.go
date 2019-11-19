package testdata

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func ErrStr() string {
	return "error occurred '%s'"
}

func main() {
	err := fmt.Errorf("this is an error")
	foo := "foo"
	log.Print(errors.Wrapf(err, ErrStr(), foo))
}
