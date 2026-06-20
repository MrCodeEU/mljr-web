package scheduler

import (
	"context"
	"hash/fnv"
	"sync"
	"time"
)

// HoroscopeProvider returns a short daily horoscope blurb for a zodiac
// sign. Implementations must be safe for concurrent use and should treat
// any failure (timeout, unreachable API, etc.) as non-fatal — callers
// swallow errors and simply omit that sign's blurb from the edition.
type HoroscopeProvider interface {
	Daily(ctx context.Context, sign string, date time.Time) (string, error)
}

// horoscopeLines is a small fixed pool of generic, sign-agnostic affirming
// lines. staticHoroscopeProvider picks one deterministically per (sign,
// date) so a sign gets a stable line for the whole day, varying day to day.
var horoscopeLines = []string{
	"Today favors small, steady steps over grand gestures — trust the process.",
	"An unexpected conversation could open a door you weren't watching.",
	"Your energy is contagious today — use it to lift someone else up.",
	"Patience pays off now; what feels stalled is quietly building momentum.",
	"A good day to tidy up loose ends before starting something new.",
	"Trust your first instinct today — overthinking will only muddy it.",
	"Connection matters more than productivity today — reach out to someone.",
	"Something you've been putting off is easier than you're imagining.",
	"Your curiosity is your sharpest tool today — follow it somewhere new.",
	"A quiet day is still a productive one — rest counts as progress too.",
}

// staticHoroscopeProvider is the default HoroscopeProvider: a local,
// network-free generator. No currently-verified free, reliable, no-signup
// horoscope API exists (the usual answer, aztro, is defunct), so this is
// the safe always-available fallback rather than a fabricated dependency.
type staticHoroscopeProvider struct{}

func (staticHoroscopeProvider) Daily(_ context.Context, sign string, date time.Time) (string, error) {
	h := fnv.New32a()
	_, _ = h.Write([]byte(sign + "|" + date.Format("2006-01-02")))
	return horoscopeLines[int(h.Sum32())%len(horoscopeLines)], nil
}

// cachingHoroscopeProvider memoizes Daily by (sign, date) so multiple
// members sharing a sign — across one group's compose or across groups in
// the same cron tick — only resolve each sign once per day. Mostly
// future-proofing for when a real HTTP-backed provider is wired in below;
// the static provider has no real cost to repeat, but the cache costs
// nothing either.
type cachingHoroscopeProvider struct {
	inner HoroscopeProvider
	mu    sync.Mutex
	cache map[string]string
}

func newCachingHoroscopeProvider(inner HoroscopeProvider) *cachingHoroscopeProvider {
	return &cachingHoroscopeProvider{inner: inner, cache: map[string]string{}}
}

func (p *cachingHoroscopeProvider) Daily(ctx context.Context, sign string, date time.Time) (string, error) {
	key := sign + "|" + date.Format("2006-01-02")

	p.mu.Lock()
	if blurb, ok := p.cache[key]; ok {
		p.mu.Unlock()
		return blurb, nil
	}
	p.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	blurb, err := p.inner.Daily(ctx, sign, date)
	if err != nil {
		return "", err
	}

	p.mu.Lock()
	p.cache[key] = blurb
	p.mu.Unlock()
	return blurb, nil
}

// defaultHoroscopeProvider builds the provider RunScan uses each tick: the
// static local generator wrapped in the shared per-day cache. A real
// HTTP-backed provider can be substituted here later once a reliable free
// API is verified — intentionally left unconfigured until then.
func defaultHoroscopeProvider() HoroscopeProvider {
	return newCachingHoroscopeProvider(staticHoroscopeProvider{})
}
