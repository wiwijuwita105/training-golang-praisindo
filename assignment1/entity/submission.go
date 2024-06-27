package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Submission struct representing the submission table
type Submission struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null"`
	User         *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Answer       []byte    `gorm:"type:json"`
	RiskScore    int       `gorm:"not null"`
	RiskCategory string    `gorm:"size:255"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// JSON custom type to handle JSON data
// type JSON map[string]interface{}
type JSON []byte

type Answer struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

// Scan implements the sql.Scanner interface for JSON
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &j)
}

// Value implements the driver.Valuer interface for JSON
func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}
