// Package greenergrass provides a utility API to assist with integration needs
// regarding data transformation, normalization, cleanup, etc.
package greenergrass

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

const testVersion = 3

// Name contains a common list fields that could be combined as a person's full name
type Name struct {
	First, Middle, Last, Prefix, Suffix string
	full, formatted, initials           string
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
	_, err := titleFiles("", false)
	if err != nil {
		return err
	}
	return nil
}

// LoadTitleDataCSV works the same as LoadTitleData(), but expects a file path as input with should contain a .csv file extension.
func LoadTitleDataCSV(path string) error {
	if strings.Contains(path, ".csv") || strings.Contains(path, ".CSV") {
		_, err := titleFiles(path, true)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid file - must be .csv")
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

// Initials returns the acronym for the Name.  The argument del will default to an empty string, but a common input would be "." producing a result with periods (eg. J.K.M.)
func (n *Name) Initials(dots bool) string {
	if len(n.initials) > 0 {
		return n.initials
	}
	var res []rune
	tokens := strings.Split(n.formatted, " ")

	for _, token := range tokens {
		temp, _ := utf8.DecodeRuneInString(token)
		res = append(res, temp)
		if dots {
			res = append(res, '.')
		}
	}
	n.initials = strings.ToUpper(string(res))
	return n.initials
}

var titleList = make(map[string]struct{})

func titleFiles(filePath string, isCSV bool) (map[string]struct{}, error) {

	if filePath == "" {
		// filePath = "titles.csv"
		filePath = "default_titles.json"
	}

	// Checks if the desired file is a csv and will process the fields by line accordingly.
	if isCSV {
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

	// Load default data contained locally with the package.
	defaultData := gjson.GetBytes(Defaults, "titles.#.title")
	defaultData.ForEach(func(key, value gjson.Result) bool {
		titleList[value.String()] = struct{}{}
		return true
	})
	return titleList, nil
}

// FormatName creates a struct member with a formatted full name.
func (n *Name) FormatName() {
	n.formatted = fmt.Sprintf("%v %v %v", n.First, n.Middle, n.Last)
}

// FormattedName returns the formatted full name string created by FormatFullName()
func (n *Name) FormattedName() string {
	return n.formatted
}
