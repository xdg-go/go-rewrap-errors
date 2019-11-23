package rewriter

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

func isEqual(t *testing.T, got, expect, label string) {
	t.Helper()
	if got != expect {
		t.Errorf("%s:\n\n%s", label, diff.LineDiff(expect, got))
	}
}

func TestRewrite(t *testing.T) {

	cases := []string{
		"wrap_string",
		"wrap_fcn",
		"wrap_var",
		"wrap_sprintf",
		"wrap_sprintf_var",
		"wrapf_string",
		"wrapf_fcn",
		"wrapf_var",
		"wrapf_sprintf",
		"errorf_string",
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
			isEqual(t, so, se, fmt.Sprintf("Incorrect rewrite of %s", c))
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

func TestSubstPkg(t *testing.T) {
	in := "package main\n\nimport \"github.com/pkg/errors\"\n"
	out, err := Rewrite("test", []byte(in))
	if err != nil {
		t.Errorf("Rewrite error: %v", err)
	}
	expect := "package main\n\nimport \"errors\"\n"
	isEqual(t, string(out), expect, "Incorrect rewrite of single import")

	in = "package main\n\nimport (\n\t\"bytes\"\n\t\"fmt\"\n\n\t\"github.com/pkg/errors\"\n)\n"
	out, err = Rewrite("test", []byte(in))
	if err != nil {
		t.Errorf("Rewrite error: %v", err)
	}
	expect = "package main\n\nimport (\n\t\"bytes\"\n\t\"errors\"\n\t\"fmt\"\n)\n"
	isEqual(t, string(out), expect, "Incorrect rewrite of multi import")
}
