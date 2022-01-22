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
	IP     string
	Params Params
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
	Method string `json:"method"`
	Params Params `json:"params"`
}

func (bulb *Bulb) SetState(timeout int) (string, error) {
	udpSession, err := udp.NewSession(bulb.IP+":38899", time.Duration(timeout)*time.Second)
	if err != nil {
		return "", err
	}
	defer udpSession.Close()

	msg := Message{
		Method: "setState",
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

	buf := make([]byte, 1024)
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
