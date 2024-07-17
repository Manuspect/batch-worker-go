package entities

type ResponseYolo struct {
	Yolo []struct {
		Xmin       float64   `json:"xmin"`
		Ymin       float64   `json:"ymin"`
		Xmax       float64   `json:"xmax"`
		Ymax       float64   `json:"ymax"`
		Confidence float64   `json:"confidence"`
		Class      int       `json:"class"`
		Name       string    `json:"name"`
		Embedding  []float64 `json:"embedding"`
	} `json:"yolo"`
	Ocr [][]any `json:"ocr"`
}
