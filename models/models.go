package models

type User struct {
	Id       uint   `json:"id"`
	UserName string `gorm:"unique; not null" json:"userName"`
	Email    string `gorm:"unique; not null" json:"email"`
	Password []byte `gorm:"not null" json:"-"`
	IsAdmin  bool   `json:"isAdmin"`
}

// Further considerations: User can have a profile photo and bio

type Post struct {
	Id       uint      `json:"id"`
	Body     string    `json:"body"`
	Title    string    `json:"title"`
	Tag      string    `json:"tag"`
	UserId   uint      `json:"userId"`
	User     User      `json:"user" gorm:"foreignKey:UserId"`
	Comments []Comment `json:"comments" gorm:"foreignKey:PostId"`
	Likes    int       `json:"likes"`
	IsEdited bool      `isEdited`
}

// A post can have a single user
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
