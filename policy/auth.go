package policy

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

func NewClient(installid int64) (*github.Client, error) {
	appId, _ := strconv.Atoi(os.Getenv("APP_ID"))
	privateKey := os.Getenv("PRIVATE_KEY")

	trans := http.DefaultTransport

	gh, err := ghinstallation.NewKeyFromFile(trans, int64(appId), installid, privateKey)

	if err != nil {
		return &github.Client{}, err
	}

	return github.NewClient(&http.Client{Transport: gh}), nil

}

/*Get all Changed files in github pull request*/
func (pr *PullRequest) ChangedFilesFromPullRequest(client *github.Client) []*github.CommitFile {

	files, _, err := client.PullRequests.ListFiles(context.Background(), pr.PullRequestData.Head.User.Login,
		pr.PullRequestData.Head.Repo.Name, int(pr.Number), &github.ListOptions{})

		fmt.Println("file: ",files)

	if err != nil {
		log.Println(err.Error())
		return []*github.CommitFile{}
	}
	return files
}

/*Merge the pull request*/
func (pr *PullRequest) MergePullRequest(client *github.Client) (*github.PullRequestMergeResult, error) {

	log.Printf("Merging pull request %d on %s",
		int(pr.Number),
		pr.PullRequestData.Base.Repo.Name)

	result, _, err := client.PullRequests.Merge(context.Background(),
		pr.PullRequestData.Base.User.Login,
		pr.PullRequestData.Base.Repo.Name, int(pr.Number),
		fmt.Sprintf("Merging pull request %d", pr.Number),
		&github.PullRequestOptions{})

	return result, err

}

/*closes the pull request*/
func (pr *PullRequest) ClosePullRequest(client *github.Client) {

	log.Printf("Closing pull request %d on %s",
		int(pr.Number),
		pr.PullRequestData.Base.Repo.Name)

	result, res, err := client.PullRequests.Edit(context.Background(), pr.PullRequestData.Base.User.Login,
		pr.PullRequestData.Base.Repo.Name, int(pr.Number), &github.PullRequest{
			State: github.String("closed"),
		})
	if err != nil {
		log.Println(err.Error())
	}
	
	log.Println(res.StatusCode)
	log.Println(*result.State)
}

func (pr PullRequest) CommentOnPullRequest(client *github.Client, comment *github.PullRequestComment) {

	_, res, err := client.PullRequests.CreateComment(context.Background(),
		pr.PullRequestData.Base.User.Login,
		pr.PullRequestData.Base.Repo.Name, int(pr.Number),
		comment)

	if err != nil {
		log.Println(err.Error())
	}
	log.Println(res.Status)

}

func (pr PullRequest) CreateReview(client *github.Client, event string, body string) {

	_, res, err := client.PullRequests.CreateReview(context.Background(),
		pr.PullRequestData.Base.User.Login,
		pr.PullRequestData.Base.Repo.Name,
		int(pr.Number),
		&github.PullRequestReviewRequest{
			Body:  github.String(body),
			Event: github.String(event),
		})

	if err != nil {
		log.Println(err.Error())
	}

	log.Println(res.Status)



}

func (pr PullRequest) GetFiledata(client *github.Client, sha string) []byte {

	//to get raw blob of given file's SHA hash
	blob, _, err := client.Git.GetBlobRaw(context.Background(),
		pr.PullRequestData.User.Login,
		pr.PullRequestData.Head.Repo.Name,
		sha)
	if err != nil {
		log.Println(err.Error())
	}

	return blob
}

func (pr PullRequest) GetFileFromMain(client *github.Client, path string) []byte {

	//to get certain file from main respoitory
	fc, _, st, err := client.Repositories.GetContents(
		context.Background(),
		pr.PullRequestData.Base.User.Login,
		pr.PullRequestData.Head.Repo.Name,
		path, &github.RepositoryContentGetOptions{
			Ref: pr.PullRequestData.Base.Sha,
		})

	if err != nil {
		log.Println(err.Error())
	}

	log.Println(st.StatusCode)
	data, _ := fc.GetContent()
	return []byte(data)
}

func VerifySignature(payload []byte, signature string) bool {

	key := hmac.New(sha256.New, []byte(os.Getenv("WEBHOOK_SECRET")))
	key.Write([]byte(string(payload)))
	computedSignature := "sha256=" + hex.EncodeToString(key.Sum(nil))

	return computedSignature == signature
}
