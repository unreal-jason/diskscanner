package utils

import (
	"diskscanner/internal/types"
	"encoding/csv"
	"fmt"
	"os"
)

// FormatSize 格式化文件大小
func FormatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// SaveToCSV 保存结果到CSV文件
func SaveToCSV(files []types.FileInfo, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"大小(字节)", "格式化大小", "文件路径"})

	for _, f := range files {
		writer.Write([]string{
			fmt.Sprintf("%d", f.Size),
			FormatSize(f.Size),
			f.Path,
		})
	}
	return nil
}
