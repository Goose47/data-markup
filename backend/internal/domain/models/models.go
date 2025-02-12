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
//	UserID  uint  ``
//	BatchID uint  ``
//	User    User  `gorm:"foreignKey:UserID;references:ID"`
//	Batch   Batch `gorm:"foreignKey:BatchID;references:ID"`
//}

type Markup struct {
	ID          uint         `gorm:"primaryKey"`
	BatchID     uint         ``
	Data        string       `gorm:"type:text"`
	Batch       Batch        `gorm:"foreignKey:BatchID;references:ID"`
	Assessments []Assessment `gorm:"foreignKey:MarkupID;references:ID"`
}

type MarkupType struct {
	ID      uint              `gorm:"primaryKey" json:"id"`
	BatchID *uint             `gorm:"null" json:"batch_id"`
	Name    string            `gorm:"not null" json:"name"`
	ChildID *uint             `gorm:"null" json:"child_id"`
	UserID  *uint             `gorm:"null" json:"user_id"`
	Fields  []MarkupTypeField `gorm:"foreignKey:MarkupTypeID;references:ID;onDelete:CASCADE" json:"fields"`
	Batch   Batch             `gorm:"foreignKey:BatchID;references:ID" json:"-"`
}

type MarkupTypeField struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	MarkupTypeID     uint           `json:"markup_type_id"`
	AssessmentTypeID uint           `json:"assessment_type_id"`
	Name             string         `gorm:"not null" json:"name"`
	GroupID          uint           `json:"group_id"`
	MarkupType       MarkupType     `gorm:"foreignKey:MarkupTypeID;references:ID" json:"-"`
	AssessmentType   AssessmentType `gorm:"foreignKey:AssessmentTypeID;references:ID" json:"assessment_type"`
}

type AssessmentType struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

type Assessment struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint ``
	MarkupID  uint ``
	CreatedAt time.Time
	IsPrior   bool
	Fields    []AssessmentField `gorm:"foreignKey:AssessmentID;references:ID"`
	User      User              `gorm:"foreignKey:UserID;references:ID"`
	Markup    Markup            `gorm:"foreignKey:MarkupID;references:ID"`
}

type AssessmentField struct {
	ID                uint            `gorm:"primaryKey"`
	AssessmentID      uint            ``
	MarkupTypeFieldID uint            ``
	Assessment        Assessment      `gorm:"foreignKey:AssessmentID;references:ID"`
	MarkupTypeField   MarkupTypeField `gorm:"foreignKey:MarkupTypeFieldID;references:ID"`
}
