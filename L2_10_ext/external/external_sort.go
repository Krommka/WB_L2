package external

import (
	"bufio"
	"fmt"
	parsingflags "main/parsingFlags"
	"main/sorting"
	"os"
	"sort"
)

const (
	ChunkSize    = 1 * 1024 * 1024 * 1024 // 1GB
	MaxOpenFiles = 32
)

type Chunk struct {
	file   *os.File
	reader *bufio.Reader
	line   string
	index  int
}

type ChunkHeap []*Chunk

func (h ChunkHeap) Len() int           { return len(h) }
func (h ChunkHeap) Less(i, j int) bool { return h[i].line < h[j].line }
func (h ChunkHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ChunkHeap) Push(x interface{}) {
	*h = append(*h, x.(*Chunk))
}

func (h *ChunkHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func ExternalSort(lines []string, config *parsingflags.Config, outputWriter *os.File) error {
	// Если данные помещаются в память, используем обычную сортировку
	if estimateMemoryUsage(lines) < ChunkSize {
		return sortInMemory(lines, config, outputWriter)
	}

	// Разбиваем на чанки и сортируем их
	chunkFiles, err := createSortedChunks(lines, config)
	if err != nil {
		return err
	}
	defer cleanupChunkFiles(chunkFiles)

	// Сливаем чанки
	return mergeChunks(chunkFiles, config, outputWriter)
}

func estimateMemoryUsage(lines []string) int64 {
	var total int64
	for _, line := range lines {
		total += int64(len(line))
	}
	return total
}

func sortInMemory(lines []string, config *parsingflags.Config, outputWriter *os.File) error {
	sorted := make([]string, len(lines))
	copy(sorted, lines)

	sort.Slice(sorted, func(i, j int) bool {
		return sorting.CompareStrings(sorted[i], sorted[j], config)
	})

	if config.Unique {
		sorted = sorting.RemoveDuplicates(sorted)
	}

	for _, line := range sorted {
		if _, err := fmt.Fprintln(outputWriter, line); err != nil {
			return err
		}
	}

	return nil
}

func createSortedChunks(lines []string, config *parsingflags.Config) ([]string, error) {
	var chunkFiles []string
	var currentChunk []string
	var currentSize int64

	flushChunk := func() error {
		if len(currentChunk) == 0 {
			return nil
		}

		// Сортируем чанк
		sort.Slice(currentChunk, func(i, j int) bool {
			return compareLines(currentChunk[i], currentChunk[j], config)
		})

		if config.Unique {
			currentChunk = removeDuplicates(currentChunk)
		}

		// Сохраняем во временный файл
		chunkFile, err := saveChunkToFile(currentChunk)
		if err != nil {
			return err
		}

		chunkFiles = append(chunkFiles, chunkFile)
		currentChunk = nil
		currentSize = 0

		return nil
	}

	for _, line := range lines {
		lineSize := int64(len(line))

		if currentSize+lineSize > ChunkSize && len(currentChunk) > 0 {
			if err := flushChunk(); err != nil {
				return nil, err
			}
		}

		currentChunk = append(currentChunk, line)
		currentSize += lineSize
	}

	// Флашим последний чанк
	if err := flushChunk(); err != nil {
		return nil, err
	}

	return chunkFiles, nil
}
