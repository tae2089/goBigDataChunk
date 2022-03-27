package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	// file open

	f, err := os.Open("access.log")
	defer f.Close()
	if err != nil {
		log.Fatal(err)

	}
	start := time.Now()

	err = Process(f)

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

func Process(f *os.File) error {

	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 500*1024)
		return lines
	}}
	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}
	slicePool := sync.Pool{New: func() interface{} {
		lines := make([]string, 100)
		return lines
	}}
	//var wg sync.WaitGroup
	r := bufio.NewReader(f)
	d := 0
	//mutex := new(sync.Mutex)
	for {
		buf := linesPool.Get().([]byte)
		n, err := r.Read(buf)
		buf = buf[:n]
		//fmt.Println(n)
		if n == 0 {
			if err != nil {
				fmt.Println(err)
				break
			}

			if err == io.EOF {
				fmt.Println(err)
				break
			}
			return err
		}
		//linesPool을 통해 받은 데이터에 끝부분이 \n아닐 수도 있기에 추가적으로 읽기 진행
		nextUntillNewline, err := r.ReadBytes('\n')
		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}
		//wg.Add(1)
		d += ProcessChunk(buf, &linesPool, &stringPool, &slicePool)
		//mutex.Lock()
		//go func() {
		//	d += ProcessChunk(buf, &linesPool, &stringPool, &slicePool)
		//	wg.Done()
		//}()
		//mutex.Unlock()
	}
	//wg.Wait()

	fmt.Println(d)

	return nil
}

func ProcessChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool, slicePool *sync.Pool) int {

	//another wait group to process every chunk further
	//var wg2 sync.WaitGroup
	//stringPool 불러오기
	logs := stringPool.Get().(string)
	//청크 데이터를 문자열로 변환
	logs = string(chunk)
	linesPool.Put(chunk) //put back the chunk in pool
	//slicePool 불러오기
	logsSlice := slicePool.Get().([]string)
	//\n을 기준으로 string 배열 생성하기
	logsSlice = strings.Split(logs, "\n")
	//stringPool 반환
	stringPool.Put(logs)
	//100개의 로그만 읽기위한 chunksize 지정
	chunkSize := 100
	length := len(logsSlice)
	//traverse the chunk

	for i := 0; i < length; i += chunkSize {
		//wg2.Add(1)
		//청크 계산하기
		s := i * chunkSize
		e := int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice))))
		for i := s; i < e; i++ {
			text := logsSlice[i]
			if len(text) == 0 {
				continue
			}
		}
		//go func(s int, e int) {
		//	for i := s; i < e; i++ {
		//		text := logsSlice[i]
		//		if len(text) == 0 {
		//			continue
		//		}
		//
		//		/*
		//			여기에 전처리 코드를 입력하시면 됩니다:)
		//		*/
		//	}
		//	wg2.Done()
		//}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))

	}
	//청크가 다 끝날때까지 기달리기
	//wg2.Wait()
	// slicePool 반환
	slicePool.Put(logsSlice)
	return 1
}
