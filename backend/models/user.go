package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-"` // Never expose password in JSON responses
	Name      string    `json:"name"`
	Role      string    `json:"role"` // admin or user
	Bio       string    `json:"bio"`  // VULNERABILITY #3: XSS if not sanitized
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
