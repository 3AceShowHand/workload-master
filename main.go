package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/3AceShowHand/workload-master/cmd"
	"github.com/spf13/cobra"
)

var (
	dbName         string
	host           string
	port           int
	user           string
	password       string
	threads        int
	acThreads      int
	driver         string
	totalTime      time.Duration
	totalCount     int
	dropData       bool
	ignoreError    bool
	outputInterval time.Duration
	isolationLevel int
	silence        bool
	pprofAddr      string
	metricsAddr    string
	maxProcs       int
	connParams     string

	globalDB  *sql.DB
	globalCtx context.Context
)

func main() {
	var rootCmd = &cobra.Command {
		Use: "workload-master",
		Short: "workload-master master all different kind of workload.",
		Long: "workload-master can run different kind of database workload, such as go-tpc, go-ycsb, sysbench",
	}

	rootCmd.PersistentFlags().StringVarP(&dbName, "db", "D", "test", "Database name")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Database host")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "U", "root", "Database user")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Database password")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "P", 4000, "Database port")
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "T", 1, "Thread concurrency")

	cobra.EnablePrefixMatching = true

	cmd.RegisterVersionInfo(rootCmd)
	// registerTpcc(rootCmd)
	// registerTpch(rootCmd)
	// registerCHBenchmark(rootCmd)

	var cancel context.CancelFunc
	globalCtx, cancel = context.WithCancel(context.Background())

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	closeDone := make(chan struct{}, 1)
	go func() {
		sig := <-sc
		fmt.Printf("\nGot signal [%v] to exit.\n", sig)
		cancel()

		select {
		case <-sc:
			// send signal again, return directly
			fmt.Printf("\nGot signal [%v] again to exit.\n", sig)
			os.Exit(1)
		case <-time.After(10 * time.Second):
			fmt.Print("\nWait 10s for closed, force exit\n")
			os.Exit(1)
		case <-closeDone:
			return
		}
	}()

	rootCmd.Execute()

	cancel()
}