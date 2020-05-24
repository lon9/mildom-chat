package mildom

import (
	"fmt"
	"time"
)

const nonopara = "fr=web`sfr=pc`devi=OS X 10.14.6 64-bit`la=ja`gid=%s`na=Japan`loc=Japan|Tokyo`clu=aws_japan`wh=1440*900`rtm=%s`ua=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36`aid=10000157`live_type=2`live_subtype=2`game_key=Apex_Legends`game_type=pc`host_official_type=official_game`isHomePage=false"
const serverURL = "https://im.mildom.com/?room_id=%d"

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

// ChatMessage is chat object
type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}
