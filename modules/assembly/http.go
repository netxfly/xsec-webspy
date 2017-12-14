/*

Copyright (c) 2017 xsec.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package assembly

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/google/gopacket/layers"

	"xsec-webspy/models"
	"xsec-webspy/modules/vars"
	"xsec-webspy/modules/logger"

	"bufio"
	"io"
	"time"
	"net/http"
	"strings"
	"fmt"
)

type httpStreamFactory struct{}

type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.run()
	return &hstream.r
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			return
		} else if err == nil {
			defer req.Body.Close()

			clientIp, dstIp := SplitNet2Ips(h.net)
			srcPort, dstPort := Transport2Ports(h.transport)

			httpReq := models.NewHttpReq(req, clientIp, dstIp, dstPort)

			// send to sever
			go func(addr string, req *models.HttpReq, ) {
				reqInfo := fmt.Sprintf("%v:%v -> %v(%v:%v), %v, %v, %v, %v", httpReq.Client, srcPort, httpReq.Host, httpReq.Ip,
					httpReq.Port, httpReq.Method, httpReq.URL, httpReq.Header, httpReq.ReqParameters)
				logger.Log.Warnf(reqInfo)

				SendHTML(reqInfo)
				//if !CheckSelfHtml(addr, req) {
				//	SendHTML(req)
				//}
			}(vars.HttpHost, httpReq)
		}
	}
}

func SplitNet2Ips(net gopacket.Flow) (client, host string) {
	ips := strings.Split(net.String(), "->")
	if len(ips) > 1 {
		client = ips[0]
		host = ips[1]
	}
	return client, host
}

func Transport2Ports(transport gopacket.Flow) (src, dst string) {
	ports := strings.Split(transport.String(), "->")
	if len(ports) > 1 {
		src = ports[0]
		dst = ports[1]
	}
	return src, dst
}

func CheckSelfHtml(host string, req *models.HttpReq) (ret bool) {
	if host == strings.Split(req.Host, ":")[0] {
		ret = true
	}
	return ret
}

func ProcessPackets(packets chan gopacket.Packet) {
	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:

			if packet == nil {
				return
			}

			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}

func SendHTML(reqInfo string) {
	vars.Data.Put(reqInfo)
}
