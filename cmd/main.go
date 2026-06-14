package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	USAGE_LIMIT        = 85
	RAM_SIZE_KB        int
	RAM_USAGE_LIMIT_KB int

	CHECK_INTERVAL_SECONDS time.Duration = 15 * time.Second
)

func main() {
	result := exec.Command("grep", "MemTotal", "/proc/meminfo")
	output, err := result.Output()
	if err != nil {
		panic(err)
	}
	args := strings.Split(string(output), " ")
	size_int, err := strconv.Atoi(strings.TrimSpace(args[8]))
	if err != nil {
		panic(err)
	}
	RAM_SIZE_KB = size_int
	fmt.Printf("Total RAM SIZE (KB): %d\n", RAM_SIZE_KB)

	usage_limit := RAM_SIZE_KB * USAGE_LIMIT / 100
	RAM_USAGE_LIMIT_KB = usage_limit
	fmt.Printf("RAM USAGE LIMIT (KB): %d%% ( %d of %d) \n", USAGE_LIMIT, RAM_USAGE_LIMIT_KB, RAM_SIZE_KB)

	// CHECK PROCESS USING MOST RAM
	ctx, cancel := context.WithCancel(context.Background())

	go monitorRamUsage(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	fmt.Println("Received signal to stop monitoring.")
	cancel()
	fmt.Printf("Exiting in %s \n ", CHECK_INTERVAL_SECONDS.String())
	time.Sleep(CHECK_INTERVAL_SECONDS)
	fmt.Println("Exited")
}

func monitorRamUsage(ctx context.Context) {
	ticker := time.NewTicker(CHECK_INTERVAL_SECONDS)
	runCheck() // Run immediately on start
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping RAM usage monitoring...")
			return
		case <-ticker.C:
			runCheck()
		}
	}
}

func runCheck() {
	currentUsage := getCurrentRamUsage()

	if checkRamUsageExceedsLimit(currentUsage) {
		usage_percent := float64(currentUsage) / float64(RAM_SIZE_KB) * 100
		exec.Command("notify-send", "RAM Usage Alert", fmt.Sprintf("Current: %.2f%% ( %d KB of %d KB), Limit: %d%%", usage_percent, currentUsage, RAM_SIZE_KB, USAGE_LIMIT)).Run()
		getTopRamConsumingProcesses()
	}
}
func getCurrentRamUsage() int {
	result := exec.Command("grep", "MemAvailable", "/proc/meminfo")
	output, err := result.Output()
	if err != nil {
		panic(err)
	}

	s := strings.Split(string(output), " ")
	available_ram_kb, err := strconv.Atoi(strings.TrimSpace(s[4]))
	if err != nil {
		panic(err)
	}
	used_ram_kb := RAM_SIZE_KB - available_ram_kb
	return used_ram_kb
}

func checkRamUsageExceedsLimit(currentUsage int) bool {
	return currentUsage > RAM_USAGE_LIMIT_KB
}

func getTopRamConsumingProcesses() {
	exec.Command("")
}
