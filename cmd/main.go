package main

import (
	"cqupt-learning-platform/interval"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	username, password := interval.GetLoginUser()
	sessionValue := interval.Login(username, password)
	if sessionValue == "" {
		return
	}
	interval.SessionId = sessionValue
	if user := interval.UserInfo(); user != nil {
		log.Printf("欢迎使用\t %s \t %s \t %s", user.ClassNo, user.Username, user.Role)
	}
	for {
		fmt.Println("-----------------------(～￣▽￣)～-----------------------")
		fmt.Printf("1.获取试题答案到题库\n2.开始日常练习\n3.查询练习统计信息\n")
		fmt.Printf("请输入数字选择对应功能，输入 C 退出：")
		var input string
		if _, err := fmt.Scan(&input); err != nil {
			log.Printf("输入数据错误，请重试...")
			continue
		}
		switch strings.ToLower(input) {
		case "1":
			getAllTopicAnswer()
			break
		case "2":
			startReplyTopic()
		case "3":
			states, err := interval.PractiseState()
			if err != nil {
				log.Printf("获取练习统计信息失败, message: %s", err)
				break
			}
			if len(states) == 0 {
				log.Printf("当前还没有练习统计信息")
			} else {
				for _, state := range states {
					fmt.Printf("%s \t 试题数：%s \t 练习题数：%s \t 正确率：%s \n", state.ChapterName, state.Total, state.Read, state.Rate)
				}
			}
			break
		case "c":
			os.Exit(0)
		default:
			log.Println("输入数据有误，请重新输入...")
		}
	}
}

// 开始练习
func startReplyTopic() {
	chapters, err := interval.QueryChapters()
	if err != nil {
		log.Printf("获取章节列表失败，请重试！")
		return
	}
	banks := interval.ReadAnswerFile()
	if banks == nil {
		log.Printf("读取题库文件错误，请先获取试题答案或重试！")
		return
	}
	oldBankSize := len(banks)
	for i := selectChapter(chapters); i < len(chapters); i++ {
		chapter := chapters[i]
		log.Printf("------------章节 %s 开始刷题------------", chapter.ChapterName)
		interval.ReplyChapterTopic(chapter.ChapterId, &banks)
		log.Printf("------------章节 %s 完成------------", chapter.ChapterName)
		if i < len(chapters)-1 && !isNext() {
			break
		}
	}
	log.Printf("刷题结束")
	// 如果题库发生了变化 那么重新保存题库
	if len(banks) != oldBankSize {
		interval.SaveAnswersToFile(banks)
	}
}

// 获取全部试题答案
func getAllTopicAnswer() {
	chapters, err := interval.QueryChapters()
	if err != nil {
		log.Printf("获取章节列表失败，请重试！")
		return
	}
	log.Printf("------------开始获取试题答案------------")
	banks := interval.ReadAnswerFile()
	if banks == nil {
		log.Println("读取题库文件错误，请删除或创建后重试！")
		return
	}
	for i := selectChapter(chapters); i < len(chapters); i++ {
		chapter := chapters[i]
		log.Printf("开始获取章节 %s %s 的试题答案", chapter.ChapterId, chapter.ChapterName)
		interval.GetChapterAnswer(chapter.ChapterId, &banks)
		log.Printf("章节 %s 答案获取完成", chapter.ChapterName)
		if i < len(chapters)-1 && !isNext() {
			break
		}
	}
	// 将题库保存到本地文件中
	interval.SaveAnswersToFile(banks)
	log.Printf("------------所有试题答案获取完成------------")
}

// 是否继续下一章
func isNext() bool {
	fmt.Printf("当前章节已完成，是否继续下一章节？（输入 1 继续）")
	var input string
	_, _ = fmt.Scan(&input)
	return input == "1"
}

// 选择开始章节
func selectChapter(chapters []interval.Chapter) int {
	for i, chapter := range chapters {
		fmt.Printf("%d. %s\n", i, chapter.ChapterName)
	}
	startIndex := -1
	for {
		fmt.Printf("当前共有 %d 个章节，请输入起始章节下标：", len(chapters))
		if _, err := fmt.Scan(&startIndex); err != nil || startIndex < 0 || startIndex >= len(chapters) {
			log.Printf("输入数据不合法, 请重新输入!")
		} else {
			break
		}
	}
	return startIndex
}
