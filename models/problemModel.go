// models/problem.go
package models

type Problem struct {
	Pid         int32      `json:"pid,omitempty" bson:"pid,omitempty"`
	Title       string     `json:"title" bson:"title"`
	Description string     `json:"description" bson:"description"`
	Constraints string     `json:"constraints" bson:"constraints"`
	TestCases   []TestCase `json:"test_cases" bson:"test_cases"`
	AuthorID    string     `json:"author_id" bson:"author_id"`
	Visibility  bool       `json:"visibility" bson:"visibility"`
}

type TestCase struct {
	Input  string `json:"input" bson:"input"`
	Output string `json:"output" bson:"output"`
}
