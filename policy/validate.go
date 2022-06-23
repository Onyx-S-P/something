package policy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
)

func PullRequestEvent(payload []byte) {
	pr := PullRequest{}
	if err := json.Unmarshal(payload, &pr); err != nil {
		log.Println(err.Error())
		return
	}
	client, err := NewClient(pr.Installation.Id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if pr.Action == "opened" || pr.Action == "reopened" || pr.Action == "synchronize" {
		files := pr.ChangedFilesFromPullRequest(client)
		for _, f := range files {

			if filepath.Ext(f.GetFilename()) == ".yaml" && f.GetStatus() == "modified" {

				// if any change in base
				if f.GetDeletions() > 0 && f.GetAdditions() == 0 {

					body := "you have deleted something in base"
					pr.CreateReview(client, "COMMENT", body)
					pr.ClosePullRequest(client)

				} else if f.GetDeletions() > 0 && f.GetAdditions() > 0 {
					//TODO: further checkup with basef.GetAdditions() > 0 {
					body := "you have deleted and added something in base"

					pr.CreateReview(client, "REQUEST_CHANGES", body)

				} else if f.GetDeletions() == 0 && f.GetAdditions() > 0 {
					//Accept the changes from pull request
					log.Println("No change in base everything is fine...", pr.Number)
					pr.MergePullRequest(client)
					pr.CreateReview(client, "APPROVE", "")

				}
			}
		}
	}
}

func getDeletedLineNo(file *github.CommitFile) int {

	buff := bytes.NewBuffer([]byte(file.GetPatch()))
	scanner := bufio.NewScanner(buff)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		count++
		if strings.HasPrefix(line, "-") {
			return count
		}
	}
	return count
}
