package handle

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestAnzalyBodyHead(t *testing.T) {
	fileName := "../1-1.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	bodyByte := []byte{}
	var i = 0
	for {
		line, err := buf.ReadBytes('\n')
		fmt.Println(line)
		if i == 1 {
			line = line[:]
			bodyByte = line
			break
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
		i ++
	}
	err = AnzalyBodyHead(bodyByte)
	if err != nil {
		fmt.Println(err)
	}
}


func TestAnzalyBody(t *testing.T) {
	b := []byte{13,10}
	fmt.Println(string(b))

	fileName := "../2.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var size = stat.Size()
	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)
	bodyByte := make([]byte,110)
	//var i = 0
	//for {

		_,err = buf.Read(bodyByte)

		//line, err := buf.ReadBytes('\n')
		//if i >= 1 && i <= 4{
		//	fmt.Println(string(line))
		//	fmt.Println(string(line[len(line)-1:]))
		//	fmt.Println(len(line))
		//
		//	if len(bodyByte)>0 {
		//		bodyByte = append(bodyByte,line[:len(line)-1]...)
		//	}else{
		//		bodyByte = line[:len(line)-1]
		//	}
		//}
		//if i > 6 {
		//	break
		//}

	//	if err != nil {
	//		if err == io.EOF {
	//			fmt.Println("File read ok!")
	//			break
	//		} else {
	//			fmt.Println("Read file error!", err)
	//			return
	//		}
	//	}
	//	i ++
	//}
	info,err := AnzalyBody(bodyByte)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(info)
}
