package models

import (
	"time"
)

// Define struct for entities
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Hashtag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UserIds     []int     `json:"user_ids"`
	HashtagIds  []int     `json:"hashtag_ids"`
}

type DenormalizedProject struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Users       []User    `json:"users"`
	Hashtags    []Hashtag `json:"hashtags"`
}

type Payload struct {
	Doc    DenormalizedProject `json:"doc"`
	Method string              `json:"method"`
}
