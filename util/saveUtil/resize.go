package saveUtil

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
)

func Compress(buf []byte) ([]byte, error) {
	//var width uint = 200
	//var height uint = 200

	//文件压缩
	decodeBuf, layout, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	// 修改图片的大小
	//set := resize.Resize(width, height, decodeBuf, resize.Lanczos3)
	NewBuf := bytes.Buffer{}
	switch layout {
	case "png":
		err = png.Encode(&NewBuf, decodeBuf)
	case "jpeg", "jpg":
		err = jpeg.Encode(&NewBuf, decodeBuf, &jpeg.Options{Quality: 80})
	default:
		return nil, errors.New("该图片格式不支持压缩")
	}
	if err != nil {
		return nil, err
	}
	if NewBuf.Len() < len(buf) {
		buf = NewBuf.Bytes()
	}
	return buf, nil
}
