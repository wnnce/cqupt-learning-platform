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

const (
	fileFolder = "./data/"
	fileName   = "banks.json"
)

// ReadAnswerFile 读取保存在本地的题库文件
func ReadAnswerFile() map[string]*StorageTopic {
	readFile, err := os.ReadFile(fileFolder + fileName)
	banks := make(map[string]*StorageTopic)
	if os.IsNotExist(err) {
		log.Println("题库文件不存在，返回空题库!")
		return banks
	} else if err != nil {
		log.Printf("读取题库文件错误，message: %s", err)
		return nil
	}
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
	saveFolder := fileFolder
	_, err := os.Stat(saveFolder)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(fileFolder, 0755); err != nil {
			log.Printf("创建 %s 文件夹失败，请手动创建，已经题库文件保存到当前目录!", fileFolder)
			saveFolder = ""
		}
	}
	// 最终的保存路径
	finalSavePath := saveFolder + fileName
	if err := os.WriteFile(finalSavePath, jsonValue, 0644); err != nil {
		log.Printf("写入题库到文件错误，message: %s", err)
	}
	log.Printf("题库写入到题库成功，Path: %s", finalSavePath)
}
