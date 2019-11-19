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
		"wrap_fcn",
		"wrap_var",
		"wrapf_string",
		"wrapf_fcn",
		"wrapf_var",
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
			if err != nil {
				t.Errorf("Rewrite error: %v", err)
			}
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
	if i < 0 {
		return in
	}
	return in[i:]
}
