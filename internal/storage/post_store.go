// internal/storage/post_store.go
package storage

import (
	"bufio"
	"fmt"
	"github.com/XCmafu/go_project_work1/internal/models"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	posts     = make(map[int][]models.Post) // 存储回帖数据的映射
	postIDSeq = 1                           // 用于生成唯一的回帖 ID
	postLock  sync.RWMutex                  // 用于保护对 posts 的并发访问
	postFile  *os.File                      // 用于操作回帖数据文件
)

var (
	postDataFilePath = "data/posts.txt" // 回帖数据文件路径
)

func init() {
	// 初始化回帖数据文件，如果文件不存在则创建
	var err error
	postFile, err = os.OpenFile(postDataFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Failed to initialize post data file:", err)
		return
	}

	// 加载已有数据到 posts 中
	loadExistingPosts()
}

func loadExistingPosts() {
	postLock.Lock()         // 获取写锁，保护对 posts 的并发访问
	defer postLock.Unlock() // 函数执行完后释放写锁

	postFile.Seek(0, 0) // 将文件指针移到文件开头
	scanner := bufio.NewScanner(postFile)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "|")
		if len(fields) != 4 {
			continue
		}
		postID, _ := strconv.Atoi(fields[0])
		topicID, _ := strconv.Atoi(fields[1])
		createTime, _ := time.Parse(time.RFC3339, fields[3])
		post := models.Post{
			ID:         postID,
			TopicID:    topicID,
			Content:    fields[2],
			CreateTime: createTime,
		}
		if _, found := posts[topicID]; !found {
			// 如果话题不存在，就创建一个新的切片，并将 post 添加到这个切片中。
			posts[topicID] = []models.Post{post}
		} else {
			// 如果话题存在，将新的 post 追加到话题对应的帖子切片中。
			posts[topicID] = append(posts[topicID], post)
		}
		if postID >= postIDSeq {
			postIDSeq = postID + 1 // 更新 postIDSeq
		}
	}
}

// SavePost 保存一个新的回帖到 posts 中，并将数据写入文件
func SavePost(post models.Post) {
	postLock.Lock()         // 获取写锁，保护对 posts 的并发访问
	defer postLock.Unlock() // 函数执行完后释放写锁

	// 为回帖分配一个唯一的 ID，并设置创建时间为当前时间
	post.ID = postIDSeq
	post.CreateTime = time.Now()

	// 将回帖存储到 posts 映射中，使用 topicID 作为键来分组回帖
	posts[post.TopicID] = append(posts[post.TopicID], post)
	postIDSeq++

	// 将回帖数据写入文件
	writePostToFile(post)
}

// LoadPosts 从文件中加载特定话题的所有回帖数据并返回一个回帖列表
func LoadPosts(topicID int) []models.Post {
	postLock.RLock()         // 获取读W锁，保护对 posts 的并发访问
	defer postLock.RUnlock() // 函数执行完后释放读锁

	return posts[topicID]
}

// writePostToFile 将回帖数据写入文件
func writePostToFile(post models.Post) {
	postString := fmt.Sprintf("%d|%d|%s|%s\n", post.ID, post.TopicID, post.Content, post.CreateTime.Format(time.RFC3339))
	_, _ = postFile.WriteString(postString)
}
