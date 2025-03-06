package main

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"diskscanner/internal/scanner"
	"diskscanner/internal/utils"
)

func TestMain(t *testing.T) {
	rootPath := "D:\\"
	diskScanner := scanner.NewScanner(4)

	fmt.Println("开始扫描...")
	startTime := time.Now()

	err := diskScanner.ScanDirectory(rootPath)
	if err != nil {
		fmt.Printf("扫描错误: %v\n", err)
		os.Exit(1)
	}

	files := diskScanner.GetResults()

	// 按文件大小排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	// 输出前10大文件
	fmt.Println("\n最大的10个文件:")
	fmt.Println("大小\t\t\t路径")
	fmt.Println("-------------------------------------------------------------")

	count := 10
	if len(files) < 10 {
		count = len(files)
	}

	for i := 0; i < count; i++ {
		fmt.Printf("%-15s\t%s\n", utils.FormatSize(files[i].Size), files[i].Path)
	}

	// 保存结果到CSV
	csvFile := "scan_results.csv"
	if err := utils.SaveToCSV(files, csvFile); err != nil {
		fmt.Printf("保存CSV文件失败: %v\n", err)
	} else {
		fmt.Printf("\n完整结果已保存到: %s\n", csvFile)
	}

	fmt.Printf("\n扫描完成! 耗时: %.2f 秒\n", time.Since(startTime).Seconds())
}
