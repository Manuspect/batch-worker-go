package internal

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func Upload(url string, values map[string]io.Reader) (err error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return err
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err

		}

	}
	w.Close()

	body := io.Reader(&b)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	fmt.Println(body)
	client := http.Client{
		Timeout: 33 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}

	return nil
}
