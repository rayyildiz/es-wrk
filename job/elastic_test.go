package job

import (
	"fmt"
	"testing"
)

func TestGetRandomPost(t *testing.T) {
	testCases := []struct {
		baseID   string
		idLenght int
	}{
		{"1000", 20},
		{"1", 17},
		{"", 16},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("for %s expected length %d", tc.baseID, tc.idLenght), func(t *testing.T) {
			p := GetRandomPost(tc.baseID)

			if p == nil {
				t.Fatal("Post is empty")
			}

			if len(p.ID) != tc.idLenght {
				t.Errorf("expected id lenght is %d but got %d", tc.idLenght, len(p.ID))
			}
		})
	}
}

func TestGetRandomPosts(t *testing.T) {
	testCases := []struct {
		baseID        string
		numberOfPosts int
	}{
		{"1000", 0},
		{"1", 1},
		{"200", 100},
		{"10000", 1000},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("generate %d posts starts with %s", tc.numberOfPosts, tc.baseID), func(t *testing.T) {
			posts := GetRandomPosts(tc.baseID, tc.numberOfPosts)
			if len(posts) != tc.numberOfPosts {
				t.Errorf("could not generate %d posts, got %d", tc.numberOfPosts, len(posts))
			}
		})
	}
}

func BenchmarkGetRandomPosts(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetRandomPosts("100", 100)
	}
}
