package interval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	loginUrl    = "http://172.22.214.200/ctas/ajaxpro/CExam.Login,App_Web_c318vtbf.ashx"
	userInfoUrl = "http://172.22.214.200/ctas/ajaxpro/CExam.StuIndex,App_Web_c318vtbf.ashx"
	stateUrl    = "http://172.22.214.200/ctas/ajaxpro/CExam.StudReportManage,App_Web_c318vtbf.ashx"
)

type LoginForm struct {
	Username string `json:"strUserId"`
	Password string `json:"strPwd"`
}

// User 用户数
type User struct {
	Username string `json:"CUserName"`  // 用户名
	ClassNo  string `json:"COrganName"` // 班级代号
	Role     string `json:"CUserSort"`  // 角色
}

// State 练习统计数据
type State struct {
	ChapterName string `json:"CCHAPTERNAME"`       // 章节名称
	Total       string `json:"CQUESTIONSUM"`       // 章节的总试题数量
	Read        string `json:"THE_NUMBER_OF"`      // 章节以及练习的试题数量
	Rate        string `json:"THE_CORRECTED_RATE"` // 正确率
}

// GetLoginUser 获取登录用户的账号和密码
func GetLoginUser() (string, string) {
	var username, password string
	fmt.Printf("请输入登录账号：")
	_, _ = fmt.Scan(&username)
	fmt.Printf("请输入登录密码：")
	_, _ = fmt.Scan(&password)
	return username, password
}

// Login 用户登录
func Login(username, password string) string {
	loginForm := &LoginForm{
		Username: username,
		Password: password,
	}
	body, _ := json.Marshal(loginForm)
	request, _ := GenerateCommonRequest("POST", loginUrl, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Ajaxpro-Method", "UserLogin")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("发送登录请求失败，message: %s", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("登录失败，读取登录请求响应体失败, message: %s", err)
	}
	if len(respBody) == 0 || strings.Contains(string(respBody), "true") {
		cookieValue := resp.Header.Get("Set-Cookie")
		index := strings.Index(cookieValue, ";")
		return cookieValue[:index]
	}
	log.Printf("登录失败，请检查账户密码")
	return ""
}

// UserInfo 获取登录用户信息
func UserInfo() *User {
	request, _ := GenerateCommonRequest("POST", userInfoUrl, nil)
	request.Header.Add("X-Ajaxpro-Method", "GetJSONUser")
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("获取用户数据失败，message: %s", err)
		return nil
	}
	defer resp.Body.Close()
	users, _ := ParseResponseBySlice[User](resp, func(body []byte) string {
		return strings.ReplaceAll(string(body[1:len(body)-4]), "\\", "")
	})
	return &users[0]
}

// PractiseState 获取练习统计信息
func PractiseState() ([]State, error) {
	request, _ := GenerateCommonRequest("POST", stateUrl, nil)
	request.Header.Add("X-Ajaxpro-Method", "GetJSONCTRList")
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("获取练习统计数据失败, message: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ParseResponseBySlice[State](resp, func(body []byte) string {
		return strings.ReplaceAll(string(body[1:len(body)-4]), "\\", "")
	})
}
