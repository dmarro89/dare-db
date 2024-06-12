// provides a go module to establish connection to dare-db
package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/go-while/nodare-db/logger"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
	"sync"
)

const DefaultAddr = "localhost:2420"
const DefaultAddrSSL = "localhost:2420"
const DefaultConnectTimeout = time.Duration(9 * time.Second)
const DefaultRequestTimeout = time.Duration(9 * time.Second)
const DefaultIdleCliTimeout = time.Duration(60 * time.Second)

type Options struct {
	Addr        string
	Mode        int
	SSL         bool
	SSLinsecure bool
	Auth        string
	Daemon      bool
	TestWorker  bool
	Stop        chan struct{}
}

type Client struct {
	logger     *ilog.LOG
	mux        sync.Mutex
	stop       chan struct{}
	addr       string
	url        string
	mode       int
	ssl        bool
	insecure   bool
	auth       string
	daemon     bool
	testWorker bool
	conn       net.Conn
	http       *http.Client

}

func NewClient(opts *Options) (*Client, error) {
	log.Printf("NewClient opts='%#v'", opts)
	client := &Client{
		logger:     ilog.NewLogger(ilog.GetEnvLOGLEVEL()),
		addr:       opts.Addr,
		mode:       opts.Mode,
		ssl:        opts.SSL,
		insecure:   opts.SSLinsecure,
		auth:       opts.Auth,
		daemon:     opts.Daemon,
		testWorker: opts.TestWorker,
		stop:       opts.Stop,
	}
	return client.Connect(client)
}

func (c *Client) Connect(client *Client) (*Client, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.conn != nil || c.http != nil {
		// conn is established, return no error
		c.logger.Info("connection already established!?")
		return client, nil
	}
	switch client.mode {
		case 1:
			c.Transport() // FIXME catch error!

		case 2:
			switch c.ssl {
				case true:
					if c.addr == "" {
						c.addr = DefaultAddrSSL
					}
					conf := &tls.Config{
						InsecureSkipVerify: c.insecure,
						MinVersion:         tls.VersionTLS12,
						CurvePreferences: []tls.CurveID{
							tls.CurveP521,
							tls.CurveP384,
							tls.CurveP256},
						PreferServerCipherSuites: true,
						CipherSuites: []uint16{
							tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
							//tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
							//tls.TLS_AES_256_GCM_SHA384,
							//tls.TLS_AES_128_GCM_SHA256,
						},
					}
					c.logger.Info("client connecting to tls://'%s'", c.addr)
					conn, err := tls.Dial("tcp", c.addr, conf)
					if err != nil {
						c.logger.Error("client.Connect tls.Dial err='%v'", err)
						return nil, err
					}
					c.conn = conn
				case false:
					if c.addr == "" {
						c.addr = DefaultAddr
					}
					c.logger.Info("client connecting to tcp://'%s'", c.addr)
					conn, err := net.Dial("tcp", c.addr)
					if err != nil {
						c.logger.Error("client.Connect net.Dial err='%v'", err)
						return nil, err
					}
					c.conn = conn
			} // end switch c.ssl
		default:
			c.logger.Error("client invalid mode=%d", c.mode)
	}
	c.logger.Info("client established c.conn='%v' c.http='%v' mode=%d", c.conn, c.http, c.mode)

	if c.testWorker {
		c.logger.Info("booting testWorker")
		c.worker(c.testWorker)
	}
	if c.daemon {
		go c.worker(false)
		return nil, nil
	}
	return client, nil
}

func (c *Client) Transport() {
	if c.url == "" {
		switch c.ssl {
		case true:
			c.url = "https://" + c.addr
		case false:
			c.url = "http://" + c.addr
		}
	}
	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		// We use ABSURDLY large keys, and should probably not.
		TLSHandshakeTimeout: 60 * time.Second,
	}
	c.http = &http.Client{
		Transport: t,
	}
	log.Printf("Transport c.http='%v' c.url='%s'", c.http, c.url)
}

func (c *Client) Get(key string) (string, error) {
	c.mux.Lock() // we lock so nobody else can use the connection at the same time
	defer c.mux.Unlock()
	resp, err := c.http.Get(c.url+"/get/" + key)
	if err != nil {
		c.logger.Error("c.http.Get err='%v'", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("c.http.Get respBody err='%v'", err)
		return "", err
	}
	c.logger.Debug("c.http.Get resp='%#v'", resp)
	return string(body), nil
}

func (c *Client) Set(key string, value string) (string, error) {
	c.mux.Lock() // we lock so nobody else can use the connection at the same time
	defer c.mux.Unlock()
	if c.http == nil {
		c.logger.Error("c.http.Set c.http == nil")
		return "", fmt.Errorf("set failed c.http is nil")
	}
	resp, err := http.Post(c.url+"/set", "application/json", bytes.NewBuffer([]byte(`{"`+key+`":"`+value+`"}`)))
	if err != nil {
		c.logger.Error("c.http.Set err='%v'", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("c.http.Set respBody err='%v'", err)
		return "", err
	}
	c.logger.Debug("c.http.Set resp='%#v'", resp)
	return string(body), nil
}

func (c *Client) Del(key string) (string, error) {
	c.mux.Lock() // we lock so nobody else can use the connection at the same time
	defer c.mux.Unlock()
	resp, err := c.http.Get(c.url+"/del/" + key)
	//resp, err := c.http.Get(c.url+"/del/" + key)
	if err != nil {
		c.logger.Error("c.http.Del err='%v'", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("c.http.Del respBody err='%v'", err)
		return "", err
	}
	c.logger.Debug("c.http.Del resp='%#v'", resp)
	return string(body), nil
}

func (c *Client) worker(testWorker bool) {
	defer c.logger.Info("worker left")
}
