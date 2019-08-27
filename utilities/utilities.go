package utilities

import (
	"bufio"
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
		if strings.EqualFold(val, valueToMatch) {
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

// StageCurrentChanges stages the current changes in git
func StageCurrentChanges() {
	fmt.Println("staging current changes...")
	cmd := exec.Command("git", "add", "--all")
	var outBuff, errBuff bytes.Buffer
	cmd.Stdout = &outBuff
	cmd.Stderr = &errBuff
	err := cmd.Run()

	if err != nil {
		fmt.Println("ERROR: Error occurred while staging the current changes: ")
		fmt.Println(errBuff.String())
		os.Exit(1)
	}
}

// RebaseBranch will rebase the current branch against the provided
// target branch on the specified origin. Exits with status 1 if
// the rebase fails and returns the message from git. If there
// are conflicts the user will have to pick it up form here to resolve
func RebaseBranch(rebaseTarget string, remote string) {
	fmt.Println("rebasing your current branch on " + remote + "/" + rebaseTarget)
	createCommit("tempCommit")
	cmd := exec.Command("git", "rebase", remote+"/"+rebaseTarget)
	var errBuff bytes.Buffer
	cmd.Stderr = &errBuff
	err := cmd.Run()

	if err != nil {
		fmt.Println("ERROR: error trying to rebase your branch, you may have conflicts: ")
		fmt.Println(errBuff.String())
		os.Exit(1)
	}
}

// ResetBranch will reset the current target against the
// target branch. It will call this after a rebase and will
// soft reset in order to prepare the changes for a single
// commit
func ResetBranch(resetTarget string, remote string) {
	fmt.Println("resetting branch against " + remote + "/" + resetTarget)
	cmd := exec.Command("git", "reset", remote+"/"+resetTarget)
	var errBuff bytes.Buffer
	cmd.Stderr = &errBuff
	err := cmd.Run()

	if err != nil {
		fmt.Println("ERROR: error while trying to reset your branch: ")
		fmt.Println(errBuff.String())
		os.Exit(1)
	}
}

// PushCommit will use the branch name to determine the story number,
// construct a commit in the [changeType] (storyNumber): <description>
// format, and push all the changes up with it. If there is no upstream
// branch to push to, it will prompt the user and ask if it should create
// one with the --set-upstream <remote> <branchName>
func PushCommit(branchName string, changeType string, description string, remote string) {
	fmt.Println("preparing to push commit....")
	if strings.Index(branchName, "/") == -1 {
		fmt.Println("Required branch naming not preset, please name the branch with a feature/<storyCardNumber> syntax")
		os.Exit(1)
	}
	branchParts := strings.Split(branchName, "/")
	storyNumber := branchParts[1]
	commitMessage := "[" + strings.ToLower(changeType) + "] " + "(" + storyNumber + "): " + description
	var pushBuff bytes.Buffer

	err := createCommit(commitMessage)

	if err != nil {
		fmt.Println("ERROR: error while trying to add the commit: ")
		os.Exit(1)
	}
	fmt.Println("pushing your commit with message: " + commitMessage)
	pushCmd := exec.Command("git", "push", "--force")
	pushCmd.Stderr = &pushBuff
	pushErr := pushCmd.Run()

	if pushErr != nil {
		errMsg := string(pushBuff.String())
		if strings.Index(errMsg, "upstream") != -1 {
			handleBranchUpstream(branchName, remote, commitMessage)
		}

		fmt.Println("ERROR: error while trying to push your commit")
		fmt.Println(errMsg)
		os.Exit(1)
	}

	fmt.Println("Succesfully pushed your commit with the message: " + commitMessage)
	fmt.Println("exiting..")
	os.Exit(0)

}

// RunFetch will fetch the changes from the remote
func RunFetch() {
	cmd := exec.Command("git", "fetch")
	var errBuff bytes.Buffer
	cmd.Stderr = &errBuff
	err := cmd.Run()

	if err != nil {
		fmt.Println("ERROR: error while trying to run fetch")
		fmt.Println(errBuff.String())
		os.Exit(1)
	}
}

func handleBranchUpstream(branchName string, remote string, commitMessage string) error {
	inputBuf := bufio.NewReader(os.Stdin)
	fmt.Println("no upstream for branch " + branchName + " on remote " + remote + ", create one? y/n")
	response, err := inputBuf.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error trying to parse answer")
		fmt.Println(err.Error())
	}
	responseString := strings.TrimSuffix(string(response), "\n")
	if !strings.EqualFold(responseString, "y") && !strings.EqualFold(responseString, "n") {
		recursErr := handleBranchUpstream(branchName, remote, commitMessage)
		return recursErr
	}

	if strings.EqualFold(responseString, "y") {
		cmd := exec.Command("git", "push", "--set-upstream", remote, branchName)
		var errBuffer bytes.Buffer
		cmd.Stderr = &errBuffer
		err := cmd.Run()

		if err != nil {
			fmt.Println("ERROR: error while trying to push and set the upstream: ")
			fmt.Println(errBuffer.String())
			os.Exit(1)
		}

		fmt.Println("Success, the commit was pushed with the message: " + commitMessage)
		os.Exit(0)
	} else {
		fmt.Println("Cannot push the branch without an upstream set, exiting..")
		os.Exit(1)
	}

	return nil
}

func createCommit(commitMessage string) error {
	commitCmd := exec.Command("git", "commit", "-m", commitMessage)
	var commitOut, commitErrBuff bytes.Buffer
	commitCmd.Stderr = &commitErrBuff
	commitCmd.Stdout = &commitOut
	err := commitCmd.Run()

	if strings.Index(commitOut.String(), "nothing to commit") != -1 {
		return nil
	}

	if err != nil {
		fmt.Println("ERROR: error while trying to add the commit: ")
		fmt.Println(commitErrBuff.String())
		os.Exit(1)
	}

	return nil
}
