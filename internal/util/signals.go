package util

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GracefulShutDown is meant to be used in tandem with context, to catch signals, always use with a go-routine as it's blocking,
// unless you know what you're doing or, you've started your server in a go-routine, then this should be called without go keyword.
func GracefulShutDown(cancel context.CancelFunc, srv *http.Server) {
	defer cancel()
	c := make(chan os.Signal, 1)

	signal.Notify(
		c,
		os.Interrupt,
		os.Kill, // it doesn't catch kill command, still kept here, so that we don't forget this.
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// Block until a signal is received.
	sig := <-c
	//log.Errorf("Received signal: %+v", sig)
	log.Infof("Shutting down http server... Received signal: %v", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Http Server gracefully stopped")
}

// GracefulShutDownGrpc works in exactly same way as GracefulShutDown, but works for grpc servers
func GracefulShutDownGrpc(srv *grpc.Server) {
	c := make(chan os.Signal, 1)

	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// Block until a signal is received.
	sig := <-c
	//log.Errorf("Received signal: %+v", sig)
	log.Infof("Shutting down grpc server... Received signal: %v", sig)

	srv.GracefulStop()

	log.Println("Grpc Server gracefully stopped")
}
