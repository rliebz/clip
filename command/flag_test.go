package command

import (
	"fmt"
	"testing"

	"github.com/rliebz/clip/flag"
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func TestParse(t *testing.T) {
	var testCases = []struct {
		args         []string
		childSFlag   string
		parentSFlag  string
		childFlag    bool
		parentFlag   bool
		childCalled  bool
		parentCalled bool
	}{
		{
			args:         []string{"foo"},
			parentCalled: true,
		},
		{
			args:         []string{"foo", "-f"},
			parentFlag:   true,
			parentCalled: true,
		},
		{
			args:         []string{"foo", "--flag"},
			parentFlag:   true,
			parentCalled: true,
		},
		{
			args:         []string{"foo", "--sflag", "bar"},
			parentSFlag:  "bar",
			parentCalled: true,
		},
		{
			args:         []string{"foo", "-fs", "bar"},
			parentFlag:   true,
			parentSFlag:  "bar",
			parentCalled: true,
		},
		{
			args:         []string{"foo", "--sflag", "bar", "--flag"},
			parentFlag:   true,
			parentSFlag:  "bar",
			parentCalled: true,
		},
		{
			args:         []string{"foo", "--flag=true"},
			parentFlag:   true,
			parentCalled: true,
		},
		{
			args:        []string{"foo", "child"},
			childCalled: true,
		},
		{
			args:        []string{"foo", "child", "--flag"},
			childFlag:   true,
			childCalled: true,
		},
		{
			args:        []string{"foo", "--flag", "child"},
			childCalled: true,
			parentFlag:  true,
		},
		{
			args:        []string{"foo", "--flag=true", "child"},
			childCalled: true,
			parentFlag:  true,
		},
		{
			args:        []string{"foo", "--flag", "child", "--flag"},
			childFlag:   true,
			childCalled: true,
			parentFlag:  true,
		},
		{
			args:        []string{"foo", "-fs=bar", "child"},
			childCalled: true,
			parentFlag:  true,
			parentSFlag: "bar",
		},
		{
			args:        []string{"foo", "-f", "-s=bar", "child"},
			childCalled: true,
			parentFlag:  true,
			parentSFlag: "bar",
		},
		{
			args:        []string{"foo", "--sflag", "bar", "child"},
			childCalled: true,
			parentSFlag: "bar",
		},
		{
			args:        []string{"foo", "-fs", "bar", "child", "-fs", "baz"},
			childFlag:   true,
			childSFlag:  "baz",
			childCalled: true,
			parentFlag:  true,
			parentSFlag: "bar",
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("args: %v", tt.args), func(t *testing.T) {
			childFlag := false
			childSFlag := ""
			childCalled := false
			child := New(
				"child",
				WithFlag(flag.NewString(&childSFlag, "sflag", flag.WithShort("s"))),
				WithFlag(flag.NewBool(&childFlag, "flag", flag.WithShort("f"))),
				WithAction(func(ctx *Context) error {
					childCalled = true
					return nil
				}),
			)

			parentFlag := false
			parentSFlag := ""
			parentCalled := false
			cmd := New(
				"foo",
				WithFlag(flag.NewString(&parentSFlag, "sflag", flag.WithShort("s"))),
				WithFlag(flag.NewBool(&parentFlag, "flag", flag.WithShort("f"))),
				WithCommand(child),
				WithAction(func(ctx *Context) error {
					parentCalled = true
					return nil
				}),
			)

			assert.NilError(t, cmd.Execute(tt.args))

			assert.Check(t, cmp.Equal(childCalled, tt.childCalled))
			assert.Check(t, cmp.Equal(childFlag, tt.childFlag))
			assert.Check(t, cmp.Equal(childSFlag, tt.childSFlag))
			assert.Check(t, cmp.Equal(parentCalled, tt.parentCalled))
			assert.Check(t, cmp.Equal(parentFlag, tt.parentFlag))
			assert.Check(t, cmp.Equal(parentSFlag, tt.parentSFlag))
		})
	}

}

func TestParseError(t *testing.T) {
	var testCases = []struct {
		args []string
		err  string
	}{
		{
			args: []string{"foo", "--bad"},
			err:  "unknown flag: bad",
		},
		{
			args: []string{"foo", "-b"},
			err:  "unknown shorthand flag: 'b' in -b",
		},
		{
			args: []string{"foo", "-bad"},
			err:  "unknown shorthand flag: 'd' in -bad",
		},
		{
			args: []string{"foo", "bad"},
			err:  "undefined sub-command: bad",
		},
		{
			args: []string{"foo", "child", "--bad"},
			err:  "unknown flag: bad",
		},
		{
			args: []string{"foo", "child", "-b"},
			err:  "unknown shorthand flag: 'b' in -b",
		},
		{
			args: []string{"foo", "child", "-bad"},
			err:  "unknown shorthand flag: 'd' in -bad",
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("args: %v", tt.args), func(t *testing.T) {
			childFlag := false
			childCalled := false
			child := New(
				"child",
				WithFlag(flag.NewBool(&childFlag, "flag")),
				WithAction(func(ctx *Context) error {
					childCalled = true
					return nil
				}),
			)

			parentFlag := false
			parentCalled := false
			cmd := New(
				"foo",
				WithFlag(flag.NewBool(&parentFlag, "flag")),
				WithCommand(child),
				WithAction(func(ctx *Context) error {
					parentCalled = true
					return nil
				}),
			)

			assert.Error(t, cmd.Execute(tt.args), tt.err)
			assert.Check(t, cmp.Equal(childCalled, false))
			assert.Check(t, cmp.Equal(parentCalled, false))
		})
	}
}
