package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	logFi "github.com/gofiber/fiber/v2/log"

	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"processing-worker/internal/entities"
	"strings"
)

func main() {
	filePath := "/home/user/Screenshot 2024-07-10 143615.png"
	// filePath := "/home/user/Изображения/cat.jpg"
	// filePath := "/home/user/Изображения/dog_1.png"
	addr := "http://localhost:8001/yolo/img_object_detection_to_json"
	res, err := doUpload(addr, filePath)
	if err != nil {
		logFi.Error("doUpload", err)
		fmt.Printf("upload file [%s] error: %s", filePath, err)
		return
	}
	fmt.Printf("file: %v\n", *res)
	fmt.Printf("upload file [%s] ok\n", filePath)
}

func createReqBody(filePath string) (string, io.Reader, error) {
	var err error

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf) // body writer

	f, err := os.Open(filePath)
	if err != nil {
		logFi.Error("createReqBody", err)
		return "", nil, err
	}
	defer f.Close()

	_, fileName := filepath.Split(filePath)
	fw1, err := CreateFormFileImage(bw, "file", fileName)
	if err != nil {
		logFi.Error("createReqBody", err)
		return "", nil, err
	}

	io.Copy(fw1, f)

	bw.Close() //write the tail boundry
	return bw.FormDataContentType(), buf, nil
}

func doUpload(addr, filePath string) (*entities.ResponseYolo, error) {

	contType, reader, err := createReqBody(filePath)
	if err != nil {
		logFi.Error("doUpload", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", addr, reader)
	if err != nil {
		logFi.Error("doUpload", "request send error:", err)
		return nil, err
	}

	req.Header.Add("Content-Type", contType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logFi.Error("doUpload", "response error:", err)
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logFi.Error("doUpload", "error:", err)
		return nil, err
	}

	var res *entities.ResponseYolo
	err = json.Unmarshal(b, &res)
	if err != nil {
		logFi.Error("doUpload", "can't unmarshal response send error:", err)
		return nil, err
	}

	resp.Body.Close()

	return res, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func CreateFormFileImage(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	arr := strings.Split(filename, ".")

	name1 := "application"
	name2 := "octet-stream"

	if len(arr) > 1 {
		name1 = "image"
		name2 = arr[len(arr)-1]

	}

	h.Set("Content-Type", fmt.Sprintf("%s/%s", name1, name2))
	return w.CreatePart(h)
}
