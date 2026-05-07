package sensor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/sensors"

	"hacompanion/entity"
	"hacompanion/util"
)

type CPUTemp struct {
	UseCelsius bool
}

func NewCPUTemp(m entity.Meta) *CPUTemp {
	c := &CPUTemp{}
	if m.GetBool("celsius") {
		c.UseCelsius = true
	}
	return c
}

func (c CPUTemp) Run(ctx context.Context) (*entity.Payload, error) {
	temps, err := sensors.TemperaturesWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("sensors.TemperaturesWithContext failed: %w", err)
	}

	p := entity.NewPayload()
	var fallback string

	for _, t := range temps {
		if t.Temperature <= 0 {
			continue
		}

		temp := t.Temperature
		if !c.UseCelsius {
			temp = temp*9/5 + 32
		}

		value := fmt.Sprintf("%.1f", temp)
		key := util.ToSnakeCase(t.SensorKey)
		p.Attributes[key] = value

		if fallback == "" {
			fallback = value
		}

		sensorKey := strings.ToLower(t.SensorKey)
		if p.State == "" &&
			(strings.Contains(sensorKey, "cpu") ||
				strings.Contains(sensorKey, "package") ||
				strings.Contains(sensorKey, "tctl") ||
				sensorKey == "tc0p" ||
				sensorKey == "tc0d" ||
				sensorKey == "tc0h") {
			p.State = value
		}
	}

	if p.State == "" {
		p.State = fallback
	}
	if p.State == "" {
		return nil, fmt.Errorf("no valid temperatures found")
	}

	return p, nil
}

type CPUUsage struct{}

func NewCPUUsage() *CPUUsage {
	return &CPUUsage{}
}

func (c CPUUsage) Run(ctx context.Context) (*entity.Payload, error) {
	p := entity.NewPayload()

	total, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("cpu.PercentWithContext(total) failed: %w", err)
	}
	if len(total) > 0 {
		p.State = util.RoundToTwoDecimals(total[0])
	}

	cores, err := cpu.PercentWithContext(ctx, time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("cpu.PercentWithContext(cores) failed: %w", err)
	}
	for i, v := range cores {
		p.Attributes[fmt.Sprintf("core_%d", i)] = util.RoundToTwoDecimals(v)
	}

	return p, nil
}
