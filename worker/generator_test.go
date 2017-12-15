package worker

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

type TestConfidence struct {
	Label      string  `json:"label"`
	Confidence float32 `json:"confidence"`
}

func TestElastic_GetRandomElem(t *testing.T) {

	api, err := NewGenerator(reflect.TypeOf(TestConfidence{}))
	if err != nil {
		t.Fatalf("could not create api, %v", err)
	}
	p := api.getRandomElem()

	if p == nil {
		t.Fatal("Post is empty")
	}

}

func TestGetRandomPosts(t *testing.T) {
	testCases := []struct {
		numberOfPosts int
	}{
		{0},
		{1},
		{100},
		{1000},
	}
	api, err := NewGenerator(reflect.TypeOf(TestConfidence{}))
	if err != nil {
		t.Fatalf("could not create api, %v", err)
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("generate %d elements ", tc.numberOfPosts), func(t *testing.T) {
			posts := api.GetRandomElements(tc.numberOfPosts)
			if len(posts) != tc.numberOfPosts {
				t.Errorf("could not generate %d posts, got %d", tc.numberOfPosts, len(posts))
			}
		})
	}
}

func TestElastic_RandomWords(t *testing.T) {
	e, err := NewGenerator(reflect.TypeOf(TestConfidence{}))

	if err != nil {
		t.Fatalf("could not create elastic api, %v", err)
	}

	str := e.randomWords(15)
	if len(str) < 30 {
		t.Errorf("could not create words")
	}

	log.Println(str)
}

func BenchmarkGetRandomPosts(b *testing.B) {
	api, err := NewGenerator(reflect.TypeOf(TestConfidence{}))
	if err != nil {
		b.Fatalf("could not create api, %v", err)
	}

	for n := 0; n < b.N; n++ {
		api.GetRandomElements(100)
	}
}

func BenchmarkElastic_RandomWords(b *testing.B) {
	api, err := NewGenerator(reflect.TypeOf(TestConfidence{}))
	if err != nil {
		b.Fatalf("could not create api, %v", err)
	}

	for n := 0; n < b.N; n++ {
		str := api.randomWords(10)
		if len(str) < 11 {
			b.Errorf("could not create words")
		}
	}

}
