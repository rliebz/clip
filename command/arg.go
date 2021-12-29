package command

// WithArgs allows access to arbitrary arguments on a command.
//
// The slice of strings provided will be overridden with the non-flag command-
// line arguments.
func WithArgs(args *[]string) Option {
	return func(c *config) {
		c.argAction = func(ctx *Context) error {
			*args = ctx.args()

			return nil
		}
	}
}
