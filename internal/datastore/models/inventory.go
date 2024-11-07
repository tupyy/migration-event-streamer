package models

import "io"

type Datastore struct {
	FreeCapacityGB  int    `json:"freeCapacityGB"`
	TotalCapacityGB int    `json:"totalCapacityGB"`
	Type            string `json:"type"`
}

type Inventory struct {
	ID            string      `json:"id"`
	EventTime     string      `json:"event_time"`
	SourceID      string      `json:"source_id"`
	TotalCpuCores int         `json:"total_cpu_cores"`
	TotalMemory   int         `json:"total_memory"`
	VMs           int         `json:"vms"`
	VMsMigratable int         `json:"vms_migratable"`
	Datastores    []Datastore `json:"datastores"`
}

type Event struct {
	Index string
	ID    string
	Body  io.ReadSeeker
}
