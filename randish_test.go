package randish_test

import (
	"sync"
	"testing"

	"github.com/TechMDW/randish"
)

func TestSeed(t *testing.T) {
	const numIterations = 100000

	duplicates := sync.Map{}

	for i := 0; i < numIterations; i++ {
		seed := randish.Seed()
		_, loaded := duplicates.LoadOrStore(seed, true)
		if loaded {
			t.Fatalf("Found duplicate number %d after %d iterations", seed, i)
		}

	}
}

func TestSeedParallel(t *testing.T) {
	const numIterations = 2000000
	const numWorkers = 50

	duplicates := make(map[int64]bool)

	var mu sync.Mutex

	lim := make(chan struct{}, numWorkers)
	wg := sync.WaitGroup{}
	for i := 0; i < numIterations; i++ {
		lim <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-lim
				wg.Done()
			}()
			seed := randish.Seed()

			mu.Lock()
			if duplicates[seed] {
				t.Errorf("Found duplicate number %d after %d iterations", seed, i)
			}

			duplicates[seed] = true
			mu.Unlock()
		}(i)
	}

	wg.Wait()
}

func TestDistribution(t *testing.T) {
	const numIterations = 100000
	const numValues = 2
	const tolerance = 0.01

	counts := make([]int, numValues)
	r := randish.Rand()

	for i := 0; i < numIterations; i++ {
		value := r.Intn(numValues)
		counts[value]++
	}

	for i, count := range counts {
		percent := float64(count) / float64(numIterations)
		if percent < 1.0/numValues-tolerance || percent > 1.0/numValues+tolerance {
			t.Errorf("Value %d is outside the expected range: got %f, expected about %f", i, percent, 1.0/numValues)
		}
	}
}

func BenchmarkRandish(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randish := randish.Rand()
		randish.Int()
	}
}
