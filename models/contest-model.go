// models/contest.go
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contest struct {
	ContestID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	StartTime   int64              `json:"start_time" bson:"start_time"`
	EndTime     int64              `json:"end_time" bson:"end_time"`
	HostID      string             `json:"host_id" bson:"host_id"`
	Problems    []int32            `json:"problems" bson:"problems"` // Array of problem PIDs
}

type Participant struct {
	ParticipantId primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ContestID     primitive.ObjectID `json:"contest_id,omitempty" bson:"contest_id,omitempty"`
	UserID        string             `json:"user_id" bson:"user_id"`
	Score         int32              `json:"score" bson:"score"`
	SubmissionID  string             `json:"submission_id,omitempty" bson:"submission_id,omitempty"`
}

type Leaderboard struct {
	ContestID    primitive.ObjectID `json:"contest_id,omitempty" bson:"contest_id,omitempty"`
	Participants []Participant      `json:"participants" bson:"participants"`
}
