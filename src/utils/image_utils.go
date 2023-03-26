package utils

import (
	"io/ioutil"
	"os/exec"
)

// H264ToJpeg 将H.264帧转换为JPEG格式
func H264ToJpeg(h264Bytes []byte) ([]byte, error) {
	// 创建FFmpeg进程
	cmd := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-f", "h264", "-i", "pipe:", "-f", "mjpeg", "-")
	// 将H.264数据传递给FFmpeg进程的标准输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	// 执行FFmpeg进程并将输出捕获到变量中
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if _, err := stdin.Write(h264Bytes); err != nil {
		return nil, err
	}
	if err := stdin.Close(); err != nil {
		return nil, err
	}

	jpegBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return jpegBytes, nil
}

func main() {
	// 将H.264帧转换为JPEG
	h264Bytes := []byte{0x00, 0x00, 0x01, 0x67, 0x42, 0x80, 0x0d, 0xda, 0x38, 0x9b, 0x60}
	jpegBytes, err := H264ToJpeg(h264Bytes)
	if err != nil {
		panic(err)
	}

	// 将JPEG数据写入文件
	if err := ioutil.WriteFile("/Users/lixiaolong/tmp/20230326/output.jpg", jpegBytes, 0644); err != nil {
		panic(err)
	}
}
