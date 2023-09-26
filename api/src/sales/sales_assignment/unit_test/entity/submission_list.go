package entity

type Submission struct {
	Finish_date string
	Task        int
	Validation  string
}

type SubmissionDetail struct {
	Id         int8
	Task       int8
	Validation string
}
