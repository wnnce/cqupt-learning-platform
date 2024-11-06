package interval

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	queryTopicUrl   = "http://172.22.214.200/ctas/ajaxpro/CExam.CPractice,App_Web_tzfdzrj8.ashx"
	submitAnswerUrl = "http://172.22.214.200/ctas/ajaxpro/CExam.CPractice,App_Web_tzfdzrj8.ashx"
)

// Topic 练习试题信息
type Topic struct {
	TopicId      string `json:"CQuestionID"`      // 试题Id
	TopicContent string `json:"CQuestionContent"` // 试题题目内容
	TopicCount   string `json:"CQuestionCount"`   // 当前程序的试题数量
	Level        string `json:"CDiffLevel"`       // 试题的难度级别
	TopicAnswer  string // 试题的答案列表
}

type Result struct {
	Msg   string `json:"msg"`
	Topic *Topic `json:"CQuestion"`
}

// 查询题目
func queryTopic(programId string, topicIndex int) (*Topic, error) {
	requestBody := "{\"strTestParam\":\"<cTest><cProgram>" + programId + "</cProgram><cQuestionIndex>" + strconv.Itoa(topicIndex) + "</cQuestionIndex></cTest>\"}"
	request, _ := GenerateCommonRequest("POST", queryTopicUrl, bytes.NewBufferString(requestBody))
	request.Header.Add("X-Ajaxpro-Method", "GetJSONTest")
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("请求程序题目失败, programId: %s, index: %d, message: %s", programId, topicIndex, err)
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ParseResponseByStruct[Result](resp, func(body []byte) string {
		return strings.ReplaceAll(string(body[1:len(body)-4]), "\\", "")
	})
	if err != nil {
		return nil, err
	}
	topicContent := strings.ReplaceAll(result.Topic.TopicContent, "&lt;br&gt;", "")
	answerIndex := strings.LastIndex(topicContent, "A")
	result.Topic.TopicContent = topicContent[:answerIndex]
	result.Topic.TopicAnswer = topicContent[answerIndex:]
	return result.Topic, nil
}

// 提交题目答案
func submitTopicAnswer(topicId, answer string) (bool, error) {
	requestBody := "{\"strTestParam\":\"<cTestParam><cQuestion>" + topicId + "</cQuestion><cUserAnswer>" + answer + "</cUserAnswer></cTestParam>\"}"
	request, _ := GenerateCommonRequest("POST", submitAnswerUrl, bytes.NewBufferString(requestBody))
	request.Header.Add("X-Ajaxpro-Method", "IsOrNotTrue")
	resp, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	if bodyString[0] == '1' {
		return true, nil
	}
	return false, nil
}

// 通过穷举获取题目的答案
func getTopicAnswer(topicId string) string {
	for _, answer := range []string{"A", "B", "C", "D"} {
		result, err := submitTopicAnswer(topicId, answer)
		if err != nil {
			log.Fatalf("判断 %s 的答案 %s 是否正确错误，message: %s，跳过当前答案。", topicId, answer, err)
		}
		if result {
			return answer
		}
		// 休眠20ms
		time.Sleep(20 * time.Millisecond)
	}
	return ""
}
