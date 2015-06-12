package tb

import (
	"sync"
	"time"

	"github.com/jaracil/clk"
)

type Tb struct {
	last  time.Duration
	size  float64
	level float64
	rate  float64
	lock  sync.Mutex
}

var ccache *clk.Cache = clk.NewCache(time.Millisecond)

// New returns new token bucket, size is the bucket size and rate
// is the rate at which the bucket is filled expressed in tokens/second.
func New(size, rate float64) *Tb {
	return &Tb{last: ccache.Lap(), size: size, level: size, rate: rate}
}

// Pull extracts n tokens from bucket and returns true if it had enough tokens.
func (p *Tb) Pull(n float64) bool {
	p.lock.Lock()
	now := ccache.Lap()
	lap := now - p.last
	p.last = now
	if p.size <= 0 {
		p.lock.Unlock()
		return true
	}
	if lap > 0 {
		p.level += lap.Seconds() * p.rate
		if p.level > p.size {
			p.level = p.size
		}
	}
	if n > p.level {
		p.lock.Unlock()
		return false
	}
	p.level -= n
	p.lock.Unlock()
	return true
}

// Last returns the last monotonic timestamp the bucket was pulled.
func (p *Tb) Last() time.Duration {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.last
}

// Size returns the bucket size
func (p *Tb) Size() float64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.size
}

// Rate returns the bucket fill rate in tokens/second.
func (p *Tb) Rate() float64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.rate
}

// Level returns the number of tokens present in the bucket.
func (p *Tb) Level() float64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	now := ccache.Lap()
	lap := now - p.last
	p.last = now
	if lap > 0 {
		p.level += lap.Seconds() * p.rate
		if p.level > p.size {
			p.level = p.size
		}
	}
	return p.level
}

// SetSize changes the bucket size (in tokens).
func (p *Tb) SetSize(v float64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.size = v
}

// SetRate changes the rate at which the bucket is filled.
func (p *Tb) SetRate(v float64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.rate = v
}

// SetLevel changes the number of tokens in the bucket.
func (p *Tb) SetLevel(v float64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.level = v
	if p.level > p.size {
		p.level = p.size
	}
}
