package model

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex;not null;" validate:"required"`
    Password string `validate:"required"`
}
