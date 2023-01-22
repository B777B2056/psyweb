package utils

import (
	"io"
	"log"
	"os/exec"
	"sync"
)

/* 深度学习子进程相关操作 */
type MsgHandler func([]byte)

type deepLearning struct {
	cmd     *exec.Cmd
	reader  io.ReadCloser
	writer  io.WriteCloser
	handler MsgHandler
	wg      sync.WaitGroup
}

func (dl *deepLearning) RegistMsgHandler(handler MsgHandler) {
	dl.handler = handler
}

func (dl *deepLearning) Start() (err error) {
	// 成员初始化
	dl.cmd = exec.Command("python3", "deeplearning/main.py") // 运行画图脚本
	dl.reader, err = dl.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	dl.writer, err = dl.cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = dl.cmd.Start()
	if err != nil {
		return err
	}
	// 启动读协程
	dl.wg.Add(1)
	go func() {
		for {
			dl.handleRead()
		}
	}()
	return err
}

func (dl *deepLearning) handleRead() {
	var msg []byte
	for {
		char := make([]byte, 1)
		n, err := dl.reader.Read(char)
		if err != nil {
			log.Println(err)
			return
		}
		if n == len(char) {
			if char[0] == '\n' {
				dl.handler(msg)
				break
			} else {
				msg = append(msg, char[0])
			}
		}
	}
}

func (dl *deepLearning) handleWrite(msg []byte) error {
	for n := 0; n < len(msg); {
		cur_n, err := dl.writer.Write(msg[n:])
		if err != nil {
			return err
		}
		n += cur_n
	}
	return nil
}

func (dl *deepLearning) Do(eeg_data_path string) error {
	msg := eeg_data_path
	msg += "\n"
	return dl.handleWrite([]byte(msg))
}

func (dl *deepLearning) Stop() error {
	msg := "stop\n"
	err := dl.handleWrite([]byte(msg))
	dl.wg.Wait()
	return err
}

var deepLearningInstance *deepLearning
var deepLearningOnce sync.Once

func GetDeepLearningInstance() *deepLearning {
	deepLearningOnce.Do(func() {
		deepLearningInstance = &deepLearning{}
	})
	return deepLearningInstance
}
