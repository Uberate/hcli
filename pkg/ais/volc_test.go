package ais

import "os"

func TestVolcAI() {
	engine := NewVolcEngineAI(VolcConfig{
		ApiKey: os.Getenv("API_KEY"),
		T
	})
}
