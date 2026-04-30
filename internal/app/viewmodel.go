package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PageData struct {
	Network Network
	Stats   StatsResponse
	Error   string
	Now     time.Time
}

func (p PageData) Sequence() int64 {
	if p.Stats.Hero.LatestLedger.Sequence != 0 {
		return p.Stats.Hero.LatestLedger.Sequence
	}
	return p.Stats.Header.LatestLedgerSequence
}

func (p PageData) ClosedAt() time.Time {
	if !p.Stats.Hero.LatestLedger.ClosedAt.IsZero() {
		return p.Stats.Hero.LatestLedger.ClosedAt
	}
	return p.Stats.Header.LatestLedgerClosedAt
}

func (p PageData) HasClosedAt() bool {
	return !p.ClosedAt().IsZero()
}

func (p PageData) HealthStatus() string {
	if p.Stats.Hero.Health.Status == "" {
		return "unknown"
	}
	return p.Stats.Hero.Health.Status
}

func (p PageData) IsCaughtUp() bool {
	if strings.EqualFold(p.HealthStatus(), "halted") {
		return false
	}
	if !p.HasClosedAt() {
		return true
	}
	return p.Now.Sub(p.ClosedAt()) <= 2*time.Minute
}

func (p PageData) Lag() string {
	if !p.HasClosedAt() {
		return "unknown"
	}
	return HumanDuration(p.Now.Sub(p.ClosedAt()))
}

func FormatInt(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	pre := len(s) % 3
	if pre == 0 {
		pre = 3
	}
	b.WriteString(s[:pre])
	for i := pre; i < len(s); i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}
	return b.String()
}

func FormatFloat(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.UTC().Format("2006-01-02 15:04:05 UTC")
}

func ShortHash(h string) string {
	if len(h) <= 20 {
		return h
	}
	return h[:10] + "…" + h[len(h)-10:]
}

func HumanDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	if d < 48*time.Hour {
		return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
	}
	return fmt.Sprintf("%dd %dh", int(d.Hours())/24, int(d.Hours())%24)
}
