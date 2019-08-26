package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirgrantt/gitnomo/utilities"
)

func main() {
	changeType := flag.String("changeType", "", "bug|patch|hotfix|minor|major|release|test")
	allowedChangeTypes := []string{"bug", "patch", "hotfix", "minor", "major", "release", "test"}
	description := flag.String("description", "", "The description that will show up in the commit history")
	rebaseTarget := flag.String("rebaseTarget", "dev", "the target branch name to rebase against [default: dev]")
	resetTarget := flag.String("resetTarget", "dev", "the branch to reset against [default: dev]")
	remote := flag.String("remote", "origin", "the remote to use for the branches [default: origin]")
	flag.Parse()

	if strings.TrimSpace(*changeType) == "" {
		fmt.Println("ERROR: must provide a valid change type, see --help for options")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if indx := utilities.StringSliceIndexOf(allowedChangeTypes, *changeType); indx == -1 {
		fmt.Println("ERROR: Invalid commit change type, see --help for valid options")
		os.Exit(1)
	}

	if strings.TrimSpace(*description) == "" {
		fmt.Println(("ERROR: must provide a non empty description for the commit"))
		os.Exit(1)
	}

	if strings.TrimSpace(*rebaseTarget) == "" {
		fmt.Println("ERROR: rebaseTarget cannot be empty")
		os.Exit(1)
	}

	if strings.TrimSpace(*resetTarget) == "" {
		fmt.Println("ERROR: resetTarget cannot be empty")
		os.Exit(1)
	}

	if strings.TrimSpace(*remote) == "" {
		fmt.Println("ERROR: remote cannot be empty")
		os.Exit(1)
	}

	// TODO: Run git add --all before the rest because git won't like
	// performing the below actions with unstaged changes
	utilities.StageCurrentChanges()
	branchName := utilities.GetBranchName()
	utilities.RebaseBranch(*rebaseTarget, *remote)
	utilities.ResetBranch(*resetTarget, *remote)
	utilities.StageCurrentChanges()
	utilities.PushCommit(branchName, *changeType, *description, *remote)
	fmt.Println(branchName, *changeType)
}
