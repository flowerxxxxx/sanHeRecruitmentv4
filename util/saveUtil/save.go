package saveUtil

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image"
	"io/ioutil"
	"log"
	"mime/multipart"
)

func SaveFile(file *multipart.FileHeader, fileAddr string) {
	fileContent, eo := file.Open()
	if eo != nil {
		log.Println("SaveFile file.Open failed,err:", eo)
		return
	}
	defer fileContent.Close()
	byteContainer, er := ioutil.ReadAll(fileContent)
	if er != nil {
		log.Println("SaveFile ioutil.ReadAll failed,err:", eo)
		return
	}
	//保存到新文件中
	ei := ioutil.WriteFile(fileAddr, byteContainer, 0644)
	if ei != nil {
		log.Println("SaveFile Compress failed,err:", eo)
		return
	}
	return
}

func SaveCompressFile(file *multipart.FileHeader, fileAddr string) error {
	fileContent, eo := file.Open()
	if eo != nil {
		log.Println("SaveFile file.Open failed,err:", eo)
		return eo
	}
	defer fileContent.Close()
	byteContainer, er := ioutil.ReadAll(fileContent)
	if er != nil {
		log.Println("SaveFile ioutil.ReadAll failed,err:", eo)
		return er
	}
	afterResize, ec := Compress(byteContainer)
	if ec != nil {
		log.Println("SaveFile Compress failed,err:", eo)
		return ec
	}
	//保存到新文件中
	ei := ioutil.WriteFile(fileAddr, afterResize, 0644)
	if ei != nil {
		log.Println("SaveFile Compress failed,err:", eo)
		return ei
	}
	return nil
}

func SaveCompressCutImg(file *multipart.FileHeader, fileAddr string) error {
	fileContent, eo := file.Open()
	if eo != nil {
		log.Println("SaveCompressCutImg file.Open failed,err:", eo)
		return eo
	}
	defer fileContent.Close()
	byteContainer, er := ioutil.ReadAll(fileContent)
	if er != nil {
		log.Println("SaveCompressCutImg ioutil.ReadAll failed,err:", er)
		return er
	}
	decodeBuf, _, eid := image.Decode(bytes.NewReader(byteContainer))
	if eid != nil {
		log.Println("SaveCompressCutImg  image.Decode failed,err:", eid)
		return eid
	}
	dsc := imaging.Fill(decodeBuf, 300, 300, imaging.Center, imaging.Lanczos)
	es := imaging.Save(dsc, fileAddr)
	if es != nil {
		log.Println("SaveCompressCutImg Save failed,err:", es)
		return eid
	}
	return nil
}
