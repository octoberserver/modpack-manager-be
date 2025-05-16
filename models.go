package main

import "time"

type Modpack struct {
	ID        string    `gorm:"primaryKey" json:"id" binding:"required"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	MCVersion string    `json:"mc_version"`
	Links     Links     `gorm:"embedded" json:"links"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Links struct {
	Launcher string `json:"launcher"`
	CF       string `json:"cf"`
	TRD      string `json:"trd"`
}

type LatestModpack struct {
	Server    string  `gorm:"primaryKey" json:"server" binding:"required"`
	ModpackID string  `json:"modpack_id"`
	Modpack   Modpack `gorm:"foreignKey:ModpackID" json:"modpack"`
}
