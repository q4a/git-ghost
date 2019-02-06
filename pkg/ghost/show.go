package ghost

import (
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"
	"io"

	log "github.com/Sirupsen/logrus"
)

// ShowOptions represents arg for Pull func
type ShowOptions struct {
	types.WorkingEnvSpec
	*types.LocalBaseBranchSpec
	*types.PullableLocalModBranchSpec
	// if you want to consume and transform the output of `ghost.Show()`,
	// Please use `io.Pipe()` as below,
	// ```
	// r, w := io.Pipe()
	// go func() { ghost.Show(ShowOptions{ Writer: w }); w.Close()}
	// ````
	// Then, you can read the output from `r` and transform them as you like.
	Writer io.Writer
}

func pullAndshow(branchSpec types.PullableGhostBranchSpec, we types.WorkingEnv, writer io.Writer) error {
	branch, err := branchSpec.PullBranch(we)
	if err != nil {
		return err
	}
	if branch != nil {
		return branch.Show(we, writer)
	}
	return nil
}

// Show writes ghost branches contents to option.Writer
func Show(options ShowOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")

	if options.LocalBaseBranchSpec != nil {
		we, err := options.WorkingEnvSpec.Initialize()
		if err != nil {
			return err
		}
		defer we.Clean()
		err = pullAndshow(options.LocalBaseBranchSpec, *we, options.Writer)
		if err != nil {
			return err
		}
	}

	if options.PullableLocalModBranchSpec != nil {
		we, err := options.WorkingEnvSpec.Initialize()
		if err != nil {
			return err
		}
		defer we.Clean()
		return pullAndshow(options.PullableLocalModBranchSpec, *we, options.Writer)
	}

	log.WithFields(util.ToFields(options)).Warn("show command has nothing to do with")
	return nil
}