package migrate

import (
	"context"
	"database/sql"
	"djson/pkg/sqlparser"
	"fmt"
	"os"
)

type Migrate struct {
	path           string
	db             *sql.DB
	migrationFiles []DirEntryWithPrefix
	txn            *sql.Tx
}

// this NewMigrate function takes datatbase and directory path and uses it to generate a struct of type Migrate and returns it
func NewMigrate(db *sql.DB, dirPath string) Migrate {
	return Migrate{
		path: dirPath,
		db:   db,
	}
}

// method to run migration

func (m *Migrate) RunMigrations() error {
	// will read dir specified by m.path
	rawEntries, err := os.ReadDir(m.path)
	if err != nil {
		fmt.Println("ERR 1")
	}
	usableEntries := m.filterSqlFileWithNumberPrefix(m.getFilesFromDirEntries(rawEntries))

	err = m.checkForSamePrefix(usableEntries)
	if err != nil {
		fmt.Printf("ERR2")
		return err
	}
	version, err1 := m.getVersion()
	if err1 != nil {
		fmt.Println("ERR3")

		return err1
	}

	//this means we are already on latestDB state
	if version == len(usableEntries) {
		fmt.Println("ERR4")
		return nil
	}

	//this means db was never created
	if version == -1 {
		m.migrationFiles = usableEntries
	} else {
		m.migrationFiles = usableEntries[version:]
	}

	m.txn, err = m.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer m.txn.Rollback()
	err = m.parseFilesAndMigrateDb()
	if err != nil {
		return err
	}
	// clear file
	err = os.Truncate(m.path+"/migrate.log", 0)
	if err != nil {
		return err
	}

	// writing latest db version to file
	latest := []byte(fmt.Sprintf("%d", len(usableEntries)))
	outFile, err2 := os.OpenFile(m.path+"/migrate.log", os.O_RDWR, 0777)
	if err2 != nil {
		return err
	}
	_, err2 = outFile.Write(latest)
	if err2 != nil {
		return err
	}

	err = outFile.Close()
	if err != nil {
		return err
	}

	err = m.txn.Commit()
	if err != nil {
		return err
	}

	return nil

}

func (m *Migrate) parseFilesAndMigrateDb() error {
	for _, file := range m.migrationFiles {
		filePath := m.path + "/" + file.Dir.Name()
		fmt.Printf("Reading File %s\n", file.Dir.Name())
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		content := string(bytes)
		commands := sqlparser.ParseSqlFile(content)
		for _, command := range commands {
			_, err = m.txn.Exec(command)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
