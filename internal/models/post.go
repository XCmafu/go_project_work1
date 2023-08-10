// internal/models/post.go
package models

import "time"

// Post 表示一个回帖的数据结构
type Post struct {
	ID         int       // 回帖的唯一标识符
	TopicID    int       // 回帖所属话题的 ID
	Content    string    // 回帖的内容
	CreateTime time.Time // 回帖的创建时间
}
