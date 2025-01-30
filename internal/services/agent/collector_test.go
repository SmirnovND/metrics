package agent

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"runtime"
	"testing"
)

func BenchmarkUpdateMemoryUsage(b *testing.B) {
	metrics := &domain.Metrics{
		Data: make(map[string]*domain.Metric),
	}

	var memStatsBefore, memStatsAfter runtime.MemStats

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&memStatsBefore)
		Update(metrics, BaseMetric)
		runtime.ReadMemStats(&memStatsAfter)
		b.ReportMetric(float64(memStatsAfter.Alloc-memStatsBefore.Alloc), "bytes_alloc")
	}
}
