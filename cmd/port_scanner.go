package cmd

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/pda13/gonet/pkg/slices"
	"github.com/spf13/cobra"
)

// gonet scan-ports --host localhost --from 1 --to 99999

var portScannerCmd = &cobra.Command{
	Use:   "scan-ports",
	Short: "Scan range of ports on a specific host",
	Run:   GetScanPortCommand(),
}

func init() {
	rootCmd.AddCommand(portScannerCmd)

	portScannerCmd.Flags().String("host", "localhost", "host")
	portScannerCmd.Flags().Int("from", 1, "port from")
	portScannerCmd.Flags().Int("to", 9090, "port to")
}

func GetScanPortCommand() func(cmd *cobra.Command, args []string) {
	ps := NewPortScanner(10)
	return ps.GetCommand()
}

type PortScanner struct {
	workersNum int
	wg         sync.WaitGroup
	openPorts  chan int
	jobs       chan scanPortJob
}

type scanPortJob struct {
	host string
	port int
}

func NewPortScanner(workersNum int) *PortScanner {
	return &PortScanner{
		workersNum: workersNum,
		openPorts:  make(chan int),
		jobs:       make(chan scanPortJob),
	}
}

func (ps *PortScanner) Run() {
	ps.wg.Add(ps.workersNum)
	for i := 0; i < ps.workersNum; i++ {
		go ps.worker()
	}
}

func (ps *PortScanner) worker() {
	defer ps.wg.Done()
	for job := range ps.jobs {
		address := net.JoinHostPort(job.host, strconv.Itoa(job.port))
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)

		if err == nil {
			conn.Close()
			ps.openPorts <- job.port
		}
	}
}

func (ps *PortScanner) GetCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Fatal(err)
		}

		from, err := cmd.Flags().GetInt("from")
		if err != nil {
			log.Fatal(err)
		}

		to, err := cmd.Flags().GetInt("to")
		if err != nil {
			log.Fatal(err)
		}

		go ps.Run()

		go func() {
			ps.wg.Wait()
			close(ps.openPorts)
		}()

		// Send jobs
		go func() {
			for port := from; port <= to; port++ {
				ps.jobs <- scanPortJob{host: host, port: port}
			}
			close(ps.jobs)
		}()

		openPorts := []int{}
		for openPort := range ps.openPorts {
			openPorts = append(openPorts, openPort)
		}

		fmt.Println("open ports:")
		slices.PrettyPrint(openPorts)
	}
}
