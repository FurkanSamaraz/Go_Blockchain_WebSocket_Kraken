package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")
var askPrc string

func converter(con json.Number) (s string) {

	client := redis.NewClient(&redis.Options{
		Addr:     "89.19.7.50:6379",
		Password: "MWNkNjEzYzAxNzU5YzAyY2ZiZGY3ZjAz",
		DB:       0,
	})

	val, err := client.Get("btcoley:provider_fee").Result()
	val_2, err := client.Get("btcoley:exchange_fee").Result()
	//val_3, err := client.Get("btcoley:USD_TRY_with_fee").Result()

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("provider fee => ", val)
	//fmt.Println("exchange_fee => ", val_2)
	//fmt.Println("USD_TRY_with_fee => ", val_3)
	f, _ := con.Float64()
	f2, err := strconv.ParseFloat(val, 8)
	f1, err := strconv.ParseFloat(val_2, 8)
	c1 := f1 + 1
	c2 := f2 + 1
	aa := f * c1
	askPrice_3 := aa * c2
	s = fmt.Sprintf("%f", askPrice_3)
	fmt.Println("USD_TRY_with_fee => ", c1)
	fmt.Println("USD_TRY_with_fee => ", c2)
	return s
}

func main() {

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
			case []ws.Trade:
				log.Printf("----Trade of %s----", update.Pair)
				for i := range data {
					log.Printf("Price: %s", data[i].Price.String())
					log.Printf("Volume: %s", data[i].Volume.String())
					log.Printf("Time: %s", data[i].Time.String())
					log.Printf("Order type: %s", data[i].OrderType)
					log.Printf("Side: %s", data[i].Side)
					log.Printf("Misc: %s", data[i].Misc)
					fmt.Println("===================================================")

				}

			case ws.OrderBookUpdate:
				for {
					log.Printf("----Aksk of %s----", update.Pair)
					for _, ask := range data.Asks {

						askPrc = converter(ask.Price)
						var askVlm = converter(ask.Volume)
						//fmt.Println(askStr)
						log.Printf("%s | %s", askPrc, askVlm)
						channells(askPrc)
						channells(askVlm)
					}

					fmt.Println("===================================================")
					log.Printf("----Bids of %s----", update.Pair)
					for _, bid := range data.Bids {
						var bidPrc = converter(bid.Price)
						var bidVlm = converter(bid.Volume)
						log.Printf("%s | %s", bidPrc, bidVlm)
					}

				}
			default:
			}

		}

	}

}

func channells(yaz string) {
	done = make(chan interface{})    // Alici İşleyicisinin tamamlandiğini gösteren kanal
	interrupt = make(chan os.Signal) // Kesintisiz bir şekilde sonlandirmak için kesme sinyalini dinleyecek kanal

	signal.Notify(interrupt, os.Interrupt) // SIGINT için kesme kanalini bilgilendir

	socketUrl := "ws://localhost:8080" + "/socket"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Websocket Sunucusuna bağlanirken hata oluştu:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	// İstemci için ana döngümüz
	// İlgili paketlerimizi buraya gönderiyoruz
	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			//Her saniye bir yanki paketi gönder
			err := conn.WriteMessage(websocket.TextMessage, []byte(yaz))
			if err != nil {
				log.Println("Websocket'e yazarken hata:", err)
				return
			}

			//log.Println("SIGINT kesme sinyali alindi. Bekleyen tüm bağlantilari kapatma")

			// Websocket bağlantimizi kapatin
			erre := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if erre != nil {
				//	log.Println("Websocket'i kapatirken hata:", err)
				return
			}

			select {
			case <-done:
			//	log.Println("Alici Kanali Kapatildi! Çikiliyor….")
			case <-time.After(time.Duration(1) * time.Second):
				//	log.Println("Alici kanali kapatirken zaman aşimi. Çikiliyor…")
			}
			return
		}
	}
}

var done chan interface{}
var interrupt chan os.Signal

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			//log.Println("alma hatasi:", err)
			return
		}
		log.Printf("Gönderilen: %s\n", msg)
	}
}

/*package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/webdeveloppro/golang-websocket-client/pkg/client"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

var addr = flag.String("addr", ":8000", "http service address")
var askPrc string

func converter(con json.Number) (s string) {
	f, _ := con.Float64()

	askPrice_3 := f * 20
	s = fmt.Sprintf("%f", askPrice_3)
	return s
}

func main() {
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
			case []ws.Trade:
				log.Printf("----Trades of %s----", update.Pair)
				for i := range data {
					log.Printf("Price: %s", data[i].Price.String())
					log.Printf("Volume: %s", data[i].Volume.String())
					log.Printf("Time: %s", data[i].Time.String())
					log.Printf("Order type: %s", data[i].OrderType)
					log.Printf("Side: %s", data[i].Side)
					log.Printf("Misc: %s", data[i].Misc)
					fmt.Println("===================================================")
				}
			case ws.OrderBookUpdate:

				log.Printf("----Aksk of %s----", update.Pair)
				for _, ask := range data.Asks {
					askPrc = converter(ask.Price)
					var askVlm = converter(ask.Volume)
					//fmt.Println(askStr)
					log.Printf("%s | %s", askPrc, askVlm)

				}
				flag.Parse()

				client, err := client.NewWebSocketClient(*addr, "frontend")
				if err != nil {
					panic(err)
				}
				fmt.Println("Connecting")

				go func() {
					// write down data every 100 ms
					ticker := time.NewTicker(time.Millisecond * 1500)
					i := 0
					for range ticker.C {
						err := client.Write(askPrc)
						if err != nil {
							fmt.Printf("error: %v, writing error\n", err)
						}
						i++
					}
				}()

				// Close connection correctly on exit
				sigs := make(chan os.Signal, 1)

				// `signal.Notify` registers the given channel to
				// receive notifications of the specified signals.
				signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

				// The program will wait here until it gets the
				<-sigs
				client.Stop()
				fmt.Println("Goodbye")
				fmt.Println("===================================================")
				log.Printf("----Bids of %s----", update.Pair)
				for _, bid := range data.Bids {
					var bidPrc = converter(bid.Price)
					var bidVlm = converter(bid.Volume)
					log.Printf("%s | %s", bidPrc, bidVlm)
				}
			default:
			}
		}
	}
}
*/
