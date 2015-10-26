package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"

	"github.com/localhots/satan"
	"github.com/localhots/satan/example/daemons"
	"github.com/localhots/satan/example/kafka"
	"github.com/localhots/satan/stats"
)

func main() {
	var brokers string

	flag.StringVar(&brokers, "brokers", "127.0.0.1:9092", "Kafka broker addresses separated by space")
	flag.Parse()

	kafka.Initialize(strings.Split(brokers, " "))
	defer kafka.Shutdown()

	statsLogger := stats.NewStdoutLogger(0)
	defer statsLogger.Print()

	s := satan.Summon()
	s.SubscribeFunc = kafka.Subscribe
	s.DaemonStats = statsLogger

	s.AddDaemon(&daemons.NumberPrinter{})
	s.AddDaemon(&daemons.PriceConsumer{})

	s.StartDaemons()
	defer s.StopDaemons()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
