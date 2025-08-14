package basic_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TimeTruncate(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)

	// UTC+8 2025-08-13 01:02:03.000000004
	// UTC   2025-08-12 17:02:03.000000004
	d := time.Date(2025, 8, 13, 1, 2, 3, 4, loc)

	tests := []struct {
		name     string
		duration time.Duration
		want     string // form "2006-01-02 15:04:05Z08"
	}{
		{
			name:     "truncate_to_hour",
			duration: time.Hour,
			want:     "2025-08-13T01:00:00+08:00",
		}, {
			name:     "truncate_to_30_minutes",
			duration: 30 * time.Minute,
			want:     "2025-08-13T01:00:00+08:00",
		}, {
			name:     "truncate_to_minute",
			duration: time.Minute,
			want:     "2025-08-13T01:02:00+08:00",
		}, {
			name:     "truncate_to_second",
			duration: time.Second,
			want:     "2025-08-13T01:02:03+08:00",
		}, {
			name:     "truncate_to_day",
			duration: 24 * time.Hour,
			want:     "2025-08-12T08:00:00+08:00",
		}, {
			name:     "truncate_to_2_hours",
			duration: 2 * time.Hour,
			want:     "2025-08-13T00:00:00+08:00",
		}, {
			name:     "truncate_to_5_minutes",
			duration: 5 * time.Minute,
			want:     "2025-08-13T01:00:00+08:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.Truncate(tt.duration).Format(time.RFC3339)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_TimeTruncateWithTimezone(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)

	tests := []struct {
		name string

		date       time.Time
		duration   time.Duration
		wantFormat string // form "2006-01-02 15:04:05Z08"
		wantUnix   int64
	}{
		{
			name: "UTC 和 UTC+8 时间在同一天",
			// UTC+8 2025-08-13 09:02:03.000000004
			// UTC   2025-08-13 01:02:03.000000004
			// UNIX 时间戳 1755046923
			date:       time.Date(2025, 8, 13, 9, 2, 3, 4, loc),
			duration:   24 * time.Hour,
			wantFormat: "2025-08-13T08:00:00+08:00",
			wantUnix:   1755043200, // 差值 3723 = 1h 2m 3s
		},
		{
			name: "UTC 和 UTC+8 时间在不同天",
			// UTC+8 2025-08-13 01:02:03.000000004
			// UTC   2025-08-12 17:02:03.000000004
			// UNIX 时间戳 1755018123
			date:       time.Date(2025, 8, 13, 1, 2, 3, 4, loc),
			duration:   24 * time.Hour,
			wantFormat: "2025-08-12T08:00:00+08:00",
			wantUnix:   1754956800, // 差值 61323 = 17h 2m 3s
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("date: %v", tt.date.Unix())

			gotT := tt.date.Truncate(tt.duration)
			// 对比 format 后的字符串
			got := gotT.Format(time.RFC3339)
			assert.Equal(t, tt.wantFormat, got)

			gotUnix := gotT.Unix()
			assert.Equal(t, tt.wantUnix, gotUnix)
		})
	}
}
