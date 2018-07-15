package main

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: git gitlab [subcommand]

Subcommand:
    help          -- Show this message.
    merged-branch -- List source branch name of merged MRs.
`)
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
	token := os.Getenv("GITLAB_API_TOKEN")
	if token == "" {
		return fmt.Errorf("GITLAB_API_TOKEN is not set.")
	}
	git := gitlab.NewClient(nil, token)
	git.SetBaseURL("http://localhost:8080/api/v4")

	opt := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String("opened"),
	}
	mrs, _, err := git.MergeRequests.ListProjectMergeRequests("root/test", opt)
	if err != nil {
		return err
	}
	fmt.Printf("%v %v %v->%v", mrs[0].IID, mrs[0].Title, mrs[0].SourceBranch, mrs[0].TargetBranch)
	return nil
}
