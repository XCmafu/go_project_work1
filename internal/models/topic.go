// internal/models/topic.go
package models

import "time"

// Topic 表示一个话题的数据结构
type Topic struct {
	ID         int       // 话题的唯一标识符
	Title      string    // 话题的标题
	Content    string    // 话题的内容
	CreateTime time.Time // 话题的创建时间
}
