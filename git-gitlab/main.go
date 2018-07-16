package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: git gitlab [subcommand]

Subcommand:
    help          -- Show this message.
    merged-branch -- List source branch name of merged MRs.`)
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "help" {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	switch cmd {
	case "merged-branch":
		err := mergedBranch()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand %q\n", cmd)
		os.Exit(2)
	}
}

func mergedBranch() error {
	url := os.Getenv("GITLAB_URL")
	if url == "" {
		return fmt.Errorf("GITLAB_URL is not set.")
	}
	repo := os.Getenv("GITLAB_REPO")
	if repo == "" {
		return fmt.Errorf("GITLAB_REPO is not set.")
	}
	token := os.Getenv("GITLAB_API_TOKEN")
	if token == "" {
		return fmt.Errorf("GITLAB_API_TOKEN is not set.")
	}
	info := gitlabInfo{url + "/api/v4", repo, token}

	mrs, err := getMergedRequests(info)
	if err != nil {
		return err
	}
	for _, mr := range mrs {
		fmt.Printf("!%v %v (%v -> %v)\n", mr.id, mr.title, mr.sourceBranch, mr.targetBranch)
	}

	branches, err := getRemoteBranches("origin")
	if err != nil {
		return err
	}
	fmt.Println(branches)

	return nil
}

type gitlabInfo struct {
	url   string
	repo  string
	token string
}

type mergeRequest struct {
	id           int
	title        string
	sourceBranch string
	targetBranch string
}

func getMergedRequests(info gitlabInfo) ([]mergeRequest, error) {
	git := gitlab.NewClient(nil, info.token)
	git.SetBaseURL(info.url)

	opt := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String("merged"),
	}
	mrs, _, err := git.MergeRequests.ListProjectMergeRequests(info.repo, opt)
	if err != nil {
		return nil, err
	}
	var mergeRequests []mergeRequest
	for _, mr := range mrs {
		mergeRequests = append(mergeRequests, mergeRequest{
			mr.IID, mr.Title, mr.SourceBranch, mr.TargetBranch,
		})
	}
	return mergeRequests, nil
}

func getRemoteBranches(remote string) ([]string, error) {
	out, err := exec.Command("git", "branch", "-a").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")

	var branches []string
	prefix := fmt.Sprintf("remotes/%s/", remote)
	for _, line := range lines {
		s := strings.Trim(line, " ")
		if strings.Index(s, prefix) == 0 {
			branches = append(branches, strings.TrimPrefix(s, prefix))
		}
	}
	return branches, nil
}
