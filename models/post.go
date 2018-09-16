package models

import (
	"bufio"
	"strings"

	"github.com/jinzhu/gorm"
)

// Post is used to read and write new posts.
type Post struct {
	gorm.Model
	Title string
	Body  []byte `gorm:"type:varbinary(8192)"`
}

// AllPosts returns all posts stored in the DB.
func AllPosts(db *gorm.DB) ([]*Post, error) {
	var posts []*Post
	err := db.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPost returns a post for the given ID.
func GetPost(db *gorm.DB, id uint64) (*Post, error) {
	post := new(Post)
	err := db.First(&post, id).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

// GetPost2 returns a post for the given ID.
func GetPost2(db *gorm.DB, id uint64) (string, error) {
	post := new(Post)
	err := db.First(&post, id).Error
	html := getHTML(post)
	if err != nil {
		return "", err
	}

	return html, nil
}

// CreatePost inserts a new post in the DB.
func CreatePost(db *gorm.DB, body []byte) (uint, error) {
	db.AutoMigrate(&Post{})
	post := Post{Body: body}
	post.Title = post.getTitle()
	err := db.Create(&post).Error
	if err != nil {
		return 0, err
	}

	return post.ID, nil
}

func (post *Post) getTitle() string {
	body := string(post.Body)

	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#+TITLE:") {
			return strings.SplitAfter(line, "#+TITLE: ")[1]
		}
	}

	return ""
}
