package main

import (

	"bufio"
	"fmt"
	"github.com/peterh/liner"
	"goyacc-sql/parser"
	"goyacc-sql/types"
	"os"
	"path/filepath"
	"strings"
	"time"
)


const historyCommmandFile="~/.minisql_history"
const firstPrompt="minisql>"
const secondPrompt="      ->"

func expandPath(path string) (string,error)  {
	if strings.HasPrefix(path, "~/") {
		parts := strings.SplitN(path, "/", 2)
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, parts[1]), nil
	}
	return path, nil
}
func loadHistoryCommand() (*os.File,error) {
	var file *os.File
	path,err:=expandPath(historyCommmandFile)
	if err!=nil{
		return nil,err
	}
	_,err=os.Stat(path)
	if os.IsNotExist(err) {
		file,err=os.Create(path)
		if err!=nil {
			return nil,err
		}
	} else {
		file,err=os.OpenFile(path,os.O_RDWR,0666)
		if err !=nil {
			return nil,err
		}
	}
	return file,err

}

func HandleOneParse(input chan types.DStatements,output chan error)  {
	for item:=range input {
		fmt.Println(item.GetOperationType())
		//DO something you want
		output<-nil // put the error return
	}
	close(output)
}

func runShell(r chan<- error)  {
	ll:=liner.NewLiner()
	defer ll.Close()
	ll.SetCtrlCAborts(true)
	file,err:= loadHistoryCommand()
	if err !=nil {
		panic(err)
	}
	defer func() {
		_,err:=ll.WriteHistory(file)
		if err !=nil{
			panic(err)
		}
		_=file.Close()
	}()
	s:= bufio.NewScanner(file)
	for s.Scan() {
		//fmt.Println(s.Text())
		ll.AppendHistory(s.Text())
	}

	StatementChannel:=make(chan types.DStatements,500)  //用于传输操作指令通道
	FinishChannel:=make(chan error,500) //用于api执行完成反馈通道
	go HandleOneParse(StatementChannel,FinishChannel)  //begin the runtime for exec
	var beginSQLParse=false
	var sqlText=make([]byte,0,100)
	for { //each sql
LOOP:
		beginSQLParse=false
		sqlText=sqlText[:0]
		var input string
		var err error
		for {  //each line
			if beginSQLParse{
				input, err = ll.Prompt(secondPrompt)
			} else {
				input, err = ll.Prompt(firstPrompt)
			}
			if err !=nil{
				if  err ==liner.ErrPromptAborted {
					goto  LOOP
				}
			}
			trimInput:=strings.TrimSpace(input) //get the input without front and backend space
			if len(trimInput)!=0 {
				ll.AppendHistory(input)
				if !beginSQLParse&&(trimInput=="quit"||strings.HasPrefix(trimInput,"quit;")) {
					close(StatementChannel)
					for _=range FinishChannel {

					}
					r<-err
					return
				}
				sqlText=append(sqlText,append([]byte{' '},[]byte(trimInput)[0:]...)...)
				if !beginSQLParse {
					beginSQLParse=true
				}
				if strings.Contains(trimInput,";") {
					break
				}
			}
		}
		beginTime:=time.Now()
		err=parser.Parse(strings.NewReader(string(sqlText)),StatementChannel)
		//fmt.Println(string(sqlText))
		if err != nil {
			fmt.Println(err)
			continue
		}
		<-FinishChannel //等待指令执行完成
		durationTime:=time.Since(beginTime)
		fmt.Println("Finish operation at: ",durationTime)
	}
}
func main() {
	//errChan 用于接收shell返回的err
	errChan:=make(chan error)
	go runShell(errChan) //开启shell协程

	err:=<-errChan
	fmt.Println("bye")
	if err!=nil	 {
		fmt.Println(err)
	}
}