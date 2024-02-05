/*
	Metadata 생성 초기화

	@author: yhan.lee shiftcats@gmail.com
    @date 2023-11-01
	@version: 0.1.0
*/

package item

import (
	"example.com/message"
	"fmt"
	"github.com/go-faker/faker/v4"
)

var datamap map[message.ServiceCode]*metadata

type supplier func() string

func fakeData(cnt int, supply supplier) *[]string {
	var arrData []string
	for i := 0; i < cnt; i++ {
		arrData = append(arrData, supply())
	}
	return &arrData
}

func LoadData(svc message.ServiceCode) *metadata {
	op := fakeData(30, func() string {
		return fmt.Sprintf("/api/%s/%s", faker.Username(), faker.FirstName())
	})
	ch := fakeData(10, func() string {
		return fmt.Sprintf("%s,%s", faker.URL(), faker.IPv4())
	})
	host := fakeData(10, func() string {
		return fmt.Sprintf("%s,%s", faker.URL(), faker.IPv4())
	})
	dest := fakeData(10, func() string {
		return fmt.Sprintf("%s,%s", faker.URL(), faker.IPv4())
	})
	return NewMetadata(svc, op, ch, host, dest)
}

func InitDataMap() {
	datamap = make(map[message.ServiceCode]*metadata)
	datamap[message.AP064210] = LoadData(message.AP064210)
	datamap[message.OU089012] = LoadData(message.OU089012)
	datamap[message.PG004040] = LoadData(message.PG004040)
	datamap[message.SP016739] = LoadData(message.SP016739)
	datamap[message.AJ054580] = LoadData(message.AJ054580)
}

func RandomItem(code message.ServiceCode) *Item {
	return datamap[code].randomItem()
}
