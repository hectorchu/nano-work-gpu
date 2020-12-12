package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hectorchu/gonano/pow"
)

func main() {
	var (
		root      rootFlag
		target    difficultyFlag
		benchmark = flag.Int("b", 0, "benchmark: number of iterations")
	)
	target.Set("fffffff800000000")
	flag.Var(&root, "r", "root: 32-byte hex string")
	flag.Var(&target, "d", "difficulty: 8-byte hex string")
	flag.Parse()
	if *benchmark > 0 {
		runBenchmark(target.n, *benchmark)
		return
	}
	if root == nil {
		fmt.Println("root is required")
		os.Exit(1)
	}
	work, err := pow.GenerateGPU(root, target.n)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rev(work)
	fmt.Printf("%x\n", work)
}

func runBenchmark(target uint64, n int) {
	root := make([]byte, 32)
	var total time.Duration
	for i := 0; i < n; i++ {
		rand.Read(root)
		start := time.Now()
		work, err := pow.GenerateGPU(root, target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		elapsed := time.Now().Sub(start)
		total += elapsed
		rev(work)
		fmt.Printf("%d: %x %x %s\n", i, root, work, elapsed)
	}
	fmt.Printf("Average %s\n", total/time.Duration(n))
}

func rev(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
