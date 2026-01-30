package models

type EncryptRequest struct {
	DataSourcePath     string `json:"dataSourcePath"`
	CarrierDir         string `json:"carrierDir"`
	CarrierImagePath   string `json:"carrierImagePath"`
	OutputDir          string `json:"outputDir"`
	OutputFileName     string `json:"outputFileName"`
	Password           string `json:"password"`
	Scatter            *bool  `json:"scatter"`
	Identifier         string `json:"identifier"`
	AutoSelectCarrier  bool   `json:"autoSelectCarrier"`
	PreferLargestImage bool   `json:"preferLargestImage"`
}

type DecryptRequest struct {
	ImagePath  string `json:"imagePath"`
	OutputDir  string `json:"outputDir"`
	Password   string `json:"password"`
	Identifier string `json:"identifier"`
}

type GenerateRequest struct {
	OutputDir    string `json:"outputDir"`
	TargetBytes  int64  `json:"targetBytes"`
	Count        int    `json:"count"`
	Prefix       string `json:"prefix"`
	RandomSeed   int64  `json:"randomSeed"`
	NoiseEnabled bool   `json:"noiseEnabled"`
}

type ProgressEvent struct {
	TaskID   string `json:"taskId"`
	Progress int    `json:"progress"`
	Message  string `json:"message"`
	Current  int    `json:"current"`
	Total    int    `json:"total"`
	Error    string `json:"error,omitempty"`
	Done     bool   `json:"done,omitempty"`
}

type AppInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	BuildHash string `json:"buildHash"`
	Author    string `json:"author"`
	GitHub    string `json:"github"`
}
