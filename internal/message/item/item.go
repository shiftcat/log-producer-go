/*
	이벤트 로그 데이터 모델

	@author: yhan.lee shiftcats@gmail.com
    @date 2023-11-01
	@version: 0.1.0
*/

package item

import (
	"bytes"
	"example.com/message"
	"example.com/message/util"
	"fmt"
	"github.com/go-faker/faker/v4"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type Caller struct {
	Channel   string `json:"channel"`
	ChannelIp string `json:"channelIp"`
}

type Host struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type Response struct {
	Type     message.ResponseType `json:"type"`
	Status   int                  `json:"status"`
	Duration int64                `json:"duration"`
}

type User struct {
	Id    string `faker:"first_name" json:"id"`
	Ip    string `faker:"ipv4" json:"ip"`
	Agent string `faker:"oneof: web, ios, android" json:"agent"`
}

type ItemType struct {
	Type    message.LogType `json:"log_type"`
	TraceId string          `json:"trace_id"`
	SpanId  string          `json:"span_id"`
}

type Item struct {
	ItemType
	Svc           message.ServiceCode `json:"service"`
	Operation     string              `json:"operation"`
	Caller        *Caller             `json:"caller"`
	Host          *Host               `json:"host"`
	Destination   *Host               `json:"destination"`
	User          *User               `json:"user"`
	EventDatetime int64               `json:"event_dt"`
	duration      int64               `json:"-"`
	indent        int                 `json:"-"`
	refItem       *Item               `json:"-"`
	spans         []*Item
}

type ResItem struct {
	*Item
	Response *Response `json:"response"`
}

func NewItem(svc message.ServiceCode, op string, caller *Caller, host *Host, dest *Host) *Item {
	itm := new(Item)
	itm.Svc = svc
	itm.Operation = op
	itm.Caller = caller
	itm.Host = host
	itm.Destination = dest
	itm.indent = 0
	return itm
}

func NewFakeItem(svc message.ServiceCode) *Item {
	op := fmt.Sprintf("/rest/%s/%s", faker.Username(), faker.FirstName())
	caller := &Caller{
		Channel:   faker.URL(),
		ChannelIp: faker.IPv4(),
	}
	host := &Host{
		Name: faker.DomainName(),
		Ip:   faker.IPv4(),
	}
	dest := &Host{
		Name: faker.Username(),
		Ip:   faker.IPv4(),
	}
	return NewItem(svc, op, caller, host, dest)
}

func (i *Item) AddSpan(span *Item) {
	i.spans = append(i.spans, span)
}

func (i *Item) copy() *Item {
	item := NewItem(i.Svc, i.Operation, i.Caller, i.Host, i.Destination)
	item.TraceId = i.TraceId
	item.SpanId = i.SpanId
	item.indent = i.indent
	return item
}

func fakeUser() *User {
	user := &User{}
	err := faker.FakeData(user)
	if err != nil {
		log.Printf("%v\n", err)
	}
	// user.Id = faker.LastName()
	return user
}

// PairingItem REQ/RES 분할된 새로운 아이템 배열 생성
func (i *Item) PairingItem() [2]*Item {
	fakeUser := fakeUser()
	reqItem := i.copy()
	reqItem.Type = message.REQ
	reqItem.User = fakeUser
	resItem := i.copy()
	resItem.Type = message.RES
	resItem.User = fakeUser
	resItem.duration = i.duration
	// RES에서 REQ를 참조할 수 있게한다.
	resItem.refItem = reqItem
	items := [2]*Item{reqItem, resItem}
	return items
}

func (i *Item) GetSpans() []*Item {
	return i.spans
}

func (i *Item) GetIndent() int {
	return i.indent
}

func pushItems(inner []*Item, outer []*Item) []*Item {
	outerSize := len(outer)
	reminder := outerSize % 2
	if reminder != 0 {
		panic("The outer size must be an even.")
	}
	mp := outerSize / 2
	ps := append(outer[:mp], append(inner[:], outer[mp:]...)...)
	return ps
}

// SpanItems 현재 아이템의 모든 spans 아이템의 참조 링크를 1차원 배열로 생성
// 이 때 하나의 아이템을 REQ/RES로 분할하며, REQ/RES 배열 사이에 span 아이템을 삽입한다.
func (i *Item) SpanItems() []*Item {
	spans := i.GetSpans()
	i.duration = int64(rand.Intn(2300) + 700)
	pairingItem := i.PairingItem()
	if spans == nil {
		return pairingItem[:]
	}
	if len(spans) < 1 {
		return pairingItem[:]
	}

	var accumulateDuration int64 = 0
	var accumulateItems []*Item
	for _, s := range spans {
		spanItems := s.SpanItems()
		accumulateDuration = s.duration + accumulateDuration
		accumulateItems = append(accumulateItems, spanItems...)
	}
	i.duration = accumulateDuration + int64(rand.Intn(500)+200)
	pairingItem[1].duration = i.duration // 응답에만 duration 반영

	return pushItems(accumulateItems, pairingItem[:])
}

func bufferedStringJoin(id ...int) string {
	var b bytes.Buffer
	for i, v := range id {
		if 0 == i {
			b.WriteString(fmt.Sprintf("%v", v))
		} else {
			b.WriteString(fmt.Sprintf(",%v", v))
		}
	}
	return b.String()
}

func (i *Item) InitTraceId(traceId string, start ...int) {
	i.TraceId = traceId
	i.SpanId = bufferedStringJoin(start...)
	i.indent = (len(start) - 1) * 2
	if i.spans == nil {
		return
	}
	if len(i.spans) == 0 {
		return
	}
	spans := i.GetSpans()
	tailN := start[len(start)-1]
	for _, s := range spans {
		tailN = tailN + 1
		ids := append(start, tailN)
		s.InitTraceId(traceId, ids...)
	}
}

// String Json 문열로 반환
func (i *Item) String() string {
	return util.ToJson(i)
}

// IndentedString 들여쓰기의 Json 문자열 반환
func (i *Item) IndentedString() string {
	return fmt.Sprint(strings.Repeat(" ", i.indent), i.String())
}

var status = []int{
	http.StatusOK, http.StatusCreated, http.StatusNoContent, http.StatusNotFound,
	http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError,
	http.StatusFound, http.StatusPermanentRedirect, http.StatusConflict, http.StatusForbidden,
}

func randomStatus() int {
	ri := rand.Intn(10)
	if (ri % 10) == 0 {
		ind := len(status)
		return status[rand.Intn(ind)]
	} else {
		return status[rand.Intn(3)]
	}
}

func responseType(status int) message.ResponseType {
	res := status / 100
	switch res {
	case 2:
		return message.S
	case 3:
		return message.I
	case 4, 5:
		return message.E
	default:
		return message.E
	}
}

func (i *Item) ToRes() *ResItem {
	resStatus := randomStatus()
	r := &Response{Type: responseType(resStatus), Status: resStatus, Duration: i.duration}
	res := new(ResItem)
	res.Item = i
	res.Response = r
	return res
}

// String Json 문열로 반환
func (r *ResItem) String() string {
	return util.ToJson(r)
}

// IndentedString 들여쓰기의 Json 문자열 반환
func (r *ResItem) IndentedString() string {
	return fmt.Sprint(strings.Repeat(" ", r.indent), r.String())
}
