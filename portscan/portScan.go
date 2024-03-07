package portscan

import (
	"net"
	"strconv"
	"sync"
	"time"
)

type Results struct {
	Port  string
	State string
}

func scanPort(protocol string, target string, port int, wg *sync.WaitGroup, results chan<- Results) {
	defer wg.Done()
	addr := target + ":" + strconv.Itoa(port)
	res := Results{Port: strconv.Itoa(port)}
	conn, err := net.DialTimeout(protocol, addr, 700*time.Millisecond)

	if err != nil {
		res.State = "Closed"
	} else {
		defer conn.Close()
		res.State = "Open"
	}
	results <- res
}

func LoopScan(protocol string, target string, inport int, outport int) []Results {
	var results []Results
	var wg sync.WaitGroup
	resultsChain := make(chan Results)

	for i := inport; i <= outport; i++ {
		wg.Add(1)
		go scanPort(protocol, target, i, &wg, resultsChain)
	}
	go func() {
		wg.Wait()
		close(resultsChain)
	}()
	for res := range resultsChain {
		results = append(results, res)
	}
	return results
}
