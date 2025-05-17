package command

import (
	"fmt"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"

	"github.com/rliebz/clip"
)

func TestParse(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(fmt.Sprintf("args: %v", tt.args), func(t *testing.T) {
			g := ghost.New(t)

			childFlag := false
			childSFlag := ""
			childCalled := false
			child := New(
				"child",
				WithFlag(clip.NewString(&childSFlag, "sflag", clip.FlagShort("s"))),
				WithFlag(clip.NewBool(&childFlag, "flag", clip.FlagShort("f"))),
				WithAction(func(*Context) error {
					childCalled = true
					return nil
				}),
			)

			parentFlag := false
			parentSFlag := ""
			parentCalled := false
			cmd := New(
				"foo",
				WithFlag(clip.NewString(&parentSFlag, "sflag", clip.FlagShort("s"))),
				WithFlag(clip.NewBool(&parentFlag, "flag", clip.FlagShort("f"))),
				WithCommand(child),
				WithAction(func(*Context) error {
					parentCalled = true
					return nil
				}),
			)

			g.NoError(cmd.Execute(tt.args))

			g.Should(be.Equal(childCalled, tt.childCalled))
			g.Should(be.Equal(childFlag, tt.childFlag))
			g.Should(be.Equal(childSFlag, tt.childSFlag))
			g.Should(be.Equal(parentCalled, tt.parentCalled))
			g.Should(be.Equal(parentFlag, tt.parentFlag))
			g.Should(be.Equal(parentSFlag, tt.parentSFlag))
		})
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(fmt.Sprintf("args: %v", tt.args), func(t *testing.T) {
			g := ghost.New(t)

			childFlag := false
			childCalled := false
			child := New(
				"child",
				WithFlag(clip.NewBool(&childFlag, "flag")),
				WithAction(func(*Context) error {
					childCalled = true
					return nil
				}),
			)

			parentFlag := false
			parentCalled := false
			cmd := New(
				"foo",
				WithFlag(clip.NewBool(&parentFlag, "flag")),
				WithCommand(child),
				WithAction(func(*Context) error {
					parentCalled = true
					return nil
				}),
			)

			g.Should(be.ErrorEqual(cmd.Execute(tt.args), tt.err))
			g.Should(be.False(childCalled))
			g.Should(be.False(parentCalled))
		})
	}
}
