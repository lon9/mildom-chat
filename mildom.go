package mildom

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const liveInfoURL = "https://cloudac.mildom.com/nonolive/gappserv/live/enterstudio"

func getLiveInfo(u string) (isLive bool, err error) {
	resp, err := http.Get(u)
	if err != nil {
		return
	}

	defer resp.Body.Close()

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

// GetListener starts listening the room
func GetListener(roomID int) (chan *ChatMessage, error) {
	id := uuid.New()

	guestID := "pc-gp-" + id.String()

	v := url.Values{}
	v.Add("user_id", strconv.Itoa(roomID))
	v.Add("timestamp", time.Now().String())
	v.Add("__guest_id", guestID)
	v.Add("__location", "Japan|Tokyo")
	v.Add("__country", "Japan")
	v.Add("__cluster", "aws_japan")
	v.Add("__platform", "web")
	v.Add("__la", "ja")
	v.Add("__sfr", "pc")

	ok, err := getLiveInfo(liveInfoURL + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	if !ok {
		log.Println("offline")
		return nil, errors.New("offline")
	}

	serverInfo, err := getServerInfo(fmt.Sprintf(serverURL, roomID))
	if err != nil {
		return nil, err
	}

	wsURL := "wss://" + serverInfo["wss_server"] + "?roomId=" + strconv.Itoa(roomID)
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	ch := make(chan *ChatMessage)
	go func() {
		var cnt int
		defer close(ch)
		defer c.Close()
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			res := make(map[string]interface{})
			if err := json.Unmarshal(msg, &res); err != nil {
				log.Println(err)
				break
			}
			switch res["cmd"].(string) {
			case "onChat":
				if cnt < 1000 {
					cnt++
					continue
				}
				ch <- &ChatMessage{
					Username: res["userName"].(string),
					Text:     res["msg"].(string),
				}
			case "onLiveEnd":
				// end live
				if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
					log.Println("write close:", err)
					break
				}
				break
			}
		}
	}()

	initialMsg := NewInitialMsg(roomID, guestID)
	if err = c.WriteJSON(initialMsg); err != nil {
		return nil, err
	}
	return ch, nil
}
