package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/3AceShowHand/workload-master/pkg/db"
	"github.com/3AceShowHand/workload-master/pkg/workload"
	"github.com/pingcap/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Option struct {
	dbOption *db.Options
	task     *workload.Options
}

var (
	dbName         string
	host           string
	port           int
	user           string
	password       string
	threads        int
	driver         string
	totalTime      time.Duration
	totalCount     int
	outputInterval time.Duration
	isolationLevel int

	maxProcessors int
	connParams    string

	globalCtx context.Context
)

func addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVar(&maxProcessors, "max-processors", 0, "runtime.GOMAXPROCS")
	cmd.PersistentFlags().StringVarP(&dbName, "db", "D", "test", "Database name")
	cmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Database host")
	cmd.PersistentFlags().StringVarP(&user, "user", "U", "root", "Database user")
	cmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Database password")
	cmd.PersistentFlags().IntVarP(&port, "port", "P", 4000, "Database port")
	cmd.PersistentFlags().IntVarP(&threads, "threads", "T", 1, "Thread concurrency")
	cmd.PersistentFlags().IntVar(&totalCount, "count", 0, "Total execution count, 0 means infinite")
	cmd.PersistentFlags().StringVarP(&driver, "driver", "d", "mysql", "Database driver: mysql")
	cmd.PersistentFlags().DurationVar(&totalTime, "time", 1<<63-1, "Total execution time")
	cmd.PersistentFlags().StringVar(&connParams, "conn-params", "", "session variables")
	cmd.PersistentFlags().IntVar(&isolationLevel, "isolation", 0, `Isolation Level 0: Default, 1: ReadUncommitted,
	2: ReadCommitted, 3: WriteCommitted, 4: RepeatableRead,
	5: Snapshot, 6: Serializable, 7: Linerizable`)
	cmd.PersistentFlags().DurationVar(&outputInterval, "interval", 10*time.Second, "Output interval time")
}

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "workload-master",
		Short: "workload-master run different kind of workload.",
		Long:  "workload-master can run different kind of custom database workload",
	}
}

func Run() {
	cmd := NewCmd()

	cmd.SetOut(os.Stdout)

	cmd.AddCommand(NewCmdVersion())
	cmd.AddCommand(NewCmdLargeTxn())
	cmd.AddCommand(NewCmdSQLCoverage())

	addFlags(cmd)

	cobra.EnablePrefixMatching = true

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

	if err := cmd.Execute(); err != nil {
		log.Warn("run command failed", zap.Error(err))
	}

	cancel()

}
