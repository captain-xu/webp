package main

import (
	"os"
	"os/exec"
	"bufio"
	"io"
	"path/filepath"
	"fmt"
	"strings"
)

var fileRealPath string

func init() {
	fileRealPath = "/Users/aihuishou/Downloads/images"
	//fileRealPath = getCurrentDirectory()
}

func main() {
	getFileList(fileRealPath)
}

func getFileList(filep string) {
	err := filepath.Walk(filep, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		png2webp(f.Name(), path)
		//deleteWebp(f.Name(),path)
		return nil
	})
	if err == nil {
		fmt.Println("转换完毕.")
	}
}

// 通过cwebp命令批量替换png图片
func png2webp(name string, path string) {
	var isLauncher = strings.HasPrefix(name, "ic_launcher")

	var isPng = strings.HasSuffix(name, ".png")

	var isJpg = strings.HasSuffix(name, ".jpg")

	var is9png = strings.HasSuffix(name, ".9.png")

	// 对png和pjg同时进行处理，排除掉应用图标及.9图
	if isPng && !isLauncher && !is9png || isJpg {

		var out string
		if isJpg {
			out = strings.TrimSuffix(path, ".jpg") + ".webp"
		}
		if isPng {
			out = strings.TrimSuffix(path, ".png") + ".webp"
		}

		cmdErr := execCommand("cwebp", "-q", "80", path, "-o", out)

		if cmdErr {
			println(path, " 替换为 =======> ", out)
			del := os.Remove(path)
			if del == nil {
				println(path, " 删除成功!")
			}
		}
	}
}

// 删除当前目录中的webp文件
func deleteWebp(name string, path string) {
	var isweb = strings.HasSuffix(name, ".webp")

	if isweb {
		del := os.Remove(path)
		if del == nil {
			println(path, " 删除成功!")
		}
	}
}

// 执行cwebp转换命令
func execCommand(commandName string, params ...string) bool {
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	cmd.Wait()
	return true
}

