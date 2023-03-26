package tools

import (
	"io/ioutil"
	"testing"
	"webrtc-exmple/utils"
)

func TestH264ToJpg(t *testing.T) {
	// 将H.264帧转换为JPEG
	h264Bytes, err := ioutil.ReadFile("/Users/lixiaolong/tmp/20230326/foreman_cif.h264")

	// h264Bytes := []byte{0x00, 0x00, 0x01, 0x67, 0x42, 0x80, 0x0d, 0xda, 0x38, 0x9b, 0x60}
	jpegBytes, err := utils.H264ToJpeg(h264Bytes)
	if err != nil {
		panic(err)
	}

	// 将JPEG数据写入文件
	if err := ioutil.WriteFile("/Users/lixiaolong/tmp/20230326/output.jpg", jpegBytes, 0644); err != nil {
		panic(err)
	}
}

func Test(t *testing.T) {
	//h264Bytes, err := ioutil.ReadFile("/Users/lixiaolong/tmp/20230326/foreman_cif.h264")
	//
	//d, err := decoder.New(decoder.PixelFormatBGR, decoder.H264)
	//if err != nil {
	//	panic(err)
	//}
	//
	//frames, err := d.Decode(h264Bytes)
	//if err != nil {
	//	t.Error(err)
	//}
	//if len(frames) == 0 {
	//	t.Log("no frames")
	//} else {
	//	frameCounter := 0
	//	for _, frame := range frames {
	//		img := frame.ToRGB()
	//		f, err := os.Create(fmt.Sprintf("/Users/lixiaolong/tmp/20230326/output_frame_%d.jpg", frameCounter))
	//		frameCounter++
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//		f.Close()
	//	}
	//	t.Logf("found %d frames", len(frames))
	//}
}
