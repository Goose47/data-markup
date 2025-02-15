package models

import (
	"slices"
	"strconv"
	"strings"
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
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"not null"`
	Overlaps    int          `json:"overlaps"`
	Priority    int          `json:"priority"`
	CreatedAt   time.Time    `json:"created_at"`
	IsActive    bool         `json:"is_active"`
	TypeID      uint         `json:"type_id"`
	Markups     []Markup     `json:"-" gorm:"foreignKey:BatchID;references:ID"`
	MarkupTypes []MarkupType `json:"-" gorm:"foreignKey:BatchID;references:ID"`
	Users       []User       `json:"-" gorm:"many2many:user_batches;"`
}

//type UserBatch struct {
//	ID      uint  `gorm:"primaryKey"`
//	UserID  uint  ``
//	BatchID uint  ``
//	User    User  `gorm:"foreignKey:UserID;references:ID"`
//	Batch   Batch `gorm:"foreignKey:BatchID;references:ID"`
//}

type Markup struct {
	ID                    uint         `json:"id" gorm:"primaryKey"`
	BatchID               uint         `json:"batch_id"`
	StatusID              uint         `json:"status_id"`
	Data                  string       `json:"data" gorm:"type:text"`
	CorrectAssessmentHash *string      `json:"correct_assessment_hash" gorm:"null"`
	Batch                 Batch        `json:"-" gorm:"foreignKey:BatchID;references:ID"`
	Assessments           []Assessment `json:"assessments" gorm:"foreignKey:MarkupID;references:ID"`
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
	Name             *string        `gorm:"null" json:"name"`
	Label            string         `gorm:"not null" json:"label"`
	GroupID          uint           `json:"group_id"`
	MarkupType       MarkupType     `gorm:"foreignKey:MarkupTypeID;references:ID" json:"-"`
	AssessmentType   AssessmentType `gorm:"foreignKey:AssessmentTypeID;references:ID" json:"assessment_type"`
}

type AssessmentType struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

type Assessment struct {
	ID        uint              `json:"id" gorm:"primaryKey"`
	UserID    uint              `json:"user_id"`
	MarkupID  uint              `json:"markup_id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt *time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	IsPrior   bool              `json:"is_prior"`
	Hash      *string           `json:"hash"`
	Fields    []AssessmentField `json:"fields" gorm:"foreignKey:AssessmentID;references:ID"`
	User      User              `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Markup    Markup            `json:"-" gorm:"foreignKey:MarkupID;references:ID"`
}

func (a Assessment) CalculateHash() string {
	ids := make([]uint, len(a.Fields))
	for i, field := range a.Fields {
		ids[i] = field.MarkupTypeFieldID
	}
	slices.Sort(ids)

	idsString := make([]string, len(ids))
	for i, id := range ids {
		idsString[i] = strconv.Itoa(int(id))
	}
	return strings.Join(idsString, ",")
}

type AssessmentField struct {
	ID                uint            `json:"id" gorm:"primaryKey"`
	AssessmentID      uint            `json:"assessment_id"`
	MarkupTypeFieldID uint            `json:"markup_type_field_id"`
	Text              *string         `json:"text"`
	Assessment        Assessment      `json:"-" gorm:"foreignKey:AssessmentID;references:ID"`
	MarkupTypeField   MarkupTypeField `json:"-" gorm:"foreignKey:MarkupTypeFieldID;references:ID"`
}
