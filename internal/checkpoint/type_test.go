package checkpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequential(t *testing.T) {
	c := New()
	assert.Equal(t, nil, c.Highest())
	assert.EqualValues(t, 0, c.Pending())

	res1 := c.Track(1, 1)
	res2 := c.Track(2, 1)
	res3 := c.Track(3, 1)
	assert.Equal(t, nil, c.Highest())
	assert.EqualValues(t, 3, c.Pending())

	v := res1()
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, c.Highest())
	assert.EqualValues(t, 2, c.Pending())

	v = res2()
	assert.Equal(t, 2, v)
	assert.Equal(t, 2, c.Highest())
	assert.EqualValues(t, 1, c.Pending())

	v = res3()
	assert.Equal(t, 3, v)
	assert.Equal(t, 3, c.Highest())
	assert.EqualValues(t, 0, c.Pending())

	res4 := c.Track(4, 1)
	assert.EqualValues(t, 1, c.Pending())

	v = res4()
	assert.Equal(t, 4, v)
	assert.Equal(t, 4, c.Highest())
	assert.EqualValues(t, 0, c.Pending())
}

func TestOutOfSync(t *testing.T) {
	c := New()
	assert.Equal(t, nil, c.Highest())

	res1 := c.Track(1, 1)
	res2 := c.Track(2, 1)
	res3 := c.Track(3, 1)
	res4 := c.Track(4, 1)
	assert.Equal(t, nil, c.Highest())

	v := res2()
	assert.Equal(t, nil, v)
	assert.Equal(t, nil, c.Highest())

	v = res1()
	assert.Equal(t, 2, v)
	assert.Equal(t, 2, c.Highest())

	v = res3()
	assert.Equal(t, 3, v)
	assert.Equal(t, 3, c.Highest())

	v = res4()
	assert.Equal(t, 4, v)
	assert.Equal(t, 4, c.Highest())
}

func TestSequentialLarge(t *testing.T) {
	c := New()
	var resolves []func() any

	for i := 0; i < 1000; i++ {
		resolves = append(resolves, c.Track(i, 1))
	}
	for i := 0; i < 1000; i++ {
		v := resolves[i]()
		assert.Equal(t, i, v)
		assert.Equal(t, i, c.Highest())
	}
}

func TestSequentialChunks(t *testing.T) {
	c := New()
	chunkSize := 100
	for i := 0; i < 10; i++ {
		var resolves []func() any

		for j := 0; j < chunkSize; j++ {
			offset := i*chunkSize + j
			resolves = append(resolves, c.Track(offset, 1))
		}

		for j := 0; j < chunkSize; j++ {
			offset := i*chunkSize + j
			v := resolves[j]()
			assert.Equal(t, offset, v)
			assert.Equal(t, offset, c.Highest())
		}
	}
}

func TestSequentialReverseLarge(t *testing.T) {
	c := New()
	var resolves []func() any

	for i := 0; i < 1000; i++ {
		resolves = append(resolves, c.Track(i, 1))
	}
	for i := 999; i >= 0; i-- {
		v := resolves[i]()
		if i == 0 {
			assert.Equal(t, 999, v)
			assert.Equal(t, 999, c.Highest())
		} else {
			assert.Equal(t, nil, v)
			assert.Equal(t, nil, c.Highest())
		}
	}
}

func TestSequentialRandomLarge(t *testing.T) {
	c := New()
	resolves := make([]func() any, 1000)
	indexes := map[int]struct{}{}
	for i := 0; i < 1000; i++ {
		resolves[i] = c.Track(i, 1)
		indexes[i] = struct{}{}
	}
	for i := range indexes {
		delete(indexes, i)
		v := resolves[i]()
		if len(indexes) == 0 {
			assert.Equal(t, 999, v)
			assert.Equal(t, 999, c.Highest())
		} else {
			assert.Equal(t, v, c.Highest())
			if v != nil {
				for k := range indexes {
					assert.False(t, k < v.(int))
					assert.False(t, k < c.Highest().(int))
				}
			}
		}
	}
}

func BenchmarkChunked100(b *testing.B) {
	b.ReportAllocs()
	c := New()
	chunkSize := 100
	N := b.N / chunkSize
	for i := 0; i < N; i++ {
		resolves := make([]func() any, chunkSize)

		for j := 0; j < chunkSize; j++ {
			offset := i*chunkSize + j
			resolves[j] = c.Track(offset, 1)
		}

		for j := 0; j < chunkSize; j++ {
			offset := i*chunkSize + j
			v := resolves[j]()
			if offset != v {
				b.Errorf("Wrong value: %v != %v", offset, v)
			}
		}
	}
}

func BenchmarkChunkedReverse100(b *testing.B) {
	b.ReportAllocs()
	c := New()
	chunkSize := 100
	N := b.N / chunkSize
	for i := 0; i < N; i++ {
		resolves := make([]func() any, chunkSize)

		for j := 0; j < chunkSize; j++ {
			offset := i*chunkSize + j
			resolves[j] = c.Track(offset, 1)
		}

		for j := chunkSize - 1; j >= 0; j-- {
			v := resolves[j]()
			var exp any
			if i > 0 {
				exp = (i * chunkSize) - 1
			}
			if j == 0 {
				exp = ((i + 1) * chunkSize) - 1
			}
			if exp == 0 {
				exp = nil
			}
			if exp != v {
				b.Errorf("Wrong value: %v != %v", exp, v)
			}
		}
	}
}

func BenchmarkChunkedReverse1000(b *testing.B) {
	b.ReportAllocs()
	c := New()
	chunkSize := 1000
	N := b.N / chunkSize
	for i := 0; i < N; i++ {
		resolves := make([]func() any, chunkSize)

		for j := 0; j < chunkSize; j++ {
			offset := i*chunkSize + j
			resolves[j] = c.Track(offset, 1)
		}

		for j := chunkSize - 1; j >= 0; j-- {
			v := resolves[j]()
			var exp any
			if i > 0 {
				exp = (i * chunkSize) - 1
			}
			if j == 0 {
				exp = ((i + 1) * chunkSize) - 1
			}
			if exp == 0 {
				exp = nil
			}
			if exp != v {
				b.Errorf("Wrong value: %v != %v", exp, v)
			}
		}
	}
}

func BenchmarkSequential(b *testing.B) {
	b.ReportAllocs()
	c := New()
	resolves := make([]func() any, b.N)
	for i := 0; i < b.N; i++ {
		resolves[i] = c.Track(i, 1)
	}
	for i := 0; i < b.N; i++ {
		v := resolves[i]()
		if i != v {
			b.Errorf("Wrong value: %v != %v", i, v)
		}
	}
}
