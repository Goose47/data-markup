package models

import (
	"time"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `gorm:"unique;not null"`
	CreatedAt   time.Time
	Roles       []Role       `gorm:"many2many:user_roles;"`
	Permissions []Permission `gorm:"many2many:user_permissions;"`
	Batches     []Batch      `gorm:"many2many:user_batches;"`
	Assessments []Assessment `gorm:"foreignKey:UserID;references:ID"`
}

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
	//Users []User `gorm:"many2many:user_roles;"`
}

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
	//Users []User `gorm:"many2many:user_permissions;"`
}

type Batch struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Overlaps  int
	Priority  int
	CreatedAt time.Time
	IsActive  bool
	Markups   []Markup `gorm:"foreignKey:BatchID;references:ID"`
	Users     []User   `gorm:"many2many:user_batches;"`
}

//type UserBatch struct {
//	ID      uint  `gorm:"primaryKey"`
//	UserID  uint  `gorm:"index"`
//	BatchID uint  `gorm:"index"`
//	User    User  `gorm:"foreignKey:UserID;references:ID"`
//	Batch   Batch `gorm:"foreignKey:BatchID;references:ID"`
//}

type Markup struct {
	ID          uint         `gorm:"primaryKey"`
	BatchID     uint         `gorm:"index"`
	Data        string       `gorm:"type:text"`
	Batch       Batch        `gorm:"foreignKey:BatchID;references:ID"`
	Assessments []Assessment `gorm:"foreignKey:MarkupID;references:ID"`
}

type MarkupType struct {
	ID      uint   `gorm:"primaryKey"`
	BatchID uint   `gorm:"index"`
	Name    string `gorm:"not null"`
	ChildID *uint
	Fields  []MarkupTypeField `gorm:"foreignKey:MarkupTypeID;references:ID"`
	Batch   Batch             `gorm:"foreignKey:BatchID;references:ID"`
}

type MarkupTypeField struct {
	ID               uint   `gorm:"primaryKey"`
	MarkupTypeID     uint   `gorm:"index"`
	AssessmentTypeID uint   `gorm:"index"`
	Name             string `gorm:"not null"`
	GroupID          uint
	MarkupType       MarkupType     `gorm:"foreignKey:MarkupTypeID;references:ID"`
	AssessmentType   AssessmentType `gorm:"foreignKey:AssessmentTypeID;references:ID"`
}

type AssessmentType struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type Assessment struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	MarkupID  uint `gorm:"index"`
	CreatedAt time.Time
	IsPrior   bool
	Fields    []AssessmentField `gorm:"foreignKey:AssessmentID;references:ID"`
	User      User              `gorm:"foreignKey:UserID;references:ID"`
	Markup    Markup            `gorm:"foreignKey:MarkupID;references:ID"`
}

type AssessmentField struct {
	ID                uint            `gorm:"primaryKey"`
	AssessmentID      uint            `gorm:"index"`
	MarkupTypeFieldID uint            `gorm:"index"`
	Assessment        Assessment      `gorm:"foreignKey:AssessmentID;references:ID"`
	MarkupTypeField   MarkupTypeField `gorm:"foreignKey:MarkupTypeFieldID;references:ID"`
}
