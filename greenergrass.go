// Package greenergrass provides a utility API to assist with integration needs
// regarding data transformation, normalization, cleanup, etc.
package greenergrass

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const testVersion = 2

// Name contains a common list fields that could be combined as a person's full name
type Name struct {
	full, First, Middle, Last, Prefix, Suffix string
}

// New returns a pointer to a Name and initializes full with full
func New(full string) *Name {
	return &Name{full: full}
}

// LoadTitleData creates a map consisting of title prefixes and suffixes that are common.
// This can be called optionally by the consumer if they are expecting their input data to
// include prefixes and/or suffixes
func LoadTitleData() error {
	// TODO***change param to use a env extension***
	// TODO***load default title data and add boolean argument
	_, err := titleFiles("")
	if err != nil {
		return err
	}
	return nil
}

// SeparateName uses receiver n, and parses full according to common logic and parses the full name
// with the fields separated.  If full is empty, then Name will reflect the zero values appropriately.
// If full cannot be split on sep, then Name.First will be set as the entire value of full.
func (n *Name) SeparateName(sep string) {
	if n.full == "" {
		return
	}
	if sep == "" {
		sep = " "
	}

	commaIndex := strings.IndexAny(n.full, ",")
	if commaIndex != -1 {
		n.Last = string(n.full[:commaIndex])
		n.full = string(n.full[commaIndex+1:])
		n.full = strings.TrimLeft(n.full, " ")
	}

	// parts is a slice of the full input string, or the string following the first comma if provided
	parts := strings.Split(n.full, sep)

	// check titleList to see if the first word of full is a listed prefix
	if _, ok := titleList[parts[0]]; ok {
		n.Prefix = parts[0]
		parts = parts[1:]
	}

	// check titleList to see if the last word of full is a listed suffix or title
	if _, ok := titleList[parts[len(parts)-1]]; ok {
		n.Suffix = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	}

	if len(parts) == 1 {
		n.First = parts[0]
	} else if len(parts) >= 2 && n.Last != "" {
		n.First = string(parts[0])
		n.Middle = strings.Join(parts[1:len(parts)], " ")
	} else {
		n.First = string(parts[0])
		n.Middle = strings.Join(parts[1:len(parts)-1], " ")
		n.Last = string(parts[len(parts)-1])
	}
}

var titleList = make(map[string]struct{})

func titleFiles(filePath string) (map[string]struct{}, error) {

	if filePath == "" {
		filePath = "titles.csv"
	}

	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening csv")
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "error reading csv")
	}

	for _, each := range records {
		titleList[each[0]] = struct{}{}
	}
	return titleList, nil
}
