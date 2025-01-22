package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("one item", func(tt *testing.T) {
		ll := NewList()

		// Добавление единственного элемента
		ll.PushBack(100)
		require.Equal(tt, 1, ll.Len())
		require.Equal(tt, 100, ll.Front().Value)
		require.Nil(tt, ll.Front().Next)
		require.Nil(tt, ll.Front().Prev)

		// Перемещение единственного элемента в начало списка (по факту ничего не должно измениться)
		ll.MoveToFront(ll.Front())
		require.Equal(tt, 1, ll.Len())
		require.Equal(tt, 100, ll.Front().Value)
		require.Nil(tt, ll.Front().Next)
		require.Nil(tt, ll.Front().Prev)

		// Удаление единственного элемента, список должен быть пустым
		ll.Remove(ll.Back())
		require.Equal(tt, 0, ll.Len())
		require.Nil(tt, ll.Front())
		require.Nil(tt, ll.Front())
	})

	t.Run("two item remove", func(tt *testing.T) {
		ll := NewList()

		// В списке два элемента, удаления первого элемента
		ll.PushBack(100)
		ll.PushBack(200)
		ll.Remove(ll.Front())
		require.Equal(tt, 1, ll.Len())
		require.Nil(tt, ll.Front().Next)
		require.Nil(tt, ll.Front().Prev)
		require.Equal(tt, ll.Front(), ll.Back())
		require.Equal(tt, 200, ll.Front().Value)

		// В списке два элемента, удаления последнего элемента
		ll.PushFront(100)
		ll.Remove(ll.Back())
		require.Equal(tt, 1, ll.Len())
		require.Nil(tt, ll.Front().Next)
		require.Nil(tt, ll.Front().Prev)
		require.Equal(tt, ll.Front(), ll.Back())
		require.Equal(tt, 100, ll.Front().Value)
	})

	t.Run("three item remove", func(tt *testing.T) {
		ll := NewList()

		// В списке три элемента, удаление первого элемента
		ll.PushBack(100)
		ll.PushBack(200)
		ll.PushBack(300)
		ll.Remove(ll.Front())
		require.Equal(tt, 2, ll.Len())
		require.Nil(tt, ll.Front().Prev)
		require.Nil(tt, ll.Back().Next)
		require.Equal(tt, ll.Front().Next, ll.Back())
		require.Equal(tt, ll.Back().Prev, ll.Front())
		require.Equal(tt, 200, ll.Front().Value)
		require.Equal(tt, 300, ll.Back().Value)

		// В списке три элемента, удаление центрального элемента
		ll.PushFront(100)
		ll.Remove(ll.Front().Next)
		require.Equal(tt, 2, ll.Len())
		require.Nil(tt, ll.Front().Prev)
		require.Nil(tt, ll.Back().Next)
		require.Equal(tt, ll.Front().Next, ll.Back())
		require.Equal(tt, ll.Back().Prev, ll.Front())
		require.Equal(tt, 100, ll.Front().Value)
		require.Equal(tt, 300, ll.Back().Value)

		// В списке три элемента, удаление последнего элемента
		ll.Remove(ll.Back())
		ll.PushBack(200)
		ll.PushBack(300)
		ll.Remove(ll.Back())
		require.Equal(tt, 2, ll.Len())
		require.Nil(tt, ll.Front().Prev)
		require.Nil(tt, ll.Back().Next)
		require.Equal(tt, ll.Front().Next, ll.Back())
		require.Equal(tt, ll.Back().Prev, ll.Front())
		require.Equal(tt, 100, ll.Front().Value)
		require.Equal(tt, 200, ll.Back().Value)
	})
}
