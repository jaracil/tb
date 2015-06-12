package tb

import (
	"testing"
	"time"
)

func TestTb(t *testing.T) {
	tkb := New(10, 10)
	for x := 0; x < 10; x++ {
		if !tkb.Pull(1) {
			t.Fatal("bucket empty!")
		}
	}
	if tkb.Pull(1) {
		t.Fatal("bucket not empty!")
	}
	time.Sleep(500 * time.Millisecond)
	if tkb.Level() < tkb.Size()*0.5 {
		t.Fatal("bucket not filled fast enough", tkb.Level())
	}
	if tkb.Level() > tkb.Size()*0.75 {
		t.Fatal("bucket filled too fast", tkb.Level())
	}
	time.Sleep(500 * time.Millisecond)
	if tkb.Level() != tkb.Size() {
		t.Fatal("bucket not full", tkb.Level())
	}
}

func BenchmarkTb(b *testing.B) {
	t := New(1000000, 1000000)
	for i := 0; i < b.N; i++ {
		t.Pull(100)
	}
}
