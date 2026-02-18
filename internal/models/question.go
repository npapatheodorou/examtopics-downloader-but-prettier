package models

type CommentData struct {
	User   string
	Answer string
	Text   string
}

type QuestionData struct {
	Title        string
	Header       string
	Content      string
	ExhibitURLs  []string
	Questions    []string
	Answer       string
	Timestamp    string
	QuestionLink string
	Comments     []CommentData
}
