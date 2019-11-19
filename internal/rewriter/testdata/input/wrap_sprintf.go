package testdata

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func main() {
	err := fmt.Errorf("this is an error")
	foo := "foo"
	log.Print(errors.Wrap(err, fmt.Sprintf("error occurred '%s'", foo)))
}
