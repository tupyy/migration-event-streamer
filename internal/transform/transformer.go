package transform

import (
	"time"

	"github.com/google/uuid"
	"github.com/tupyy/migration-event-streamer/internal/datastore/models"
)

func InventorySourceToElastic(source models.Source) models.Inventory {
	inventory := models.Inventory{
		ID:            uuid.New().String(),
		EventTime:     time.Now().Format(time.RFC3339),
		SourceID:      source.ID.String(),
		TotalCpuCores: source.Inventory.Data.Vms.CpuCores.Total,
		TotalMemory:   source.Inventory.Data.Vms.RamGB.Total,
		VMs:           source.Inventory.Data.Vms.Total,
		VMsMigratable: source.Inventory.Data.Vms.TotalMigratable,
		Datastores:    []models.Datastore{},
	}
	for _, d := range source.Inventory.Data.Infra.Datastores {
		inventory.Datastores = append(inventory.Datastores, models.Datastore{FreeCapacityGB: d.FreeCapacityGB, TotalCapacityGB: d.TotalCapacityGB, Type: d.Type})
	}
	return inventory
}
