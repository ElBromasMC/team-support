package survey

import "time"

type Survey struct {
	Id        int
	Title     string
	CreatedAt time.Time
}

type QuestionType string

const (
	TextType QuestionType = "TEXT"
)

type Question struct {
	Id           int
	Survey       Survey
	QuestionText string
	QuestionType QuestionType
}

type QuestionResponse struct {
	Question     Question
	ResponseText string
}

type SurveyResponse struct {
	Id                int
	Survey            Survey
	Name              string
	Email             string
	PhoneNumber       string
	Rating            int
	QuestionResponses []QuestionResponse
	Responses         []string // For csv generation
	CreatedAt         time.Time
}

type Landing struct {
	Id          int
	Title       string
	Content     string
	CreatedAt   time.Time
	IsPublished bool
	Survey      Survey
}
