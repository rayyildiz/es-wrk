package worker

import (
	"log"
	"testing"
)

func TestNewDictionary(t *testing.T) {
	dict, err := NewDictionary()

	if err != nil {
		t.Fatalf("%v", err)
	}

	log.Printf("number of words count is %d", len(dict.words))
}

func TestDict_RandomWords(t *testing.T) {

	dict, err := NewDictionary()

	if err != nil {
		t.Fatalf("%v", err)
	}

	w1, err := dict.GenerateRandomWords(5)

	if err != nil {
		t.Errorf("could not generate %d words, %v", 5, err)
	}

	if len(w1) != 5 {
		t.Errorf("could not generate %d words, got %d", 5, len(w1))
	}

}
