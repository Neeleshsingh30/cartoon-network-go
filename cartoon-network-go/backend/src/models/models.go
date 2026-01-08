package models

import "time"

/* USERS */
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time

	Likes []Like        `gorm:"foreignKey:UserID"`
	Views []CartoonView `gorm:"foreignKey:UserID"`
}

/* CARTOONS */
type Cartoon struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
	Genre       string
	AgeGroup    string
	Universe    string
	ShowTime    string
	ImdbRating  float32
	CreatedAt   time.Time

	Images     []CartoonImage `gorm:"foreignKey:CartoonID"`
	Characters []Character    `gorm:"foreignKey:CartoonID"`
	Likes      []Like         `gorm:"foreignKey:CartoonID"`
	Views      []CartoonView  `gorm:"foreignKey:CartoonID"`
}

/* CARTOON IMAGES */
type CartoonImage struct {
	ID        uint `gorm:"primaryKey"`
	CartoonID uint
	ImageURL  string
	ImageType string

	Cartoon Cartoon `gorm:"foreignKey:CartoonID"`
}

/* CHARACTERS */
type Character struct {
	ID          uint `gorm:"primaryKey"`
	CartoonID   uint
	Name        string
	Role        string
	Power       string
	Description string

	Cartoon Cartoon `gorm:"foreignKey:CartoonID"`
}

/* LIKES */
type Like struct {
	ID        uint `gorm:"primaryKey"`
	CartoonID uint
	UserID    uint

	Cartoon Cartoon `gorm:"foreignKey:CartoonID"`
	User    User    `gorm:"foreignKey:UserID"`
}

/* CARTOON VIEWS */
type CartoonView struct {
	ID        uint `gorm:"primaryKey"`
	CartoonID uint
	UserID    uint
	ViewedAt  time.Time

	Cartoon Cartoon `gorm:"foreignKey:CartoonID"`
	User    User    `gorm:"foreignKey:UserID"`
}

/* ADMINS */
type Admin struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string

	Logs []AdminActivityLog `gorm:"foreignKey:AdminID"`
}

/* ADMIN ACTIVITY LOGS */
type AdminActivityLog struct {
	ID        uint `gorm:"primaryKey"`
	AdminID   uint
	Action    string
	CreatedAt time.Time

	Admin Admin `gorm:"foreignKey:AdminID"`
}
