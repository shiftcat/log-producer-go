/*
	유형별 이벤트 로그 데이터 생성

	@author: yhan.lee shiftcats@gmail.com
    @date 2023-11-01
	@version: 0.1.0
*/

package item

import (
	"context"
	"example.com/message"
	"example.com/message/producer"
	"example.com/message/util"
	"fmt"
	"github.com/openzipkin/zipkin-go/idgenerator"
	"math/rand"
	"time"
)

type CaseFunction func() []*Item

var caseFunctions []CaseFunction

var svcs []message.ServiceCode

// 초기화 자동 호출
func init() {
	svcs = []message.ServiceCode{message.AP064210, message.PG004040, message.AJ054580, message.OU089012, message.SP016739}
	// var caseFunctions []CaseFunction
	caseFunctions = append(caseFunctions, func() []*Item {
		return Case1()
	})
	caseFunctions = append(caseFunctions, func() []*Item {
		return Case2()
	})
	caseFunctions = append(caseFunctions, func() []*Item {
		return Case3()
	})
	caseFunctions = append(caseFunctions, func() []*Item {
		return Case4()
	})
	caseFunctions = append(caseFunctions, func() []*Item {
		return Case5()
	})
	caseFunctions = append(caseFunctions, func() []*Item {
		return Fake1()
	})
}

// RandSvc 서비스 코드를 인자로 주어진 갯수 만큼 중복되지 않게 랜덤 생성하여 배열로 리턴
func RandSvc(cnt int) []message.ServiceCode {
	svcMap := make(map[message.ServiceCode]struct{})
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	list := util.NewList(append([]message.ServiceCode{}, svcs...))
	for {
		size := list.Size()
		ri := r.Intn(size)
		svc := list.Take(ri)
		svcMap[svc] = struct{}{}
		if len(svcMap) >= cnt {
			break
		}
	}

	// 맵의 키를 리스트로 변환
	var result []message.ServiceCode
	for k := range svcMap {
		result = append(result, k)
	}
	return result
}

func Case1() []*Item {
	randSvc := RandSvc(5)
	item1 := RandomItem(randSvc[0])
	item2 := RandomItem(randSvc[1])
	item3 := RandomItem(randSvc[2])
	item4 := RandomItem(randSvc[3])
	item5 := RandomItem(randSvc[4])

	item1.AddSpan(item2)
	item1.AddSpan(item3)
	item3.AddSpan(item4)
	item4.AddSpan(item5)

	var ig = idgenerator.NewRandom128()
	traceID := ig.TraceID()
	item1.InitTraceId(traceID.String(), 1)

	return item1.SpanItems()
}

func Case2() []*Item {
	randSvc := RandSvc(3)
	item1 := RandomItem(randSvc[0])
	item2 := RandomItem(randSvc[1])
	item3 := RandomItem(randSvc[2])

	item2.AddSpan(item3)
	item1.AddSpan(item2)

	var ig = idgenerator.NewRandom128()
	traceID := ig.TraceID()
	item1.InitTraceId(traceID.String(), 0)
	return item1.SpanItems()
}

func Case3() []*Item {
	randSvc := RandSvc(4)
	item1 := RandomItem(randSvc[0])
	item2 := RandomItem(randSvc[1])
	item3 := RandomItem(randSvc[2])
	item4 := RandomItem(randSvc[3])

	item1.AddSpan(item2)
	item1.AddSpan(item3)
	item1.AddSpan(item4)

	var ig = idgenerator.NewRandom128()
	traceID := ig.TraceID()
	item1.InitTraceId(traceID.String(), 0)
	return item1.SpanItems()
}

func Case4() []*Item {
	randSvc := RandSvc(1)
	item1 := RandomItem(randSvc[0])
	var ig = idgenerator.NewRandom128()
	traceID := ig.TraceID()
	item1.InitTraceId(traceID.String(), 0)
	return item1.SpanItems()
}

func Case5() []*Item {
	randSvc := RandSvc(5)
	item1 := RandomItem(randSvc[0])
	item2 := RandomItem(randSvc[1])
	item3 := RandomItem(randSvc[2])
	item4 := RandomItem(randSvc[3])
	item5 := RandomItem(randSvc[4])

	item1.AddSpan(item2)
	item1.AddSpan(item3)
	item3.AddSpan(item4)
	item1.AddSpan(item5)

	var ig = idgenerator.NewRandom128()
	traceID := ig.TraceID()
	item1.InitTraceId(traceID.String(), 1)

	return item1.SpanItems()
}

func Fake1() []*Item {
	randSvc := RandSvc(5)
	item1 := NewFakeItem(randSvc[0])
	item2 := NewFakeItem(randSvc[1])
	item3 := NewFakeItem(randSvc[2])
	item4 := NewFakeItem(randSvc[3])
	item5 := NewFakeItem(randSvc[4])

	item1.AddSpan(item2)
	item2.AddSpan(item3)
	item1.AddSpan(item4)
	item1.AddSpan(item5)

	var ig = idgenerator.NewRandom128()
	traceID := ig.TraceID()
	item1.InitTraceId(traceID.String(), 1)

	return item1.SpanItems()
}

func SendItems(ctx context.Context, maxSleep int, ch chan<- *producer.Message, items []*Item) bool {
	for _, item := range items {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped log publish")
			return false
		default:
			msg := new(producer.Message)
			msg.Key = []byte(item.TraceId)

			if item.Type == message.REQ {
				randSleep := rand.Intn(maxSleep)
				time.Sleep(time.Millisecond * time.Duration(randSleep))
				item.EventDatetime = time.Now().UnixMilli()
				msg.Value = []byte(item.String())
			} else if item.Type == message.RES {
				reqItem := item.refItem
				resTimestamp := reqItem.EventDatetime + item.duration
				dft := resTimestamp - time.Now().UnixMilli()
				if dft > 0 {
					time.Sleep(time.Millisecond * time.Duration(dft))
				}
				// time.Sleep(time.Millisecond * time.Duration(item.duration))
				// item.EventDatetime = time.Now().UnixMilli()
				item.EventDatetime = resTimestamp
				msg.Value = []byte(item.ToRes().String())
			}
			// fmt.Println(item.IndentedString())

			ch <- msg
		}
	}
	return true
}

func GenerateLog(ctx context.Context, maxSleep int, ch chan<- *producer.Message) {
	for {
		ri := rand.Intn(len(caseFunctions))
		// 특정 유형의 이벤트 로그 데이터 생성
		items := caseFunctions[ri]()
		if !SendItems(ctx, maxSleep, ch, items) {
			break
		}
	}
}
