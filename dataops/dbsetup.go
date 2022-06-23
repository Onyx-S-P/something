package dataops

import (
	"log"
	"path/filepath"

	"github.com/xedflix/auto-approval-system/config"
	"github.com/xedflix/auto-approval-system/repo"
)

func SetupBase(configDir, tmpDir string) {

	log.Println("Setup base repo db")
	cfg, err := config.ReadConfig(configDir)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(cfg.Url)
	path, err := repo.CloneGitRepo(tmpDir, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(path)

	//Creates the connection with database and creates policy table
	conn, err := NewSqliteCilent(filepath.Join(tmpDir, "policy.db"))
	conn.CreateTables()

	if err != nil {
		log.Println(err.Error())
	}

	//unmarshal all polices from given root directory
	polices := ReadPoliciesFrom(path)

	//insert each policy into database
	for _, p := range polices {
		res, err := conn.InsertPolicy(p)
		if err != nil {
			log.Println(err.Error())
		} else {
			id, _ := res.LastInsertId()
			log.Println("Last insterted id ", id)
		}
		conn.Readpolicy(p)
		/*res1, err1 := conn.Readpolicy(p)
		if err1 != nil {
			log.Println(err.Error())
		} else {
			id1, _ := res1.LastInsertId()
			log.Println("Last inserted id ", id1)
		}*/
	}
}
