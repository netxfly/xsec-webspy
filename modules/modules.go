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

package modules

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/urfave/cli"
	"github.com/sirupsen/logrus"

	"xsec-webspy/modules/arpspoof"
	"xsec-webspy/modules/logger"
	"xsec-webspy/modules/assembly"
	"xsec-webspy/modules/vars"

	"time"
	"strings"
	"fmt"
)

var (
	snapshotLen int32 = 1024
	promiscuous bool  = true
	err         error
	timeout     time.Duration = pcap.BlockForever
	handle      *pcap.Handle

	DebugMode  bool
	DeviceName = "eth0"
	filter     = ""
	Mode       = "local"


)

func Start(ctx *cli.Context) {
	if ctx.IsSet("device") {
		DeviceName = ctx.String("device")
	}

	if ctx.IsSet("mode") {
		Mode = ctx.String("mode")
	}

	if ctx.IsSet("host") {
		vars.HttpHost = ctx.String("host")
	}

	if ctx.IsSet("port") {
		vars.HttpPort = ctx.Int("port")
	}

	if ctx.IsSet("debug") {
		DebugMode = ctx.Bool("debug")
	}
	if DebugMode {
		logger.Log.Logger.Level = logrus.DebugLevel
	}

	if ctx.IsSet("length") {
		snapshotLen = int32(ctx.Int("len"))
	}
	// Open device
	handle, err = pcap.OpenLive(DeviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	if ctx.IsSet("filter") {
		filter = ctx.String("filter")
	}
	handle.SetBPFFilter(filter)

	go Serve(fmt.Sprintf("%v:%v", vars.HttpHost, vars.HttpPort))

	if strings.ToLower(Mode) == "local" {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		assembly.ProcessPackets(packetSource.Packets())
	} else {
		target := ""
		if ctx.IsSet("target") {
			target = ctx.String("target")
		}

		gateway := ""
		if ctx.IsSet("gateway") {
			gateway = ctx.String("gateway")
		}

		if target != "" && gateway != "" {
			go arpspoof.ArpSpoof(DeviceName, handle, target, gateway)

			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			assembly.ProcessPackets(packetSource.Packets())
		} else {
			logger.Log.Info("Need to provide target and gateway parameters")
		}
	}

}
