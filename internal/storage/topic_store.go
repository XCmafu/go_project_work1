// internal/storage/topic_store.go
package storage

import (
	"bufio"
	"fmt"
	"gocode/project_work1/internal/models"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	topics     = make(map[int]models.Topic) // 存储话题数据的映射
	topicIDSeq = 1                          // 用于生成唯一的话题 ID
	topicLock  sync.RWMutex                 // 用于保护对 topics 的并发访问
	topicFile  *os.File                     // 用于操作话题数据文件
)

var (
	topicDataFilePath = "data/topics.txt" // 话题数据文件路径
)

// init 初始化话题数据文件并加载已有数据
func init() {
	// 初始化话题数据文件，如果文件不存在则创建
	var err error
	topicFile, err = os.OpenFile(topicDataFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Failed to initialize topic data file:", err)
		return
	}

	// 加载已有数据到 topics 中
	loadExistingTopics()
}

// loadExistingTopics 从文件中加载已有话题数据到 topics 映射中
func loadExistingTopics() {
	topicLock.Lock()         // 获取写锁，保护对 topics 的并发访问
	defer topicLock.Unlock() // 函数执行完后释放写锁

	topicFile.Seek(0, 0) // 将文件指针移到文件开头
	scanner := bufio.NewScanner(topicFile)
	// 逐行扫描从 topicFile 文件中读取的内容。
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "|") // 将扫描器（scanner）当前行的文本内容按照竖线符号 | 进行分割，并返回一个字符串切片（数组）
		if len(fields) != 4 {
			continue
		}
		topicID, _ := strconv.Atoi(fields[0])
		createTime, _ := time.Parse(time.RFC3339, fields[3])
		topic := models.Topic{
			ID:         topicID,
			Title:      fields[1],
			Content:    fields[2],
			CreateTime: createTime,
		}
		topics[topic.ID] = topic // 将数据加载到 topics
		if topicID >= topicIDSeq {
			topicIDSeq = topicID + 1 // 更新 topicIDSeq
		}
	}
}

// SaveTopic 保存一个新的话题到 topics 中，并将数据写入文件
func SaveTopic(topic models.Topic) {
	topicLock.Lock()         // 获取写锁，保护对 topics 的并发访问
	defer topicLock.Unlock() // 函数执行完后释放写锁

	// 为话题分配一个唯一的 ID，并设置创建时间为当前时间
	topic.ID = topicIDSeq
	topic.CreateTime = time.Now()

	// 将话题存储到 topics 映射中
	topics[topic.ID] = topic
	topicIDSeq++

	// 将话题数据写入文件
	writeTopicToFile(topic)
}

// LoadTopics 从文件中加载所有话题数据并返回一个话题列表
func LoadTopics() []models.Topic {
	topicLock.RLock()         // 获取读锁，保护对 topics 的并发访问
	defer topicLock.RUnlock() // 函数执行完后释放读锁

	var topicList []models.Topic
	topicFile.Seek(0, 0) // 将文件指针移到文件开头
	scanner := bufio.NewScanner(topicFile)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "|")
		if len(fields) != 4 {
			continue
		}
		topicID, _ := strconv.Atoi(fields[0])
		createTime, _ := time.Parse(time.RFC3339, fields[3])
		topicList = append(topicList, models.Topic{
			ID:         topicID,
			Title:      fields[1],
			Content:    fields[2],
			CreateTime: createTime,
		})
	}

	return topicList
}

// LoadTopic 从文件中加载特定话题的详细信息并返回一个话题，如果找不到则返回 nil
func LoadTopic(topicID int) *models.Topic {
	topicLock.RLock()         // 获取读锁，保护对 topics 的并发访问
	defer topicLock.RUnlock() // 函数执行完后释放读锁

	topic := topics[topicID]
	if topic.ID == 0 {
		return nil // 未找到特定话题
	}
	return &topic
}

// writeTopicToFile 将话题数据写入文件
func writeTopicToFile(topic models.Topic) {
	topicString := fmt.Sprintf("%d|%s|%s|%s\n", topic.ID, topic.Title, topic.Content, topic.CreateTime.Format(time.RFC3339))
	_, _ = topicFile.WriteString(topicString)
}
