package testdata

import (
	"log"

	"github.com/pkg/errors"
)

func main() {
	foo := "foo"
	log.Print(errors.Errorf("error occurred '%s'", foo))
}
