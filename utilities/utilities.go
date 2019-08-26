package utilities

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// StringSliceIndexOf looks through the provided slice to see if
// the valueToMatch is in it or not
func StringSliceIndexOf(stringSlice []string, valueToMatch string) int {
	for i, val := range stringSlice {
		if val == valueToMatch {
			return i
		}
	}
	return -1
}

// GetBranchName will execute the git branch command and get the value
// indicated by the star as the branch in use and return that for use
func GetBranchName() string {
	cmd := exec.Command("git", "branch")
	var out, errBuff bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuff
	err := cmd.Run()

	result := strings.Split(out.String(), "\n")

	var branchName string

	for i, val := range result {
		found := strings.Index(val, "*")
		if found > -1 {
			selectedBranch := result[i]
			splitFromStar := strings.Split(selectedBranch, "*")
			branchName = strings.TrimSpace(splitFromStar[1])
		}
	}

	if err != nil {
		fmt.Println("Error while trying to get the branch name:")
		fmt.Println(errBuff.String())
		os.Exit(1)
	}
	return branchName
}

// RebaseBranch will rebase the current branch against the provided
// target branch on the specified origin. Exits with status 1 if
// the rebase fails and returns the message from git. If there
// are conflicts the user will have to pick it up form here to resolve
func RebaseBranch(rebaseTarget string, origin string) {

}

// ResetBranch will reset the current target against the
// target branch. It will call this after a rebase and will
// soft reset in order to prepare the changes for a single
// commit
func ResetBranch(resetTarget string, origin string) {

}

// PushCommit will use the branch name to determine the story number,
// construct a commit in the [changeType] (storyNumber): <description>
// format, and push all the changes up with it. If there is no upstream
// branch to push to, it will prompt the user and ask if it should create
// one with the --set-upstream <remote> <branchName>
func PushCommit(branchName string, changeType string, description string, remote string) {

}
