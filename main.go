package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hectorchu/gonano/pow"
)

func main() {
	var (
		root      rootFlag
		target    difficultyFlag
		benchmark = flag.Int("b", 0, "benchmark: number of iterations")
		server    = flag.String("s", "", "server: address to listen for RPC requests")
	)
	target.Set("fffffff800000000")
	flag.Var(&root, "r", "root: 32-byte hex string")
	flag.Var(&target, "d", "difficulty: 8-byte hex string")
	flag.Parse()
	if *benchmark > 0 {
		runBenchmark(target.n, *benchmark)
		return
	}
	if *server != "" {
		http.HandleFunc("/", rpcHandler)
		if err := http.ListenAndServe(*server, nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	var (
		buf    bytes.Buffer
		v      struct{ Action, Difficulty, Hash string }
		root   rootFlag
		target difficultyFlag
	)
	io.Copy(&buf, r.Body)
	r.Body.Close()
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		return
	}
	if v.Action != "work_generate" {
		return
	}
	if err := root.Set(v.Hash); err != nil {
		return
	}
	target.Set("fffffff800000000")
	if v.Difficulty != "" {
		if err := target.Set(v.Difficulty); err != nil {
			return
		}
	}
	work, err := pow.GenerateGPU(root, target.n)
	if err != nil {
		return
	}
	rev(work)
	buf.Reset()
	json.NewEncoder(&buf).Encode(map[string]string{
		"work": hex.EncodeToString(work),
	})
	io.Copy(w, &buf)
}

func rev(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
