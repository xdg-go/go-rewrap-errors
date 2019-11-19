package testdata

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func main() {
	err := fmt.Errorf("this is an error")
	log.Print(errors.Wrap(err, "error occurred"))
}
