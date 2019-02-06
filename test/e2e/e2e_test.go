package e2e

import (
	"os"
	"strings"
	"testing"

	"git-ghost/test/util"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	result := m.Run()
	os.Exit(result)
}

func TestAll(t *testing.T) {
	ghostDir, err := util.CreateGitWorkDir()
	if err != nil {
		t.Fatal(err)
	}
	defer ghostDir.Remove()

	t.Run("BasicScenario", CreateTestBasicScenario(ghostDir))
}

func CreateTestBasicScenario(ghostDir *util.WorkDir) func(t *testing.T) {
	return func(t *testing.T) {
		srcDir, err := util.CreateGitWorkDir()
		if err != nil {
			t.Fatal(err)
		}
		defer srcDir.Remove()

		err = setupBasicGitRepo(srcDir)
		if err != nil {
			t.Fatal(err)
		}
		srcDir.Env = map[string]string{
			"GHOST_REPO": ghostDir.Dir,
		}

		dstDir, err := util.CloneWorkDir(srcDir)
		if err != nil {
			t.Fatal(err)
		}
		dstDir.Env = map[string]string{
			"GHOST_REPO": ghostDir.Dir,
		}
		defer dstDir.Remove()

		// Make one modification
		_, _, err = srcDir.RunCommmand("bash", "-c", "echo c > sample.txt")
		if err != nil {
			t.Fatal(err)
		}

		stdout, _, err := srcDir.RunCommmand("git", "ghost", "push")
		if err != nil {
			t.Fatal(err)
		}
		diffHash := strings.TrimRight(stdout, "\n")
		assert.NotEqual(t, "", diffHash)

		_, _, err = srcDir.RunCommmand("git", "ghost", "show", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		// TODO: Do some assertion

		_, _, err = dstDir.RunCommmand("git", "ghost", "pull", diffHash)
		if err != nil {
			t.Fatal(err)
		}
		stdout, _, err = dstDir.RunCommmand("cat", "sample.txt")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "c\n", stdout)

		_, _, err = dstDir.RunCommmand("git", "ghost", "list")
		if err != nil {
			t.Fatal(err)
		}
		// TODO: Do some assertion

		// TODO: delete the ghost branches and do some assertion
	}
}

func setupBasicGitRepo(wd *util.WorkDir) error {
	var err error
	_, _, err = wd.RunCommmand("bash", "-c", "echo a > sample.txt")
	if err != nil {
		return err
	}
	_, _, err = wd.RunCommmand("git", "add", "sample.txt")
	if err != nil {
		return err
	}
	_, _, err = wd.RunCommmand("git", "commit", "sample.txt", "-m", "initial commit")
	if err != nil {
		return err
	}
	_, _, err = wd.RunCommmand("bash", "-c", "echo b > sample.txt")
	if err != nil {
		return err
	}
	_, _, err = wd.RunCommmand("git", "commit", "sample.txt", "-m", "second commit")
	if err != nil {
		return err
	}
	return nil
}