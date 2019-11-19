package testdata

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

const ErrStr = "error occurred '%s'"

func main() {
	err := fmt.Errorf("this is an error")
	foo := "foo"
	log.Print(errors.Wrapf(err, ErrStr, foo))
}
