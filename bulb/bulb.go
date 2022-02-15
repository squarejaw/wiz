package bulb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/squarejaw/wiz/udp"
	"github.com/tidwall/gjson"
)

type Bulb struct {
	IP     *string `json:"ip,omitempty"`
	Mac    *string `json:"mac,omitempty"`
	Params *Params `json:"method,omitempty"`
}

type Params struct {
	SceneID   *int  `json:"sceneId,omitempty"`
	Speed     *int  `json:"speed,omitempty"`
	Dimming   *int  `json:"dimming,omitempty"`
	Temp      *int  `json:"temp,omitempty"`
	Red       *int  `json:"r,omitempty"`
	Green     *int  `json:"g,omitempty"`
	Blue      *int  `json:"b,omitempty"`
	ColdWhite *int  `json:"c,omitempty"`
	WarmWhite *int  `json:"w,omitempty"`
	State     *bool `json:"state,omitempty"`
}

type Message struct {
	Method *string `json:"method,omitempty"`
	Params *Params `json:"params,omitempty"`
}

func (bulb *Bulb) SetState(timeout int) (string, error) {
	udpSession, err := udp.NewSession(*bulb.IP, time.Duration(timeout)*time.Second)
	if err != nil {
		return "", err
	}
	defer udpSession.Close()
	method := "setState"
	msg := Message{
		Method: &method,
		Params: bulb.Params,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	_, err = udpSession.Write(b)
	if err != nil {
		return "", err
	}

	buf := make([]byte, udp.MAX_SAFE_PAYLOAD_SIZE)
	_, _, err = udpSession.Read(buf)
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	result := gjson.GetBytes(buf, "result")
	if result.Get("success").Bool() {
		return result.String(), nil
	} else {
		errorMessage := gjson.GetBytes(buf, "error.message").String()
		return "", fmt.Errorf(errorMessage)
	}
}
