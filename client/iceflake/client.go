package iceflake

import (
	"errors"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/istyle-inc/iceflake/pbdef"
)

// Client client struct for get unique id
type Client struct {
	listenType string
	addr       string
	timeout    time.Duration
}

// New returns new Client
func New(listenType, addr string) *Client {
	return &Client{
		listenType: listenType,
		addr:       addr,
		timeout:    5 * time.Second,
	}
}

// WithTimeout set timeout
func (c *Client) WithTimeout(t time.Duration) {
	c.timeout = t
}

// Get return IceFlake struct
func (c *Client) Get() (*pbdef.IceFlake, error) {
	conn, err := net.Dial(c.listenType, c.addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	err = conn.SetReadDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return nil, err
	}
	// snowflake type id has 64bit length
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, errors.New("iceflake server returned nothing")
	}

	var result pbdef.IceFlake
	if err := proto.Unmarshal(buf[:n], &result); err != nil {
		return nil, err
	}

	return &result, nil
}
