package probe

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ProberTestSuite struct {
	suite.Suite
	prober *Prober
}

func createTCPProbeMock() *Probe {
	return &Probe{
		Probe: &v1.Probe{
			FailureThreshold: 1,
			PeriodSeconds:    1,
			ProbeHandler: v1.ProbeHandler{
				TCPSocket: &v1.TCPSocketAction{
					Port: intstr.FromInt(8080),
				},
			},
		},
	}
}

func createHTTPProbeMock(port int) *Probe {
	return &Probe{
		Probe: &v1.Probe{
			FailureThreshold: 1,
			PeriodSeconds:    1,
			ProbeHandler: v1.ProbeHandler{
				HTTPGet: &v1.HTTPGetAction{
					Host: "localhost",
					Port: intstr.FromInt(port),
				},
			},
		},
	}
}

func runTCPServer(wg *sync.WaitGroup, pingCnt *atomic.Int32) {
	defer wg.Done()

	port := 8080
	// Start listening on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error listening on port %v: %v", port, err)
		return
	}
	defer listener.Close()

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			return
		default:
			// Wait for a connection
			_, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			pingCnt.Add(1)
		}
	}
}

func (t *ProberTestSuite) TestTCPProbe() {
	wg := &sync.WaitGroup{}
	pingCnt := &atomic.Int32{}
	wg.Add(1)
	go runTCPServer(wg, pingCnt)
	// Give time for preemption
	time.Sleep(time.Millisecond)
	probe := createTCPProbeMock()
	t.prober.RunProbe(probe)
	wg.Wait()

	t.GreaterOrEqual(pingCnt.Load(), int32(10))
}

type HTTPHandler struct {
	pingCnt *atomic.Int32
}

func NewHTTPHandler(cnt *atomic.Int32) *HTTPHandler {
	return &HTTPHandler{
		pingCnt: cnt,
	}
}

func (h *HTTPHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
	h.pingCnt.Add(1)
}

func createHTTPServer(wg *sync.WaitGroup, port int, pingCnt *atomic.Int32, signal chan<- struct{}) error {
	defer wg.Done()
	httpServer := &http.Server{
		Handler:           NewHTTPHandler(pingCnt),
		Addr:              fmt.Sprintf("localhost:%d", port),
		ReadHeaderTimeout: 60 * time.Second,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("unable to listen: %v", err)
	}
	signal <- struct{}{}
	go func() {
		time.Sleep(10 * time.Second)
		if err := httpServer.Shutdown(context.Background()); err != nil {
			fmt.Printf("error while trying to shutdown the HTTP server: %v", err)
		}
	}()
	return httpServer.Serve(lis)
}

func (t *ProberTestSuite) TestHTTPProbe() {
	port := 8081
	pingCnt := &atomic.Int32{}
	signal := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go createHTTPServer(wg, port, pingCnt, signal)
	<-signal
	t.prober.RunProbe(createHTTPProbeMock(port))
	wg.Wait()

	t.GreaterOrEqual(pingCnt.Load(), int32(10))
}

func TestExampleTestSuite(t *testing.T) {
	pSuite := &ProberTestSuite{
		prober: NewProber(),
	}
	suite.Run(t, pSuite)
}
