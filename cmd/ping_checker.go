package cmd

import (
	"fmt"
	"log"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/spf13/cobra"
)

// gonet ping-check --ip google.com

var pingCheckerCmd = &cobra.Command{
	Use:   "ping-check",
	Short: "Check host availability by ICMP-requests to IP addresses",
	Run:   GetPinkCheckerCommand(),
}

func init() {
	rootCmd.AddCommand(pingCheckerCmd)

	pingCheckerCmd.Flags().String("ip", "ya.ru", "ip")
}

func GetPinkCheckerCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ip, err := cmd.Flags().GetString("ip")
		if err != nil {
			log.Fatal(err)
		}

		pinger, err := probing.NewPinger(ip)
		if err != nil {
			log.Fatal(err)
		}

		pinger.Count = 4                 // Количество ICMP-запросов
		pinger.Timeout = time.Second * 5 // Таймаут для всего процесса
		pinger.Interval = time.Second    // Интервал между запросами

		fmt.Printf("PING %s (%s):\n", ip, pinger.IPAddr())
		pinger.Run()

		// Получаем статистику
		stats := pinger.Statistics()
		fmt.Printf("\n--- ping statistics for %s ---\n", ip)
		fmt.Printf("%d packet sent, %d reveived, %v%% loss\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("Receive-sent time: min/avg/max = %v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
	}
}
