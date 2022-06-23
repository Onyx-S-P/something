package repo

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/xedflix/auto-approval-system/config"
)

func CloneGitRepo(dir string, cfg config.Config) (string, error) {

	var err error

	repoDir := RandomString(5)

	basePath := filepath.Join(dir, repoDir)

	if cfg.Password != "" && cfg.Username != "" {
		_, err = git.PlainClone(basePath, false, &git.CloneOptions{
			URL:      cfg.Url,
			Progress: os.Stdout,
			Auth: &http.BasicAuth{
				Username: cfg.Username,
				Password: cfg.Password,
			},
		})
	} else {
		_, err = git.PlainClone(basePath, false, &git.CloneOptions{
			URL:      cfg.Url,
			Progress: os.Stdout,
		})
	}
	return basePath, err
}

//
func RandomString(length int) string {

	bi := make([]byte, length)

	for i := range bi {
		bi[i] = stringWithCharaset[randFun.Intn(len(stringWithCharaset))]
	}

	return string(bi)
}

func String(length int) string {
	return RandomString(length)
}

var randFun *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const stringWithCharaset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
