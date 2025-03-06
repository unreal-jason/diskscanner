package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"diskscanner/internal/types"
	"diskscanner/internal/utils"
)

type Scanner struct {
	files     []types.FileInfo
	mutex     sync.Mutex
	totalSize int64
	fileCount int
	workers   int
}

func NewScanner(workers int) *Scanner {
	return &Scanner{
		files:   make([]types.FileInfo, 0),
		workers: workers,
	}
}

func (s *Scanner) ScanDirectory(root string) error {
	// 首先验证目录是否存在
	if _, err := os.Stat(root); err != nil {
		return fmt.Errorf("无法访问目录 %s: %v", root, err)
	}

	var wg sync.WaitGroup
	paths := make(chan string)
	done := make(chan bool)

	// 启动目录处理协程
	go s.processDirs(root, paths)

	// 启动工作协程
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go s.processFiles(paths, &wg)
	}

	// 启动进度显示
	go s.showProgress(done)

	// 等待所有工作完成
	wg.Wait()
	done <- true
	fmt.Println()

	return nil
}

// 处理目录
func (s *Scanner) processDirs(root string, paths chan<- string) {
	defer close(paths)
	dirs := []string{root}
	for len(dirs) > 0 {
		currentDir := dirs[len(dirs)-1]
		dirs = dirs[:len(dirs)-1]
		paths <- currentDir

		entries, err := os.ReadDir(currentDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				fullPath := filepath.Join(currentDir, entry.Name())
				dirs = append(dirs, fullPath)
			}
		}
	}
}

// 处理文件
func (s *Scanner) processFiles(paths <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			fullPath := filepath.Join(path, entry.Name())
			info, err := entry.Info()
			if err != nil {
				continue
			}

			s.mutex.Lock()
			s.files = append(s.files, types.FileInfo{
				Path: fullPath,
				Size: info.Size(),
			})
			s.totalSize += info.Size()
			s.fileCount++
			s.mutex.Unlock()
		}
	}
}

// 显示进度
func (s *Scanner) showProgress(done chan bool) {
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	var lastCount int
	var lastSize int64
	lastUpdate := time.Now()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			duration := now.Sub(lastUpdate)

			s.mutex.Lock()
			currentCount := s.fileCount
			currentSize := s.totalSize
			countDiff := currentCount - lastCount
			sizeDiff := currentSize - lastSize

			fileSpeed := float64(countDiff) / duration.Seconds()
			sizeSpeed := float64(sizeDiff) / duration.Seconds()

			fmt.Printf("\r已扫描: %d 个文件, 总大小: %s (%.2f MB/s), 速度: %.0f 文件/秒",
				currentCount,
				utils.FormatSize(currentSize),
				sizeSpeed/(1024*1024),
				fileSpeed)

			lastCount = currentCount
			lastSize = currentSize
			lastUpdate = now
			s.mutex.Unlock()
		case <-done:
			return
		}
	}
}

// 获取扫描结果
func (s *Scanner) GetResults() []types.FileInfo {
	return s.files
}
