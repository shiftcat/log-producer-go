package item

import (
	"example.com/message"
	"example.com/message/util"
	"fmt"
	"github.com/go-faker/faker/v4"
	"testing"
)

func TestItem(t *testing.T) {

	t.Run("test1", func(t *testing.T) {
		item := NewFakeItem(message.AJ054580)
		spanItems := item.SpanItems()
		size := len(spanItems)
		fmt.Println("size: ", size)
		fmt.Println("req: ", spanItems[0])
	})

	t.Run("test2", func(t *testing.T) {
		item1 := NewFakeItem(message.AJ054580)
		item2 := NewFakeItem(message.PG004040)
		fmt.Println(item1.String())
		item1.AddSpan(item2)
		items := item1.SpanItems()
		it := util.NewIterator(items)

		for it.HasNext() {
			i, next := it.Next()
			fmt.Println(i, " : ", next.String())
		}
	})

	t.Run("PushItem", func(t *testing.T) {
		item1 := NewFakeItem(message.AJ054580)
		item2 := NewFakeItem(message.PG004040)
		item3 := NewFakeItem(message.OU089012)
		outer := []*Item{item1, item1, item2, item2}
		items := pushItems([]*Item{item3}, outer)
		for _, i := range items {
			fmt.Println(i.String())
		}
	})

	// https://github.com/go-faker/faker
	// https://kgw7401.tistory.com/110
	t.Run("fake-user", func(t *testing.T) {
		user := &User{}
		err := faker.FakeData(user)
		if err == nil {
			fmt.Printf("Fake user : %v\n", user)
		}
	})
}
