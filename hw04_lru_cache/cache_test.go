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

	t.Run("remove item if capacity exceeded", func(t *testing.T) {
		c := NewCache(2)

		c.Set("aaa", 100)
		_, ok := c.Get("aaa")
		require.True(t, ok)

		c.Set("bbb", 200)
		c.Set("ccc", 300)

		//capacity exceeded, "aaa" should have been removed
		deleted, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, deleted)
	})

	t.Run("remove item if not least recently used", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		c.Get("aaa")
		c.Get("bbb")

		//capacity exceeded, "ccc" (not used element) should have been removed
		c.Set("ddd", 400)

		notUsedElem, ok := c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, notUsedElem)
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

	//Anton's tests
	const (
		goroutines = 100
		iterations = 200
		capacity   = 50
	)
	t.Run("concurrent set and get", func(t *testing.T) {
		c := NewCache(capacity)
		var wg sync.WaitGroup

		// Writers
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					c.Set(Key(strconv.Itoa(i*iterations+j)), i*j)
				}
			}(i)
		}

		// Concurrent readers — стартуют одновременно с writers
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					c.Get(Key(strconv.Itoa(i*iterations + j)))
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("concurrent set and clear", func(t *testing.T) {
		c := NewCache(capacity)
		var wg sync.WaitGroup

		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				c.Set(Key(strconv.Itoa(i)), i)
			}(i)
		}

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c.Clear()
			}()
		}

		wg.Wait()
	})
}
