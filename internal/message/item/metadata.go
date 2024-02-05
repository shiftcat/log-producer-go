/*
	이벤트 로그 데이터 생성을 위한 데이터 모델

	이벤트 로그 생성에 필요한 Operation, Channel, Host, Destination 데이터를 관리한다.

	@author: yhan.lee shiftcats@gmail.com
    @date 2023-11-01
	@version: 0.1.0
*/

package item

import (
	"example.com/message"
	"math/rand"
	"strings"
)

func NewMetadata(svc message.ServiceCode, op *[]string, ch *[]string, host *[]string, dest *[]string) *metadata {
	i := new(metadata)
	i.svc = svc
	i.op = op
	i.ch = ch
	i.host = host
	i.dest = dest
	return i
}

type metadata struct {
	svc  message.ServiceCode
	op   *[]string
	ch   *[]string
	host *[]string
	dest *[]string
}

func (m *metadata) getOperation() string {
	var size int = len(*m.op)
	loc := rand.Intn(size)
	return (*m.op)[loc]
}

func (m *metadata) getCaller() *Caller {
	var size = len(*m.ch)
	loc := rand.Intn(size)
	line := (*m.ch)[loc]
	split := strings.Split(line, ",")
	if len(split) > 1 {
		return &Caller{Channel: strings.Trim(split[0], " "), ChannelIp: strings.Trim(split[1], " ")}
	} else {
		return &Caller{Channel: strings.Trim(split[0], " "), ChannelIp: ""}
	}
}

func (m *metadata) getHost() *Host {
	size := len(*m.host)
	loc := rand.Intn(size)
	line := (*m.host)[loc]
	split := strings.Split(line, ",")
	if len(split) > 1 {
		return &Host{Name: strings.Trim(split[0], " "), Ip: strings.Trim(split[1], " ")}
	} else {
		return &Host{Name: strings.Trim(split[0], " "), Ip: ""}
	}
}

func (m *metadata) getDestination() *Host {
	size := len(*m.dest)
	loc := rand.Intn(size)
	line := (*m.dest)[loc]
	split := strings.Split(line, ",")
	if len(split) > 1 {
		return &Host{Name: strings.Trim(split[0], " "), Ip: strings.Trim(split[1], " ")}
	} else {
		return &Host{Name: strings.Trim(split[0], " "), Ip: ""}
	}
}

func (m *metadata) randomItem() *Item {
	itm := new(Item)
	itm.Svc = m.svc
	itm.Operation = m.getOperation()
	itm.Caller = m.getCaller()
	itm.Host = m.getHost()
	itm.Destination = m.getDestination()
	return itm
}
