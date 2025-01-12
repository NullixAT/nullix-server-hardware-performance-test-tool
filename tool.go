package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {

	cpuTaskRuntime := 5000

	for index, value := range os.Args {
		if value == "--cpu-runtime-ms" {
			ms, err := strconv.Atoi(os.Args[index+1])
			if err != nil {
				panic(err)
			}
			cpuTaskRuntime = ms
		}
	}

	start := time.Now()

	if cpuTaskRuntime != 0 {
		total := 0
		for true {
			cpuTask()
			total++
			if cpuTaskRuntime > 0 && cpuTaskRuntime <= int(time.Since(start).Milliseconds()) {
				fmt.Print(total)
				break
			}
		}
	}
}

func cpuTask() float64 {
	return math.Sqrt(rand.Float64()/rand.Float64()*rand.Float64() + rand.Float64() - rand.Float64())
}
