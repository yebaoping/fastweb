package fastweb

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type options struct {
	addr              string
	listener          net.Listener
	certFile, keyFile string
	certData, keyData []byte
}

type Option func(o *options)

func Addr(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

func Listener(ln net.Listener) Option {
	return func(o *options) {
		o.listener = ln
	}
}

func TLS(certFile, keyFile string) Option {
	return func(o *options) {
		o.certFile = certFile
		o.keyFile = keyFile
	}
}

func TLSEmbed(certData, keyData []byte) Option {
	return func(o *options) {
		o.certData = certData
		o.keyData = keyData
	}
}

func (a *App) Run(opts ...Option) error {
	options := new(options)

	for _, option := range opts {
		option(options)
	}

	if options.listener == nil {
		if options.addr != "" {
			ln, err := net.Listen("tcp4", options.addr)
			if err != nil {
				return err
			}

			options.listener = ln
		} else {
			return errors.New("address is empty")
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		<-signals

		a.log.Debug("server stopping...")
		if err := a.Shutdown(); err != nil {
			a.log.Errorf("shutdown failed,error:%+v", err)
		}
	}()

	if options.certFile != "" && options.keyFile != "" {
		return a.ServeTLS(options.listener, options.certFile, options.keyFile)
	} else if len(options.certData) > 0 && len(options.keyData) > 0 {
		return a.ServeTLSEmbed(options.listener, options.certData, options.keyData)
	} else {
		return a.Serve(options.listener)
	}
}
