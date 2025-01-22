package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		cache := NewCache(3)

		// Вытеснение
		cache.Set("aaa", 100)
		cache.Set("bbb", 200)
		cache.Set("ccc", 300)
		cache.Set("ddd", 400)

		// Этот элемент должен был вытеснен из кэша
		val, ok := cache.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		// А эти должны остаться
		val, ok = cache.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)
		val, ok = cache.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)
		val, ok = cache.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)

		// Очистка кэша
		cache.Clear()
		val, ok = cache.Get("ddd")
		require.False(t, ok)
		require.Nil(t, val)

		// Вытеснение редко используемого элемента
		cache.Set("aaa", 100)
		cache.Set("bbb", 200)
		cache.Set("ccc", 300)

		cache.Get("aaa")
		cache.Set("aaa", 400)
		cache.Get("aaa")

		cache.Set("bbb", 400)
		cache.Get("bbb")
		cache.Get("bbb")

		cache.Set("ddd", 400)
		// "ccc" редко используемый ключ элемента, он должен быть вытеснен
		val, ok = cache.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = cache.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 400, val)
		val, ok = cache.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 400, val)
		val, ok = cache.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()

	c.Clear()
	val, ok := c.Get("999999")
	require.False(t, ok)
	require.Nil(t, val)
}
