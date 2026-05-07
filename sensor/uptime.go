package sensor

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/host"

	"hacompanion/entity"
	"hacompanion/util"
)

type Uptime struct{}

func NewUptime() *Uptime {
	return &Uptime{}
}

func (u Uptime) Run(ctx context.Context) (*entity.Payload, error) {
	boot, err := host.BootTimeWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("host.BootTimeWithContext failed: %w", err)
	}

	start := time.Unix(int64(boot), 0)
	uptimeSeconds := time.Since(start).Seconds()

	p := entity.NewPayload()
	p.State = start.Format(time.RFC3339)
	p.Attributes["boot_time"] = start.Format(time.RFC3339)
	p.Attributes["uptime_seconds"] = util.RoundToTwoDecimals(uptimeSeconds)

	return p, nil
}
