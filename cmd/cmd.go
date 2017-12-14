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

package cmd

import (
	"github.com/urfave/cli"
	"xsec-webspy/modules"
)

var Start = cli.Command{
	Name:        "start",
	Usage:       "sniff local server",
	Description: "startup sniff on local server",
	Action:      modules.Start,
	Flags: []cli.Flag{
		stringFlag("mode,m", "local", "webspy running mode, local or arp"),
		stringFlag("device,i", "eth0", "device name"),
		stringFlag("host,H", "127.0.0.1", "web server listen address"),
		intFlag("port,p", 4000, "web server listen address"),
		boolFlag("debug, d", "debug mode"),
		stringFlag("target, t", "", "target ip address"),
		stringFlag("gateway, g", "", "gateway ip address"),
		stringFlag("filter,f", "", "setting filters"),
		intFlag("length,l", 1024, "setting snapshot Length"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
