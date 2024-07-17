package batch

import (
	"bufio"
	"os"
	"strings"
)

type FileCsv struct {
	Timestamp     string `csv:"timestamp"`
	Process_path  string `csv:"process_path"`
	Title         string `csv:"title"`
	Class_name    string `csv:"class_name"`
	Window_left   string `csv:"window_left"`
	Window_top    string `csv:"window_top"`
	Window_right  string `csv:"window_right"`
	Window_bottom string `csv:"window_bottom"`
	Event         string `csv:"event"`
	Mouse_x_pos   string `csv:"mouse_x_pos"`
	Mouse_y_pos   string `csv:"mouse_y_pos"`
	Modifiers     string `csv:"modifiers"`
}

func GetEvetsCsv(pathCsv string) ([]FileCsv, error) {
	file, err := os.Open(pathCsv)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []FileCsv

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		lineArr := strings.Split(line, ",")

		events = append(events, FileCsv{
			Timestamp:     lineArr[0],
			Process_path:  lineArr[1],
			Title:         lineArr[2],
			Class_name:    lineArr[3],
			Window_left:   lineArr[4],
			Window_top:    lineArr[5],
			Window_right:  lineArr[6],
			Window_bottom: lineArr[7],
			Event:         lineArr[8],
			Mouse_x_pos:   lineArr[9],
			Mouse_y_pos:   lineArr[10],
			Modifiers:     lineArr[11],
		})

	}
	return events, err
}
