package main

import (
	"conditions/tasks"
	"log"
	"os"
	"sync"
)

var (
	logg   = log.New(os.Stderr, "[INFO] => ", 13)
	logger = log.New(os.Stderr, "[WARNING] => ", 13)
)

func init() {
	logger.Printf("#################Start Client request services...")
}

var (
	authToken = []byte{}
	idleAuth  = make(chan []byte)
	wg        = sync.WaitGroup{}
)

func ligindo(task *tasks.Tasks) {
	// 登录和清理
	authToken = task.AuthTask("")
	// cleanRst := task.BookingPostPageClean(true)
	// logg.Println("clean &:", &cleanRst)
	idleAuth <- authToken
	close(idleAuth)
	wg.Done()
}

func bookdo(task *tasks.Tasks) {
	//预定 和发布
	addRst := task.BookingPostPage(true)
	logg.Println("book &:", &addRst)
	<-idleAuth
	logger.Printf("bookdo authToken:%s \n", authToken)
	wg.Done()
}
func main() {
	allTask := tasks.DoTasks()
	defer allTask.Closer()
	// 清理
	wg.Add(3)
	go func() {
		ligindo(allTask)
	}()

	// BookingPostPage 处理项目
	go func() {
		bookdo(allTask)
	}()

	go func() {
		tasks.ClientPipePath()
	}()
	wg.Wait()
	logger.Printf("#################Stoped Client request services...")
}
