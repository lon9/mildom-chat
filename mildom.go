package mildom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const nonopara = "fr=web`sfr=pc`devi=OS X 10.14.6 64-bit`la=ja`gid=%s`na=Japan`loc=Japan|Tokyo`clu=aws_japan`wh=1440*900`rtm=%s`ua=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36`aid=10000157`live_type=2`live_subtype=2`game_key=Apex_Legends`game_type=pc`host_official_type=official_game`isHomePage=false"

const serverURL = "https://im.mildom.com/?room_id=%d"
const liveInfoURL = "https://cloudac.mildom.com/nonolive/gappserv/live/enterstudio?user_id=%d&timestamp=%s&__guest_id=%s&__location=Japan%%7CTokyo&__country=Japan&__cluster=aws_japan&__platform=web&__la=ja&sfr=pc&accessToken="

func getLiveInfo(u string) (isLive bool, err error) {
	resp, err := http.Get(u)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	s, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(s))

	res := struct {
		Body struct {
			LiveMode int `json:"live_mode"`
		} `json:"body"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&res); err != nil {
		return
	}

	if res.Body.LiveMode != 0 {
		return true, nil
	}

	return
}

func getServerInfo(u string) (info map[string]string, err error) {

	info = make(map[string]string)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&info); err != nil {
		return
	}
	return
}

// Listen starts listening the room
func Listen(roomID int) (err error) {

	id := uuid.New()
	log.Println(id)

	guestID := "pc-gp-" + id.String()

	u := fmt.Sprintf(liveInfoURL, roomID, time.Now(), guestID)
	log.Println(u)
	ok, err := getLiveInfo(u)
	if err != nil {
		return
	}
	log.Println(ok)
	if !ok {
		log.Println("offline")
	} else {
		log.Println("online")
	}

	serverInfo, err := getServerInfo(fmt.Sprintf(serverURL, roomID))
	if err != nil {
		return
	}

	wsURL := "wss://" + serverInfo["wss_server"] + "?roomId=" + strconv.Itoa(roomID)
	log.Println(wsURL)
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}

	defer c.Close()

	done := make(chan interface{})

	go func() {
		var cnt int
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			res := make(map[string]interface{})
			if err := json.Unmarshal(msg, &res); err != nil {
				log.Println(err)
				return
			}
			log.Println(res["cmd"])
			switch res["cmd"].(string) {
			case "onChat":
				if cnt < 1000 {
					cnt++
					continue
				}
				log.Println(res["userName"], res["msg"])
			case "onLiveEnd":
				break

			}
		}
	}()

	initialMsg := NewInitialMsg(roomID, guestID)
	if err = c.WriteJSON(initialMsg); err != nil {
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// InitialMsg is message for initialization
type InitialMsg struct {
	AvatarDecoration int    `json:"avatarDecoration"`
	Cmd              string `json:"cmd"`
	EnterroomEffect  int    `json:"enterroomEffect"`
	GuestID          string `json:"guestId"`
	Level            int    `json:"level"`
	NobleClose       int    `json:"nobleClose"`
	NobleLevel       int    `json:"nobleLevel"`
	NobleSeatClose   int    `json:"nobleSeatClose"`
	Nonopara         string `json:"nonopara"`
	ReConnect        int    `json:"reConnect"`
	ReqID            int    `json:"reqId"`
	RoomID           int    `json:"roomId"`
	UserID           int    `json:"userId"`
	UserName         string `json:"userName"`
}

// NewInitialMsg is constructor for InitialMsg
func NewInitialMsg(roomID int, guestID string) *InitialMsg {
	return &InitialMsg{
		Cmd:      "enterRoom",
		RoomID:   roomID,
		GuestID:  guestID,
		Nonopara: fmt.Sprintf(nonopara, guestID, time.Now()),
		UserName: "test",
		Level:    1,
		ReqID:    1,
	}
}
