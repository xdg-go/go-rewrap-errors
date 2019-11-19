package rewriter

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

func TestRewrite(t *testing.T) {

	cases := []string{
		"wrap_string",
		// "wrap_fcn",
		// "wrap_var",
		"wrapf_string",
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata/input", c) + ".go")
			if err != nil {
				t.Fatal(err)
			}
			expect, err := ioutil.ReadFile(filepath.Join("testdata/expect", c) + ".go")
			if err != nil {
				t.Fatal(err)
			}
			out, err := Rewrite(c, in)
			so := trimToMain(string(out))
			se := trimToMain(string(expect))
			if so != se {
				t.Errorf("Incorrect rewrite of %s:\n%v", c, diff.LineDiff(se, so))
			}
		})
	}
}

func trimToMain(in string) string {
	i := strings.Index(in, "func main")
	return in[i:]
}
