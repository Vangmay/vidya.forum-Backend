package models

type User struct {
	Id       uint   `json:"id"`
	UserName string `gorm:"unique; not null" json:"userName"`
	Email    string `gorm:"unique; not null" json:"email"`
	Password []byte `gorm:"not null" json:"-"`
	IsAdmin  bool   `json:"isAdmin"`
}

type Post struct {
	Id       uint      `json:"id"`
	Content  string    `json:"content"`
	Title    string    `json:"title"`
	UserId   uint      `json:"userId"`
	User     User      `json:"user" gorm:"foreignKey:UserId"`
	Comments []Comment `json:"comments" gorm:"foreignKey:PostId"`
	Likes    int       `json:"likes"`
	Tag      string    `json:"tag"`
}

// A Post has many comments

type Comment struct {
	Id      uint   `json:"id"`
	Content string `json:"content"`
	PostId  uint   `json:"postId"`
	Post    Post   `json:"post" gorm:"foreignKey:PostId"`
	UserId  uint   `json:"userId"`
	User    User   `json:"user" gorm:"foreignKey:UserId"`
}

// Comment Belongs to a post
// Comment also belongs to a user
