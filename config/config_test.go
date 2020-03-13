package config

import (
	"testing"
)

func TestGetDbConfig(t *testing.T) {
	config := GetDbConfig()

	if len(config.DbName) == 0 && len(config.UsersDB) == 0 && len(config.BooksDB) == 0 {
		t.Fatal("FAILED !")
	}
	t.Logf("Testint DbName config Passed : %v %v %v", config.DbName, config.UsersDB, config.BooksDB)
}
