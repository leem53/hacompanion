package sensor

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"

	"hacompanion/entity"
	"hacompanion/util"
)

type Memory struct{}

func NewMemory() *Memory {
	return &Memory{}
}

func (m Memory) Run(ctx context.Context) (*entity.Payload, error) {
	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("mem.VirtualMemoryWithContext failed: %w", err)
	}

	swap, err := mem.SwapMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("mem.SwapMemoryWithContext failed: %w", err)
	}

	p := entity.NewPayload()

	p.State = util.RoundToTwoDecimals(float64(vm.Available) / 1024 / 1024 / 1024)
	p.Attributes["mem_available"] = util.RoundToTwoDecimals(float64(vm.Available) / 1024 / 1024 / 1024)
	p.Attributes["mem_total"] = util.RoundToTwoDecimals(float64(vm.Total) / 1024 / 1024 / 1024)
	p.Attributes["swap_free"] = util.RoundToTwoDecimals(float64(swap.Free) / 1024 / 1024 / 1024)
	p.Attributes["swap_total"] = util.RoundToTwoDecimals(float64(swap.Total) / 1024 / 1024 / 1024)

	return p, nil
}
