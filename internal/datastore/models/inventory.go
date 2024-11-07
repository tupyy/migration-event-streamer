package models

import "io"

type InventoryElastic struct {
	ID            string `json:"id"`
	SourceID      string `json:"source_id"`
	TotalCpuCores int    `json:"total_cpu_cores"`
}

type Event struct {
	Index string
	ID    string
	Body  io.ReadSeeker
}
