package dataops

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/xedflix/auto-approval-system/model"
)

/*ReadPoliciesFrom walks in given directory
and check if file is yaml file or not and if
it is yaml file it marshal it into Policy struct.
returns slice of Policy struct*/
func ReadPoliciesFrom(path string) []model.Policy {
	var policies []model.Policy
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".yaml" {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Println(err.Error())
			}
			temp := model.NewPolicyFrom(data)
			if err != nil {
				log.Println("Error marshaling file", path, err.Error())
				return nil
			} else {
				policies = append(policies, temp)

			}
		}
		return err
	})
	return policies
}
