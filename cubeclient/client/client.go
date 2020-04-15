package client

import (
	"errors"
	"github.com/bahadrix/cardinalitycube/server/cubeserver"
	"github.com/zeromq/goczmq"
	"strings"
)

type Client struct {
	dealer *goczmq.Sock
	endpoint string
}

func NewClient(endpoint string) (*Client, error) {
	client := &Client{
		dealer:   nil,
		endpoint: endpoint,
	}
	return client, client.acquireDealer()
}

func (c *Client) acquireDealer() (err error) {
	if c.dealer == nil {
		c.dealer, err = goczmq.NewDealer(c.endpoint)
	}
	return
}

func (c *Client) Execute(command string, args ...string) (reply string, err error) {
	err = c.acquireDealer()
	if err != nil {
		return
	}
	msg := strings.Join(append([]string{command}, args...), " ")
	err = c.dealer.SendFrame([]byte(msg), goczmq.FlagNone)
	if err != nil {
		return
	}

	response, err := c.dealer.RecvMessage()
	if err != nil {
		return
	}
	responseBytes := response[0]
	responseSize := len(responseBytes)

	if responseSize < 1 {
		err = errors.New("unexpected null reply")
		return
	}

	state := responseBytes[0]
	if responseSize > 1 {
		reply = string(responseBytes[1:])
	}


	if state != cubeserver.MsgPrefixSuccess {
		err = errors.New(reply)
		reply = ""
	}

	return


}

func (c *Client) Close() {
	if c.dealer != nil {
		c.dealer.Destroy()
		c.dealer = nil
	}
}