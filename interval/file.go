package interval

import (
	"encoding/json"
	"log"
	"os"
)

// StorageTopic 存储在json题库文件中的试题信息
type StorageTopic struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Answers string `json:"answers"`
	Answer  string `json:"answer"`
}

const filePath = "./data/banks.json"

// ReadAnswerFile 读取保存在本地的题库文件
func ReadAnswerFile() map[string]*StorageTopic {
	readFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("读取题库文件错误，message: %s", err)
		return nil
	}
	banks := make(map[string]*StorageTopic)
	if err = json.Unmarshal(readFile, &banks); err != nil {
		log.Printf("格式化题库文件错误，message: %s", err)
		return nil
	}
	return banks
}

// SaveAnswersToFile 将题库数据保存到文件
func SaveAnswersToFile(banks map[string]*StorageTopic) {
	// 序列化为阅读友好的json格式
	jsonValue, _ := json.MarshalIndent(banks, "", "  ")
	if err := os.WriteFile(filePath, jsonValue, 0644); err != nil {
		log.Printf("写入题库到文件错误，message: %s", err)
	}
	log.Printf("题库写入到题库成功，Path: %s", filePath)
}
