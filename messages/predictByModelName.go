package messages

type PredictByModelName struct {
	BatchSize             uint     `json:"BatchSize"`
	DesiredResultModality string   `json:"DesiredResultModality"`
	Inputs                []string `json:"InputFiles"`
	ModelName             string   `json:"ModelName"`
	Warmups               uint     `json:"NumWarmup"`
	TraceLevel            string   `json:"TraceLevel"`
	TracerAddress         string   `json:"TracerAddress"`
	UseGpu                bool     `json:"UseGpu"`
}
