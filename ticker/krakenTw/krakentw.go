package krakentw

import (
	convertertw "main/converterTw"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

var dinle string
var cst = make(chan string)

func Krken(askPrc string) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	kraken := ws.NewKraken(ws.ProdBaseURL)
	if err := kraken.Connect(); err != nil {
		log.Fatalf("Error connecting to web socket: %s", err.Error())
	}

	// subscribe to BTCUSD`s book
	if err := kraken.SubscribeBook([]string{ws.BTCUSD}, ws.Depth10); err != nil {
		log.Fatalf("SubscribeTicker error: %s", err.Error())
	}

	// subscribe to BTCUSD`s candles
	if err := kraken.SubscribeCandles([]string{ws.BTCUSD}, ws.Interval1440); err != nil {
		log.Fatalf("SubscribeTicker error: %s", err.Error())
	}

	// subscribe to BTCUSD`s ticker
	if err := kraken.SubscribeTicker([]string{ws.BTCUSD}); err != nil {
		log.Fatalf("SubscribeTicker error: %s", err.Error())
	}

	// subscribe to BTCUSD`s trades
	if err := kraken.SubscribeTrades([]string{ws.BTCUSD}); err != nil {
		log.Fatalf("SubscribeTicker error: %s", err.Error())
	}

	// subscribe to BTCUSD`s spread
	if err := kraken.SubscribeSpread([]string{ws.BTCUSD}); err != nil {
		log.Fatalf("SubscribeTicker error: %s", err.Error())
	}

	for {
		select {
		case <-signals:
			log.Warn("Stopping...")
			if err := kraken.Close(); err != nil {
				log.Fatal(err)
			}
			return
		case update := <-kraken.Listen():
			switch data := update.Data.(type) {
			case ws.TickerUpdate:
				log.Printf("----Ticker of %s----", update.Pair)
				log.Printf("Ask: %s with %s", data.Ask.Price.String(), data.Ask.Volume.String())
				log.Printf("Bid: %s with %s", data.Bid.Price.String(), data.Bid.Volume.String())
				log.Printf("Open today: %s | Open last 24 hours: %s", data.Open.Today.String(), data.Open.Last24.String())
				dinle = convertertw.Converter(data.Ask.Price)
				//var askVlm = convertertw.Converter(data.Ask.Price)
				cst <- dinle
				askPrc = <-cst
			default:
			}
		}
	}

}
