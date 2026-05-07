package sensor

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/load"

	"hacompanion/entity"
	"hacompanion/util"
)

type LoadAVG struct{}

func NewLoadAVG() *LoadAVG {
	return &LoadAVG{}
}

func (w LoadAVG) Run(ctx context.Context) (*entity.Payload, error) {
	avg, err := load.AvgWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("load.AvgWithContext failed: %w", err)
	}

	p := entity.NewPayload()
	p.State = util.RoundToTwoDecimals(avg.Load1)
	p.Attributes["5m"] = util.RoundToTwoDecimals(avg.Load5)
	p.Attributes["15m"] = util.RoundToTwoDecimals(avg.Load15)

	return p, nil
}
