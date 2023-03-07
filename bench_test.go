package main

import (
	"bufio"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	wenyxu "github.com/WenyXu/sync-adaptive-radix-tree"
	acomagu "github.com/acomagu/trie"
	arriqaaq "github.com/arriqaaq/art"
	baevik "github.com/beevik/prefixtree"
	bobotu "github.com/bobotu/opt-art"
	derekparker "github.com/derekparker/trie"
	dghubble "github.com/dghubble/trie"
	gammazero "github.com/gammazero/radixtree"
	gbrlsnchs "github.com/gbrlsnchs/radix"
	cedar "github.com/go-ego/cedar"
	doublearray "github.com/kampersanda/doublearray-go"
	kellydunn "github.com/kellydunn/go-art"
	ternary "github.com/manveru/trie"
	"github.com/openacid/slim/encode"
	slim "github.com/openacid/slim/trie"
	plar "github.com/plar/go-adaptive-radix-tree"
	recoilme "github.com/recoilme/art"
	ryuanerin "github.com/ryuanerin/ptrie"
	superfell "github.com/superfell/art"
	patrica "github.com/tchap/go-patricia/v2/patricia"
	wzshiming "github.com/wzshiming/trie"
)

const (
	prefixesPath       = "data/prefixes.csv"
	prefixesSortedPath = "data/prefixes-sorted.csv"
)

// Common

// In test format for easy run
func TestTipsLeavesToPrefixes(t *testing.T) {
	inputLeavesPath := "data/tips-leaves.csv"
	outputPrefixesPath := prefixesPath
	outputPrefixesSortedPath := prefixesSortedPath

	prefixes := make([]string, 0)

	file, err := os.Open(inputLeavesPath)
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	// skip header
	scanner.Scan()

	for scanner.Scan() {
		text := scanner.Text()
		columns := strings.Split(text, "|")

		if len(columns) < 2 {
			t.Fatal()
		}

		prefix := columns[0]
		prefixes = append(prefixes, prefix)
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	// shuffle prefixes
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(prefixes), func(i, j int) {
		prefixes[i], prefixes[j] = prefixes[j], prefixes[i]
	})

	// write prefixes
	file, err = os.Create(outputPrefixesPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, prefix := range prefixes {
		_, err = file.WriteString(prefix + "\n")
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	sort.Strings(prefixes)

	file, err = os.Create(outputPrefixesSortedPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, prefix := range prefixes {
		_, err = file.WriteString(prefix + "\n")
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func PrefixesBytes(count int, path string) ([][]byte, error) {
	input := make([][]byte, 0)
	result := make([][]byte, 0)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input = append(input, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	// if -1 - returns all
	if count == -1 {
		count = len(input)
	}

	for i := 0; i < count; i++ {
		result = append(result, input[i%len(input)])
	}

	// shuffle

	return result, nil
}

func PrefixesStrings(count int, path string) ([]string, error) {
	resultBytes, err := PrefixesBytes(count, path)
	if err != nil {
		return nil, err
	}

	resultStrs := make([]string, 0, len(resultBytes))
	for _, bytes := range resultBytes {
		resultStrs = append(resultStrs, string(bytes))
	}

	return resultStrs, nil
}

// Measure sort

func BenchmarkSortInput(b *testing.B) {
	resultBytes, err := PrefixesBytes(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	resultStrs := make([]string, 0, len(resultBytes))
	for _, bytes := range resultBytes {
		resultStrs = append(resultStrs, string(bytes))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sort.Strings(resultStrs)
	}
}

// Reference map

func BenchmarkMapInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := make(map[string]int, 0)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree[prefixes[i]] = 0
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := make(map[string]int, 0)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree[prefixes[i]] = 0
		}
	})
}

func BenchmarkMapInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := make(map[string]int)
			for j := 0; j < len(prefixes); j++ {
				tree[prefixes[j]] = 0
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkMapInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := make(map[string]int)
			for j := 0; j < len(prefixes); j++ {
				tree[prefixes[j]] = 0
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkMapInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkMapExistsSearch(b *testing.B) {
	prefixes, err := PrefixesStrings(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := make(map[string]int)
	for i := 0; i < len(prefixes); i++ {
		tree[prefixes[i]] = 0
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = tree[prefixes[i%len(prefixes)]]
	}
}

// Cedar

func BenchmarkCedarInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := cedar.New()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Insert(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := cedar.New()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Insert(prefixes[i], 0)
		}
	})
}

func BenchmarkCedarInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := cedar.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkCedarInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := cedar.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkCedarInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkCedarExistsSearch(b *testing.B) {
	prefixes, err := PrefixesBytes(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := cedar.New()
	for i := 0; i < len(prefixes); i++ {
		tree.Insert(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(prefixes[i%len(prefixes)])
	}
}

// Slim

func BenchmarkSlimInsertAll(b *testing.B) {
	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		ints := make([]int, len(prefixes))

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := slim.NewSlimTrie(encode.Dummy{}, prefixes, ints, slim.Opt{
				Complete: slim.Bool(true),
			})
			if err != nil {
				b.Fatal(err)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkSlimInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkSlimExistsSearch(b *testing.B) {
	prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
	if err != nil {
		b.Fatal(err)
	}

	ints := make([]int, len(prefixes))

	tree, err := slim.NewSlimTrie(encode.Dummy{}, prefixes, ints, slim.Opt{
		Complete: slim.Bool(true),
	})
	if err != nil {
		b.Fatal(err)
	}

	shufflePrefixes := make([]string, len(prefixes))
	copy(shufflePrefixes, prefixes)
	rand.Seed(time.Now().Unix())

	rand.Shuffle(len(shufflePrefixes), func(i, j int) {
		shufflePrefixes[i], shufflePrefixes[j] = shufflePrefixes[j], shufflePrefixes[i]
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(shufflePrefixes[i%len(shufflePrefixes)])
	}
}

// Double array

func BenchmarkDoubleArrayInsertAll(b *testing.B) {
	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		ints := make([]int, len(prefixes))

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := doublearray.Build(prefixes, ints)
			if err != nil {
				b.Fatal(err)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkDoubleArrayInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkDoubleArrayExistsSearch(b *testing.B) {
	prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
	if err != nil {
		b.Fatal(err)
	}

	ints := make([]int, len(prefixes))

	tree, err := doublearray.Build(prefixes, ints)
	if err != nil {
		b.Fatal(err)
	}

	shufflePrefixes := make([]string, len(prefixes))
	copy(shufflePrefixes, prefixes)
	rand.Seed(time.Now().Unix())

	rand.Shuffle(len(shufflePrefixes), func(i, j int) {
		shufflePrefixes[i], shufflePrefixes[j] = shufflePrefixes[j], shufflePrefixes[i]
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Lookup(shufflePrefixes[i%len(shufflePrefixes)])
	}
}

// ternary tree

func BenchmarkTernaryTreeInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := ternary.Trie{}
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTernaryTreeInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := ternary.Trie{}
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTernaryTreeInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

// ARTs

// Plar

func BenchmarkArtCurrentPlarInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := plar.New()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Insert(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := plar.New()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Insert(prefixes[i], 0)
		}
	})
}

func BenchmarkArtCurrentPlarInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := plar.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtCurrentPlarInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := plar.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtCurrentPlarInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkArtCurrentPlarExistsSearch(b *testing.B) {
	prefixes, err := PrefixesBytes(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := plar.New()
	for i := 0; i < len(prefixes); i++ {
		tree.Insert(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Search(prefixes[i%len(prefixes)])
	}
}

// gammazero

func BenchmarkArtGammazeroInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := gammazero.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtGammazeroInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := gammazero.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtGammazeroInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkArtGammazeroInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := gammazero.New()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Put(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := gammazero.New()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Put(prefixes[i], 0)
		}
	})
}

func BenchmarkArtGammazeroExistsSearch(b *testing.B) {
	prefixes, err := PrefixesStrings(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := gammazero.New()
	for i := 0; i < len(prefixes); i++ {
		tree.Put(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(prefixes[i%len(prefixes)])
	}
}

// recoilme

func BenchmarkArtRecoilmeInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := recoilme.New()
			val := []byte("val")
			for j := 0; j < len(prefixes); j++ {
				tree.Set(prefixes[j], val)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtRecoilmeInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := recoilme.New()
			val := []byte("val")
			for j := 0; j < len(prefixes); j++ {
				tree.Set(prefixes[j], val)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtRecoilmeInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

// superfell

func BenchmarkArtSuperfellInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := superfell.Tree[int]{}
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtSuperfellInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := superfell.Tree[int]{}
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtSuperfellInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkArtSuperfellInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := superfell.Tree[int]{}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Put(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := superfell.Tree[int]{}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Put(prefixes[i], 0)
		}
	})
}

func BenchmarkArtSuperfellExistsSearch(b *testing.B) {
	prefixes, err := PrefixesBytes(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := superfell.Tree[int]{}
	for i := 0; i < len(prefixes); i++ {
		tree.Put(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(prefixes[i%len(prefixes)])
	}
}

// arriqaaq

func BenchmarkArtArriqaaqInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := arriqaaq.NewTree()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtArriqaaqInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := arriqaaq.NewTree()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtArriqaaqInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

// kellydunn

func BenchmarkArtKellydunnInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := kellydunn.NewArtTree()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtKellydunnInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := kellydunn.NewArtTree()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtKellydunnInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

// Bobotu

func BenchmarkArtBobotuInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := bobotu.NewART()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtBobotuInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := bobotu.NewART()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtBobotuInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkArtBobotuInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := bobotu.NewART()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Put(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := bobotu.NewART()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Put(prefixes[i], 0)
		}
	})
}

func BenchmarkArtBobotuExistsSearch(b *testing.B) {
	prefixes, err := PrefixesBytes(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := bobotu.NewART()
	for i := 0; i < len(prefixes); i++ {
		tree.Put(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(prefixes[i%len(prefixes)])
	}
}

// Wenyxu

func BenchmarkArtWenyxuTreeInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := wenyxu.Tree[int]{}
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtWenyxuTreeInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := wenyxu.Tree[int]{}
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkArtWenyxuTreeInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkArtWenyxuInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := wenyxu.Tree[int]{}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Insert(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := wenyxu.Tree[int]{}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Insert(prefixes[i], 0)
		}
	})
}

func BenchmarkArtWenyxuExistsSearch(b *testing.B) {
	prefixes, err := PrefixesBytes(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := wenyxu.Tree[int]{}
	for i := 0; i < len(prefixes); i++ {
		tree.Insert(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Search(prefixes[i%len(prefixes)])
	}
}

// Radixes

// patrica

func BenchmarkRadixPatricaTreeInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := patrica.NewTrie()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkRadixPatricaTreeInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := patrica.NewTrie()
			for j := 0; j < len(prefixes); j++ {
				tree.Insert(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkRadixPatricaTreeInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

// wzshiming
func BenchmarkRadixWzshimingInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := wzshiming.NewTrie[int]()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkRadixWzshimingInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := wzshiming.NewTrie[int]()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkRadixWzshimingInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkRadixGbrlsnchsInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := gbrlsnchs.New(gbrlsnchs.Tnocolor)
			for j := 0; j < len(prefixes); j++ {
				tree.Add(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkRadixGbrlsnchsInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := gbrlsnchs.New(gbrlsnchs.Tnocolor)
			for j := 0; j < len(prefixes); j++ {
				tree.Add(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkRadixGbrlsnchsInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkRadixGbrlsnchsInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := gbrlsnchs.New(gbrlsnchs.Tnocolor)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Add(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := gbrlsnchs.New(gbrlsnchs.Tnocolor)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Add(prefixes[i], 0)
		}
	})
}

func BenchmarkRadixGbrlsnchsExistsSearch(b *testing.B) {
	prefixes, err := PrefixesStrings(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := gbrlsnchs.New(gbrlsnchs.Tnocolor)
	for i := 0; i < len(prefixes); i++ {
		tree.Add(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(prefixes[i%len(prefixes)])
	}
}

// Tries

// ryuanerin

func BenchmarkTrieRyuanerinInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := ryuanerin.New[int]()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieRyuanerinInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := ryuanerin.New[int]()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieRyuanerinInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkTrieRyuanerinInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := ryuanerin.New[int]()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Put(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := ryuanerin.New[int]()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Put(prefixes[i], 0)
		}
	})
}

func BenchmarkTrieRyuanerinExistsSearch(b *testing.B) {
	prefixes, err := PrefixesBytes(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := ryuanerin.New[int]()
	for i := 0; i < len(prefixes); i++ {
		tree.Put(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(prefixes[i%len(prefixes)])
	}
}

// dghubble
func BenchmarkTrieDghubbleInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := dghubble.NewRuneTrie()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieDghubbleInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := dghubble.NewRuneTrie()
			for j := 0; j < len(prefixes); j++ {
				tree.Put(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieDghubbleInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

// acomagu

func BenchmarkTrieAcomaguInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		ints := make([]interface{}, len(prefixes))

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			acomagu.New(prefixes, ints)
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieAcomaguInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesBytes(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		ints := make([]interface{}, len(prefixes))

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			acomagu.New(prefixes, ints)
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieAcomaguInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

// derekparker

func BenchmarkTrieDerekparkerInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := derekparker.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Add(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieDerekparkerInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := derekparker.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Add(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieDerekparkerInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

// beevik

func BenchmarkTrieBeevikInsertAll(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := baevik.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Add(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieBeevikInsertAll/unsorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(-1, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := baevik.New()
			for j := 0; j < len(prefixes); j++ {
				tree.Add(prefixes[j], 0)
			}
		}

		runtime.ReadMemStats(&m2)
		b.Logf("BenchmarkTrieBeevikInsertAll/sorted inuse: %f MB\n", float64(m2.HeapInuse-m1.HeapInuse)/(1024*1024))
	})
}

func BenchmarkTrieBeevikInsert(b *testing.B) {
	b.Run("unsorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(b.N, prefixesPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := baevik.New()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Add(prefixes[i], 0)
		}
	})

	b.Run("sorted", func(b *testing.B) {
		prefixes, err := PrefixesStrings(b.N, prefixesSortedPath)
		if err != nil {
			b.Fatal(err)
		}

		tree := baevik.New()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Add(prefixes[i], 0)
		}
	})
}

func BenchmarkTrieBeevikExistsSearch(b *testing.B) {
	prefixes, err := PrefixesStrings(-1, prefixesPath)
	if err != nil {
		b.Fatal(err)
	}

	tree := baevik.New()
	for i := 0; i < len(prefixes); i++ {
		tree.Add(prefixes[i], 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Find(prefixes[i%len(prefixes)])
	}
}
