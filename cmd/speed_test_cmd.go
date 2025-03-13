package cmd

import (
	"fmt"

	"github.com/showwin/speedtest-go/speedtest"
	"github.com/spf13/cobra"
)

var speedTestCmd = &cobra.Command{
	Use:   "speed-test",
	Short: "Test download and upload speed",
	Run:   GetSpeedTestCommand(),
}

func init() {
	rootCmd.AddCommand(speedTestCmd)

	speedTestCmd.Flags().Int("upload-size-mb", 10, "set size of uploading file")
}

func GetSpeedTestCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var speedtestClient = speedtest.New()

		serverList, _ := speedtestClient.FetchServers()
		targets, _ := serverList.FindServer([]int{})

		for _, s := range targets {
			s.PingTest(nil)
			s.DownloadTest()
			s.UploadTest()
			fmt.Printf("Latency: %s, Download: %s, Upload: %s\n", s.Latency, s.DLSpeed, s.ULSpeed)
			s.Context.Reset()
		}
	}
}
