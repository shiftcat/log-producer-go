package item

import (
	"example.com/message"
	"example.com/message/util"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func printItem(items []*Item) {
	//for _, v := range items {
	//	if v.Type == message.REQ {
	//		fmt.Println(v.IndentedString())
	//	} else {
	//		fmt.Println(v.ToRes().IndentedString())
	//	}
	//}

	iter := util.NewIterator(items)
	for iter.HasNext() {
		i, v := iter.Next()
		if v.Type == message.REQ {
			fmt.Println(i, " : ", v.IndentedString())
		} else {
			fmt.Println(i, " : ", v.ToRes().IndentedString())
		}
	}
}

func TestGenerate(t *testing.T) {

	assertion := assert.New(t)
	InitDataMap()

	t.Run("Case1", func(t *testing.T) {
		items := Case1()
		assertion.Condition(func() bool {
			return len(items) > 0
		})
		printItem(items)
	})

	t.Run("Case2", func(t *testing.T) {
		items := Case2()
		assertion.Condition(func() bool {
			return len(items) > 0
		})
		printItem(items)
	})

	t.Run("Case3", func(t *testing.T) {
		items := Case3()
		assertion.Condition(func() bool {
			return len(items) > 0
		})
		printItem(items)
	})

	t.Run("Case4", func(t *testing.T) {
		items := Case4()
		assertion.Condition(func() bool {
			return len(items) > 0
		})
		printItem(items)
	})

	t.Run("Case5", func(t *testing.T) {
		items := Case5()
		assertion.Condition(func() bool {
			return len(items) > 0
		})
		printItem(items)
	})

	t.Run("Fake1", func(t *testing.T) {
		items := Fake1()
		assertion.Condition(func() bool {
			return len(items) > 0
		})
		printItem(items)
	})
}
