package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username           string    `json: "username"`
	Email              string    `json: "email" `
	Encrypted_password []byte    `json: "password"`
	Profile_picture    string    `json:"profilePic"`
	Bio                string    `json: "bio"`
	IsAdmin            bool      `json:"isAdmin"`
	Posts              []Post    `json: "posts" gorm:"foreignKey:PostID" `
	Notes              []Note    `json: "notes" gorm:"foreignKey:NoteID"`
	Comments           []Comment `json: "comments" gorm:"foreignKey:CommentID"`
	// An array or notes where ever element has note_id, user_id and post_id
}
type Post struct {
	gorm.Model
	PostID   uint      `json: "postId;" gorm:"primaryKey"`
	UserID   uint      `json: "userId"`
	Image    string    `json: "image"`
	Content  string    `json: "content"`
	Title    string    `json: "title"`
	Comments []Comment `json: "comments" gorm:"foreignKey:CommentID"`
	Notes    []Note    `json: "notes" gorm:"foreignKey:NoteID"`
	// Nested inside ever post will be an array of all comments which I can then map over
}

type Comment struct {
	gorm.Model
	CommentID uint   `json: "commentId", gorm:"primaryKey"`
	PostID    uint   `json: "postId"`
	UserID    uint   `json: "userId" `
	Content   string `json: "content"`
}

type Note struct {
	gorm.Model
	NoteID uint `json: "noteId" gorm:"primaryKey"`
	UserID uint `json: "userId"`
	PostID uint `json: "postId"`
}

type Tag struct {
	gorm.Model
	TagID uint   `json:"name"gorm:"primaryKey"`
	Name  string `json: "name"`
}
