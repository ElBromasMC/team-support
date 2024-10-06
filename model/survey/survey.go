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

type Respondents struct {
	Id          int
	Survey      Survey
	Name        string
	Email       string
	PhoneNumber string
	Rating      int
}

type Responses struct {
	Id           int
	Respondents  Respondents
	Question     Question
	ResponseText string
	CreatedAt    time.Time
}

type Landing struct {
	Id          int
	Title       string
	Content     string
	CreatedAt   time.Time
	IsPublished bool
	Survey      Survey
}
