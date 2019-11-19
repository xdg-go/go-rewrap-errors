package testdata

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

const ErrStr = "error occurred"

func main() {
	err := fmt.Errorf("this is an error")
	log.Print(errors.Wrap(err, ErrStr))
}
