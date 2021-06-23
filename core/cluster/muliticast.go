/*
Copyright Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"github.com/huminghe/infini-framework/core/config"
	"github.com/huminghe/infini-framework/core/util"
	"log"
	"net"
)

const (
	maxDataSize = 4096
)

//send a Broadcast message to network to discovery the cluster
func Broadcast(config config.NetworkConfig, req *Request) {

	addr, err := net.ResolveUDPAddr("udp", config.GetBindingAddr())
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	payload := util.ToJSONBytes(req)

	c.Write(payload)

}

func ServeMulticastDiscovery(config config.NetworkConfig, h func(*net.UDPAddr, int, []byte), signal chan bool) {

	addr, err := net.ResolveUDPAddr("udp", config.GetBindingAddr())
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	l.SetReadBuffer(maxDataSize)

	signal <- true

	for {
		b := make([]byte, maxDataSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(src, n, b)
	}

}
