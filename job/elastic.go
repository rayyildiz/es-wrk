package job

import (
	"fmt"
	"math/rand"
	"time"
)

type Confidence struct {
	Label      string  `json:"label"`
	Confidence float32 `json:"confidence"`
}

type MetaData struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Feature struct {
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Feature    string  `json:"feature"`
	Text       string  `json:"text"`
	Source     int     `json:"source"`
	Confidence float32 `json:"confidence"`
}

type Sentiment struct {
	Start    int     `json:"start"`
	End      int     `json:"end"`
	Text     string  `json:"text"`
	Positive bool    `json:"positive"`
	Scale    float32 `json:"scale"`
}

type FeatureSentiment struct {
	FeatureStart   int     `json:"featureStart"`
	FeatureEnd     int     `json:"featureEnd"`
	Sentence       string  `json:"sentence"`
	Feature        string  `json:"feature"`
	FeatureText    string  `json:"featureText"`
	SentimentText  string  `json:"sentimentText"`
	SentimentStart int     `json:"sentimentStart"`
	SentimentEnd   int     `json:"sentimentEnd"`
	Positive       bool    `json:"positive"`
	SentimentScale float32 `json:"sentimentScale"`
}

type Post struct {
	ID                string             `json:"id"`
	Text              string             `json:"text"`
	CleanText         string             `json:"cleanText"`
	CompanyID         string             `json:"companyId"`
	ProjectID         string             `json:"projectId"`
	Language          string             `json:"language"`
	ProcessedAt       string             `json:"processedAt"`
	PostDate          string             `json:"postDate"`
	CurrentStatus     int                `json:"currentStatus"`
	TargetStatus      int                `json:"targetStatus"`
	MetaData          []MetaData         `json:"metaData"`
	Profiles          []Confidence       `json:"profiles"`
	Classification    *Confidence        `json:"classification"`
	Features          []Feature          `json:"features"`
	Sentiments        []Sentiment        `json:"sentiments"`
	FeatureSentiments []FeatureSentiment `json:"featureSentiments"`
}

const constLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "
const constNumbers = "0123456789"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = constLetters[rand.Intn(len(constLetters))]
	}
	return string(b)
}

func RandNumbers(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = constNumbers[rand.Intn(len(constNumbers))]
	}
	return string(b)
}

func GetRandomPost(baseID string) *Post {
	p := Post{}
	p.ID = fmt.Sprintf("%s%s", baseID, RandNumbers(16))
	p.CleanText = RandString(72)
	p.CompanyID = RandNumbers(4)
	p.CurrentStatus = 0
	p.TargetStatus = 100
	p.Language = "en"
	p.PostDate = fmt.Sprintf("%v", time.Now())
	p.ProcessedAt = fmt.Sprintf("%v", time.Now())
	p.ProjectID = RandNumbers(3)
	p.Text = RandString(84)

	for i := 0; i < rand.Intn(10); i++ {
		f := Feature{
			Confidence: rand.Float32(),
			End:        5 + rand.Intn(20),
			Start:      rand.Intn(5),
			Source:     rand.Intn(10),
			Text:       RandString(12),
		}
		p.Features = append(p.Features, f)
	}

	for i := 0; i < rand.Intn(10); i++ {
		s := Sentiment{
			Scale:    rand.Float32(),
			End:      5 + rand.Intn(20),
			Start:    rand.Intn(5),
			Text:     RandString(12),
			Positive: (rand.Intn(15)%2 == 0),
		}
		p.Sentiments = append(p.Sentiments, s)
	}

	for i := 0; i < rand.Intn(20); i++ {
		m := MetaData{
			Type:  RandString(12),
			Value: RandString(12),
		}
		p.MetaData = append(p.MetaData, m)
	}

	for i := 0; i < rand.Intn(10); i++ {
		fs := FeatureSentiment{
			SentimentScale: rand.Float32(),
			SentimentEnd:   5 + rand.Intn(20),
			SentimentStart: rand.Intn(5),
			SentimentText:  RandString(12),
			Positive:       (rand.Intn(15)%2 == 0),
			FeatureEnd:     5 + rand.Intn(20),
			FeatureStart:   rand.Intn(5),
			Feature:        RandString(12),
		}
		p.FeatureSentiments = append(p.FeatureSentiments, fs)
	}

	for i := 0; i < rand.Intn(10); i++ {
		c := Confidence{
			Label:      RandString(10),
			Confidence: rand.Float32(),
		}
		p.Profiles = append(p.Profiles, c)
	}

	p.Classification = &Confidence{
		Label:      RandString(10),
		Confidence: rand.Float32(),
	}

	return &p
}

func GetRandomPosts(baseID string, numberOfPosts int) []Post {
	var posts []Post

	for i := 0; i < numberOfPosts; i++ {
		p := GetRandomPost(baseID)
		if p != nil {
			posts = append(posts, *GetRandomPost(baseID))
		}
	}

	return posts
}
