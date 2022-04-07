package main

import (
	"encoding/json"
	"fmt"
	convertertw "main/converterTw"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

var dinle2 = make(chan string)
var dinle = make(chan string)
var askPrc string
var upgrader = websocket.Upgrader{} // use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()
	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			//	log.Println("Error during message reading:", err)
			break
		}

		log.Printf("Alinan: %s", message)
		aa <- string(message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}
func routin(wg *sync.WaitGroup) {
	time.Sleep(2 * time.Second)
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
				var aba = convertertw.Converter(data.Ask.Price)
				//var askVlm = convertertw.Converter(data.Ask.Price)
				dinle <- aba

			default:
			}
		}
	}
	defer wg.Done()
}
func routin_2(wg *sync.WaitGroup) {
	time.Sleep(time.Second)
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", dea)
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
	wg.Done()
}

var aa = make(chan string)

func dea(w http.ResponseWriter, r *http.Request) {

	sss := <-dinle
	b, _ := json.Marshal(sss)
	w.Write(b)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go routin(&wg)
	go routin_2(&wg)
	fmt.Println("selam")
	wg.Wait()
}

/*
var addr = flag.String("addr", ":8000", "http service address")
func channells(yaz string) {
	done = make(chan interface{})    // Alici İşleyicisinin tamamlandiğini gösteren kanal
	interrupt = make(chan os.Signal) // Kesintisiz bir şekilde sonlandirmak için kesme sinyalini dinleyecek kanal

	signal.Notify(interrupt, os.Interrupt) // SIGINT için kesme kanalini bilgilendir

	socketUrl := "ws://localhost:8000" + "/socket"
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
				log.Println("Alici kanali kapatirken zaman aşimi. Çikiliyor…")
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
}*/
