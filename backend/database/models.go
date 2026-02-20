package database

import "time"

// Job table
type Job struct {
	ID         uint `gorm:"primaryKey"` // GORM requires capitalized ID
	ExternalId string
	Name       string
	Number     string
}

// Material table
type Material struct {
	ID            uint `gorm:"primaryKey"`
	JobID         uint `gorm:"not null"` // foreign key to Job
	UnitOfMeasure string
	Name          string

	Job Job `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
}

// MaterialLog table
type MaterialLog struct {
	ID         uint      `gorm:"primaryKey"`
	JobID      uint      `gorm:"not null"`
	MaterialID uint      `gorm:"not null"`
	Date       time.Time `gorm:"not null"`
	Quantity   float64

	Job      Job      `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
	Material Material `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE"`
}

// Employee table
type Employee struct {
	ID         uint `gorm:"primaryKey"`
	ExternalId string
	EmployeeID string
	FirstName  string
	LastName   string
}

// CostCode table
type CostCode struct {
	ID          uint `gorm:"primaryKey"`
	ExternalId  string
	Code        string
	Description string
}

type Classification struct {
	ID         uint `gorm:"primaryKey"`
	ExternalId string
	Name       string
}
