package interval

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"
)

const listChapterUrl = "http://172.22.214.200/ctas/ajaxpro/CExam.CPractice,App_Web_tzfdzrj8.ashx"
const listProgramUrl = "http://172.22.214.200/ctas/ajaxpro/CExam.CPractice,App_Web_tzfdzrj8.ashx"

// Chapter 练习平台章节数据
type Chapter struct {
	ChapterId   string `json:"CChapterID"`
	ChapterName string `json:"CChapterName"`
}

// Program 练习平台章节程序数据
type Program struct {
	ProgramId string `json:"CProgramID"`
}

// ReplyChapterTopic 练习某个章节的试题
func ReplyChapterTopic(chapterId string, banks *map[string]*StorageTopic) {
	programs := getProgramsByChapterId(chapterId)
	if programs == nil || len(programs) == 0 {
		log.Printf("章节 %s 的程序列表为空或以达到最大重试次数，跳过此章节！", chapterId)
		return
	}
	for _, program := range programs {
		log.Printf("章节 %s 程序 %s，刷题开始", chapterId, program.ProgramId)
		count := ReplyProgramTopic(program.ProgramId, banks)
		log.Printf("章节 %s 程序 %s，刷题结束，共刷 %d 题", chapterId, program.ProgramId, count+1)
		time.Sleep(2 * time.Second)
	}
}

// ReplyProgramTopic 练习某个程序的试题
func ReplyProgramTopic(programId string, banks *map[string]*StorageTopic) int {
	index := 0
	for {
		topic := getTopicByProgramIdAndIndex(programId, index)
		if topic == nil {
			log.Printf("获取程序 %s 的第 %d 道试题信息失败, 跳过此题", programId, index+1)
			index++
			continue
		}
		storage, ok := (*banks)[topic.TopicId]
		if !ok {
			log.Printf("题库中不存在试题 %s 答案，开始获取题解", topic.TopicId)
			if answer := getTopicAnswer(topic.TopicId); answer != "" {
				log.Printf("获取试题 %s 答案成功，添加到题库", topic.TopicId)
				storage = &StorageTopic{
					Id:      topic.TopicId,
					Content: topic.TopicContent,
					Answers: topic.TopicAnswer,
					Answer:  answer,
				}
				(*banks)[topic.TopicId] = storage
			}
		}
		if storage != nil && storage.Answer != "" {
			result, err := submitTopicAnswer(topic.TopicId, storage.Answer, true)
			if err != nil {
				log.Printf("提交试题 %s 答案失败", topic.TopicId)
			} else {
				log.Printf("试题 %s , 题目内容：%s, 可选答案：%s，题库答案：%s，提交成功，结果：%v", topic.TopicId, topic.TopicContent, topic.TopicAnswer, storage.Answer, result)
			}
		}
		count, _ := strconv.Atoi(topic.TopicCount)
		if index >= count-1 {
			break
		}
		index++
		time.Sleep(100 * time.Millisecond)
	}
	return index
}

// GetChapterAnswer 获取某个章节的试题答案
func GetChapterAnswer(chapterId string, banks *map[string]*StorageTopic) {
	programs := getProgramsByChapterId(chapterId)
	if programs == nil || len(programs) == 0 {
		log.Printf("章节 %s 的程序列表为空或以达到最大重试次数，跳过此章节！", chapterId)
		return
	}
	log.Printf("章节 %s 程序列表获取成功，Size: %d", chapterId, len(programs))
	for _, program := range programs {
		log.Printf("开始获取章节 %s 程序 %s 的试题答案", chapterId, program.ProgramId)
		count := GetProgramAnswer(program.ProgramId, banks)
		log.Printf("章节 %s 程序 %s 获取试题答案完成，成功获取 %d 道题", chapterId, program.ProgramId, count+1)
	}
}

// GetProgramAnswer 获取某个章节中指定程序内的所有答案
func GetProgramAnswer(programId string, banks *map[string]*StorageTopic) int {
	index, errCount := 0, 0
	for {
		if errCount >= 3 {
			log.Printf("程序 %s 连续获取试题失败超过3道，跳过处理", programId)
			break
		}
		topic := getTopicByProgramIdAndIndex(programId, index)
		if topic == nil {
			errCount++
			index++
			log.Printf("获取程序 %s 的第 %d 道试题信息失败, 已连续获取失败 %d 道题", programId, index+1, errCount)
			continue
		}
		errCount = 0
		if _, ok := (*banks)[topic.TopicId]; ok {
			log.Printf("试题 %s 答案题库内已存在，跳过获取！", topic.TopicId)
		} else {
			log.Printf("开始获取试题 %s 的答案, 试题内容: %s, 全部答案: %s", topic.TopicId, topic.TopicContent, topic.TopicAnswer)
			answer := getTopicAnswer(topic.TopicId)
			if answer == "" {
				log.Printf("获取试题 %s 答案失败，跳过此题", topic.TopicId)
			} else {
				log.Printf("获取试题 %s 答案成功, 答案: %s", topic.TopicId, answer)
			}
			storage := &StorageTopic{
				Id:      topic.TopicId,
				Content: topic.TopicContent,
				Answers: topic.TopicAnswer,
				Answer:  answer,
			}
			(*banks)[topic.TopicId] = storage
		}
		count, _ := strconv.Atoi(topic.TopicCount)
		// 当前程序的题目以及获取完成 直接跳出
		if index >= count-1 {
			break
		}
		index++
		time.Sleep(100 * time.Millisecond)
	}
	return index
}

// getTopicByProgramIdAndIndex 通过程序id和试题小标获取试题信息 如果获取失败 在方法内部会重试三次
func getTopicByProgramIdAndIndex(programId string, index int) *Topic {
	for i := 0; i < 3; i++ {
		topic, err := queryTopic(programId, index)
		if err != nil {
			log.Printf("获取程序 %s 的第 %d 道试题信息失败，已尝试获取 %d 次", programId, index+1, i+1)
		} else {
			return topic
		}
	}
	return nil
}

// 通过章节id获取章节的程序目录 如果获取失败 在方法内部会重试三次
func getProgramsByChapterId(chapterId string) []Program {
	for i := 0; i < 3; i++ {
		programs, err := queryProgram(chapterId)
		if err != nil {
			log.Printf("获取章节 %s 的程序列表失败，已尝试获取 %d 次，message: %s", chapterId, i+1, err)
		} else {
			return programs
		}
	}
	return nil
}

// QueryChapters 查询章节列表
func QueryChapters() ([]Chapter, error) {
	request, _ := GenerateCommonRequest("POST", listChapterUrl, nil)
	request.Header.Set("X-Ajaxpro-Method", "GetJSONChapterList")
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("请求章节列表失败, message: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ParseResponseBySlice[Chapter](resp, func(body []byte) string {
		return strings.ReplaceAll(string(body[1:len(body)-5]), "\\", "")
	})
}

// 查询某个章节的程序列表
func queryProgram(chapterId string) ([]Program, error) {
	requestBodyString := "{\"cChapterID\": \"" + chapterId + "\"}"
	request, _ := GenerateCommonRequest("POST", listProgramUrl, bytes.NewBufferString(requestBodyString))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Ajaxpro-Method", "GetJSONProgramList")
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("请求章节程序列表失败，message: %s, chapterId: %s", err, chapterId)
		return nil, err
	}
	defer resp.Body.Close()
	return ParseResponseBySlice[Program](resp, func(body []byte) string {
		return strings.ReplaceAll(string(body[1:len(body)-4]), "\\", "")
	})
}
