package migrate

import (
	"errors"
	"fmt"
	"io/fs"
	"strconv"
)

type DirEntryWithPrefix struct {
	Prefix int
	Dir    fs.DirEntry
}

// filter sql queries with number prefix
func (m *Migrate) filterSqlFileWithNumberPrefix(entries []fs.DirEntry) []DirEntryWithPrefix {
	var result []DirEntryWithPrefix
	for _, v := range entries {
		prefix, err := m.getPrefixFromName(v.Name())
		if err == nil {
			result = append(result, DirEntryWithPrefix{
				Prefix: prefix, Dir: v,
			})
		}
	}
	return result
}

func (m *Migrate) getPrefixFromName(name string) (int, error) {
	if len(name) == 0 {
		return 0, errors.New("Invalid file name")
	}
	prefix := 0
	for i := 1; i < len(name); i++ {
		num, err := strconv.Atoi(name[:i])
		if err != nil {
			if i == 1 {
				return 0, errors.New("Invalid file name")
			} else {
				break
			}
		} else {
			prefix = num
		}
	}
	return prefix, nil
}

func (m *Migrate) getFilesFromDirEntries(entries []fs.DirEntry) []fs.DirEntry {
	var result []fs.DirEntry
	for _, v := range entries {
		if !v.IsDir() {
			result = append(result, v)
		}
	}
	return result
}

func (m *Migrate) checkForSamePrefix(entries []DirEntryWithPrefix) error {
	if len(entries) == 0 || len(entries) == 1 {
		return nil
	}
	for i := 0; i < len(entries)-1; i++ {
		if entries[i].Prefix == entries[i+1].Prefix {
			return fmt.Errorf("same prefix found for %d", entries[i].Prefix)
		}
	}
	return nil
}
