package app

import (
	"fmt"
	"time"
)

type PerfLogger func(module, action, details string)

func logPerf(logf PerfLogger, module, taskID, stage string, d time.Duration, extra string) {
	if logf == nil {
		return
	}
	details := fmt.Sprintf("TaskID: %s | Stage: %s | Duration: %s", taskID, stage, formatDuration(d))
	if extra != "" {
		details += " | " + extra
	}
	logf(module, "PERF", details)
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	secs := d.Seconds()
	if secs < 10 {
		return fmt.Sprintf("%.2fs", secs)
	}
	return fmt.Sprintf("%.1fs", secs)
}
