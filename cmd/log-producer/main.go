/*
	이벤트 로그 발생 프로그램 메인 실행 파일

	@author: yhan.lee shiftcats@gmail.com
    @date 2023-11-01
	@version: 0.1.0
*/

package main

import (
	"context"
	"example.com/message/config"
	"example.com/message/item"
	"example.com/message/producer"
	"example.com/message/util"
	"flag"
	"math/rand"
	"sync"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.yml", "config file")
	flag.StringVar(&configFile, "config", "config.yml", "config file")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
	// 환경설정 파일 로드
	cfg := config.LoadConfig(configFile)

	// 이벤트 로그 데이터 생성
	item.InitDataMap()
	// 카프카 프로듀서 초기화
	pd := producer.NewProducer(&cfg.Kafka)
	defer pd.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	util.ShutdownHook(cancel)

	// 이벤트 로그 generator에서 생성한 메시지를 카프카 producer로 전달하기 위해 채널을 사용한다.
	ch := make(chan *producer.Message, cfg.NumberOfWorker)
	defer close(ch)

	wg := sync.WaitGroup{}
	wg.Add(cfg.NumberOfWorker + 1)
	// 카프카로 메시지 전달을 위한 Producer 실행
	go func() {
		pd.Produce(ctx, ch)
		wg.Done()
	}()
	// 환경 설정에 지정한 수 만큼 GO 루틴 생성
	for i := 0; i < cfg.NumberOfWorker; i++ {
		go func() {
			item.GenerateLog(ctx, rand.Intn(300)+700, ch)
			wg.Done()
		}()
	}
	wg.Wait()
}
