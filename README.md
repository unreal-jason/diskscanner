# DiskScanner

DiskScanner 是一个高性能的磁盘文件扫描工具，可以快速扫描指定目录下的所有文件，并按大小排序展示结果。

## 特性

- 多协程并发扫描，性能优异
- 实时显示扫描进度和速度
- 展示最大的文件列表
- 将扫描结果导出为 CSV 文件
- 支持显示文件大小的人性化格式（KB/MB/GB/TB）

## 编译

1. 克隆仓库
```bash
git clone https://github.com/unreal-jason/diskscanner.git
cd diskscanner
```
2. 直接运行
```bash
go run main.go [目录路径]
```
3. 编译&运行
```bash
go build
diskscanner.exe [目录路径]
```
4. 示例
```bash
diskscanner.exe D:\
开始扫描...

最大的10个文件:
大小                    路径
----------------------------------------
1.5 GB          D:\Games\game1.dat
800.5 MB        D:\Videos\movie.mp4
500.2 MB        D:\Documents\backup.zip

完整结果已保存到: scan_results.csv

扫描完成! 耗时: 2.35 秒
```