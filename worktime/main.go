package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
)
const (
	filePath = "./tmp/images"
)
func ConvertSVGToPNG(ctx context.Context, svg io.Reader) (io.Reader, error) {
	makePathIfNotExists(filePath)
	fileName := filePath + "/test"
	fileNameSvg := fileName + ".svg"
	fileNamePng := fileName + ".png"

	file, err := os.Create(fileNameSvg)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, svg)
	if err != nil {
		return nil, err
	}

	_, err = exec.Command("rsvg-convert", fileNameSvg, "-h","280","-w","482","-o",fileNamePng).Output()
	if err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	go func(_ctx context.Context) {
		defer pw.Close()

		png, _ := os.Open(fileNamePng)
		_, _ = io.Copy(pw, png)
	}(ctx)
	return pr, nil
}


func makePathIfNotExists(pathName string) {
	_, err := os.Stat(pathName)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(pathName, os.ModePerm)
	} else if os.IsExist(err) {
		os.Remove(pathName)
	}
}

func main() {
	cli := http.DefaultClient
	resp, err := cli.Get("http://home-store-img.uniontech.com/svg/5503c3819e77454a82af2a6f84f9aa98.svg")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	_, err = ConvertSVGToPNG(context.Background(),bytes.NewReader(bs))
	if err != nil {
		panic(err)
	}
}