package batch

import (
	"encoding/json"
	"io"
	"os"
)

type FileJson struct {
	Id               int    `json:"id"`
	BatchId          int    `json:"batch_id"`
	UserId           int    `json:"user_id"`
	LogRecordCounter int    `json:"log_record_counter"`
	Timestamp        string `json:"timestamp"`
	WindowLeft       int    `json:"window_left"`
	WindowTop        int    `json:"window_top"`
	WindowRight      int    `json:"window_right"`
	WindowBottom     int    `json:"window_bottom"`
	MouseX           int    `json:"mouse_x"`
	MouseY           int    `json:"mouse_y"`
	ProcessPath      string `json:"process_path"`
	Title            string `json:"title"`
	ClassName        string `json:"class_name"`
	Event            string `json:"event"`
	Modifiers        string `json:"modifiers"`
}

type InfoJson struct {
	Id               int    `json:"id"`
	BatchId          int    `json:"batch_id"`
	UserId           int    `json:"user_id"`
	LogRecordCounter int    `json:"log_record_counter"`
	Timestamp        string `json:"timestamp"`
	WindowLeft       int    `json:"window_left"`
	WindowTop        int    `json:"window_top"`
	WindowRight      int    `json:"window_right"`
	WindowBottom     int    `json:"window_bottom"`
	MouseX           int    `json:"mouse_x"`
	MouseY           int    `json:"mouse_y"`
	ProcessPath      string `json:"process_path"`
	Title            string `json:"title"`
	ClassName        string `json:"class_name"`
	Event            string `json:"event"`
	Modifiers        string `json:"modifiers"`
	FilePath         string `json:"filePath"`
}

func GetEventsJson(pathJson string) ([]FileJson, error) {
	file, err := os.Open(pathJson)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var events []FileJson

	err = json.Unmarshal(data, &events)

	if err != nil {
		return nil, err
	}

	return events, err
}
