package worker

import (
	"fmt"
	"math/rand"
	"strings"
)

type dict struct {
	words []string
}

// NewDictionary loads the dictionary file give the words.
func NewDictionary() (*dict, error) {
	data, err := Asset("../data/dict.txt")
	if err != nil {
		return nil, fmt.Errorf("could not read the file, %v", err)
	}

	d := dict{}
	str := string(data[:])
	d.words = strings.Split(str, "\n")

	return &d, nil
}

// GenerateRandomWords generate random words using dictionary file.
func (d *dict) GenerateRandomWords(no int) ([]string, error) {
	if len(d.words) == 0 {
		return nil, fmt.Errorf("dictionary is empty")
	}

	var xs []string

	for i := 0; i < no; i++ {
		ind := rand.Intn(len(d.words))

		xs = append(xs, d.words[ind])
	}

	return xs, nil
}
