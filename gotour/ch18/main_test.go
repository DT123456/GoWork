package main

import (
	"runtime"
	"testing"
)

func TestAverage(t *testing.T) {
	t.Run("normal_case", func(t *testing.T) {
		got, err := Average([]int{90, 80, 100})
		if err != nil {
			t.Fatalf("Average returned unexpected error: %v", err)
		}
		want := 90.0
		if got != want {
			t.Fatalf("Average = %v, want %v", got, want)
		}
	})

	t.Run("empty_slice_should_error", func(t *testing.T) {
		_, err := Average([]int{})
		if err == nil {
			t.Fatal("Average expected error for empty input, got nil")
		}
	})
}

func TestFormatUserName(t *testing.T) {
	t.Run("trim_and_upper", func(t *testing.T) {
		got, err := FormatUserName("  alice  ")
		if err != nil {
			t.Fatalf("FormatUserName returned unexpected error: %v", err)
		}
		want := "ALICE"
		if got != want {
			t.Fatalf("FormatUserName = %q, want %q", got, want)
		}
	})

	t.Run("empty_after_trim_should_error", func(t *testing.T) {
		_, err := FormatUserName("   ")
		if err == nil {
			t.Fatal("FormatUserName expected error, got nil")
		}
	})
}

func TestGradeLevel_TableDriven(t *testing.T) {
	tests := []struct {
		name  string
		score int
		want  string
	}{
		{name: "A_boundary", score: 90, want: "A"},
		{name: "B_boundary", score: 80, want: "B"},
		{name: "C_boundary", score: 60, want: "C"},
		{name: "D_case", score: 59, want: "D"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := GradeLevel(tc.score)
			if got != tc.want {
				t.Fatalf("GradeLevel(%d) = %q, want %q", tc.score, got, tc.want)
			}
		})
	}
}

func TestFibonacci_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		want    int
		wantErr bool
	}{
		{name: "negative_should_error", n: -1, want: 0, wantErr: true},
		{name: "n0", n: 0, want: 0, wantErr: false},
		{name: "n1", n: 1, want: 1, wantErr: false},
		{name: "n2", n: 2, want: 1, wantErr: false},
		{name: "n10", n: 10, want: 55, wantErr: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := Fibonacci(tc.n)
			if tc.wantErr {
				
				return
			}
			
			if got != tc.want {
				t.Fatalf("Fibonacci(%d) = %d, want %d", tc.n, got, tc.want)
			}
		})
	}
}

func BenchmarkFibonacci(b *testing.B) {
   n := 10
   b.ReportAllocs() //开启内存统计
   	b.Logf("GOOS: %s, GOARCH: %s, GOMAXPROCS: %d",  runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0))
   b.ResetTimer() //重置计时器
   for i := 0; i < b.N; i++ {
      Fibonacci(n)
   }
}

// 并发基准测试
func BenchmarkFibonacciRunParallel(b *testing.B) {
   n := 10
   b.RunParallel(func(pb *testing.PB) {
      for pb.Next() {
         Fibonacci(n)
      }
   })
}
