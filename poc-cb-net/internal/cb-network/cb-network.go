package cbnet

import (
	"flag"
	"fmt"
	dataobjects "github.com/cloud-barista/cb-larva/poc-cb-net/internal/cb-network/data-objects"
	"github.com/cloud-barista/cb-larva/poc-cb-net/internal/ip-checker"
	cblog "github.com/cloud-barista/cb-log"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// I use TUN interface, so only plain IP packet, no ethernet header + mtu is set to 1300
const (
	// BUFFERSIZE represents a size of read buffer.
	BUFFERSIZE = 1500
	// MTU represents a maximum transmission unit.
	MTU = "1300"
	// IPv4 represents a version of IP address
	IPv4 = "IPv4"
	// IPv6 represents a version of IP address
	IPv6 = "IPv6"
)

// CBLogger represents a logger to show execution processes according to the logging level.
var CBLogger *logrus.Logger

func init() {
	// cblog is a global variable.
	configPath := filepath.Join("..", "..", "configs", "log_conf.yaml")
	CBLogger = cblog.GetLoggerWithConfigPath("cb-network", configPath)
}

// CBNetwork represents a network for the multi-cloud
type CBNetwork struct {
	CBNet             *water.Interface // Assigned cbnet0 IP from the server
	name              string           // InterfaceName of CBNet, e.g., cbnet0
	port              int              // Port used for tunneling
	MyPublicIP        string           // Inquired public IP of VM/Host
	myPrivateNetworks []string         // Inquired CIDR blocks of private network of VM/Host
	//listenConnection  *net.UDPConn                // Connection for encapsulation and decapsulation
	NetworkingRules dataobjects.NetworkingRules // Networking rule for CBNet and tunneling
	isRunning       bool

	NetworkInterfaces []dataobjects.NetworkInterface // To be Deprecated
}

// NewCBNetwork represents a constructor of CBNetwork
func NewCBNetwork(name string, port int) *CBNetwork {
	CBLogger.Debug("Start.........")

	temp := &CBNetwork{name: name, port: port}
	temp.isRunning = false
	temp.inquiryVMPublicIP()
	temp.UpdateNetworkInterfaceInfo() // To be deprecated and update "updateCIDRBlocksOfPrivateNetwork"
	temp.updateCIDRBlocksOfPrivateNetwork()

	CBLogger.Debug("End.........")
	return temp
}

func (cbnetwork *CBNetwork) inquiryVMPublicIP() {
	CBLogger.Debug("Start.........")

	resp, err := http.Get("https://ifconfig.co/")
	if err != nil {
		CBLogger.Panic(err)
	}

	// Perform error handling
	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			CBLogger.Fatal("can't close the response", errClose)
		}
	}()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		CBLogger.Panic(err)
	}
	CBLogger.Tracef("%s\n", string(data))

	cbnetwork.MyPublicIP = string(data[:len(data)-1]) // Remove last '\n'

	CBLogger.Debug("End.........")
}

func (cbnetwork *CBNetwork) updateCIDRBlocksOfPrivateNetwork() {
	CBLogger.Debug("Start.........")

	// Explore network interfaces
	for _, networkInterface := range cbnetwork.NetworkInterfaces {
		CBLogger.Trace(networkInterface)
		// Explore IPs
		for _, IP := range networkInterface.IPs {
			isPrivateIP := ipchkr.IsPrivateIP(net.ParseIP(IP.IPAddress))
			// Is private IP ?
			if isPrivateIP {
				if IP.Version == IPv4 { // Is IPv4 ?
					cbnetwork.myPrivateNetworks = append(cbnetwork.myPrivateNetworks, IP.CIDRBlock)
					CBLogger.Tracef("True v4 %s, %s\n", IP.IPAddress, IP.CIDRBlock)
				} else if IP.Version == IPv6 { // Is IPv6 ?
					CBLogger.Tracef("True v6 %s, %s\n", IP.IPAddress, IP.CIDRBlock)
				} else { // Unknown version
					CBLogger.Trace("!!! Unknown version !!!")
				}
			} else {
				CBLogger.Tracef("PublicIP %s, %s\n", IP.IPAddress, IP.CIDRBlock)
			}
		}
	}

	CBLogger.Debug("End.........")
}

//func (self CBNetworkAgent) GetCIDRBlocksOfPrivateNetwork() []string {
//	return self.myPrivateNetworks
//}

// GetVMNetworkInformation represents a function to get the network information of a VM.
func (cbnetwork CBNetwork) GetVMNetworkInformation() dataobjects.VMNetworkInformation {
	CBLogger.Debug("Start.........")

	temp := dataobjects.VMNetworkInformation{
		PublicIP:        cbnetwork.MyPublicIP,
		PrivateNetworks: cbnetwork.myPrivateNetworks,
	}
	CBLogger.Trace(temp)

	CBLogger.Debug("End.........")
	return temp
}

// SetNetworkingRules represents a function to set a networking rule
func (cbnetwork *CBNetwork) SetNetworkingRules(rules dataobjects.NetworkingRules) {
	CBLogger.Debug("Start.........")

	cbnetwork.NetworkingRules = rules

	CBLogger.Debug("End.........")
}

func (cbnetwork *CBNetwork) initCBNet() {
	CBLogger.Debug("Start.........")

	idx := cbnetwork.NetworkingRules.GetIndexOfPublicIP(cbnetwork.MyPublicIP)
	localNetwork := cbnetwork.NetworkingRules.CBNet[idx]

	localIP := flag.String("local", localNetwork, "Local tun interface IP/MASK like 192.168.3.3⁄24")
	if *localIP == "" {
		flag.Usage()
		CBLogger.Fatal("local ip is not specified")
	}

	iface, err := water.New(water.Config{
		DeviceType:             water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{Name: cbnetwork.name},
	})
	if nil != err {
		CBLogger.Fatal("Unable to allocate TUN interface:", err)
	}
	CBLogger.Info("Interface allocated:", iface.Name())

	cbnetwork.CBNet = iface
	CBLogger.Trace("=== *cbnetwork.CBNet: ", *cbnetwork.CBNet)
	CBLogger.Trace("=== cbnetwork.CBNet: ", cbnetwork.CBNet)

	// Set interface parameters
	cbnetwork.runIP("link", "set", "dev", cbnetwork.CBNet.Name(), "mtu", MTU)
	cbnetwork.runIP("addr", "add", *localIP, "dev", cbnetwork.CBNet.Name())
	cbnetwork.runIP("link", "set", "dev", cbnetwork.CBNet.Name(), "up")

	CBLogger.Debug("End.........")
}

func (cbnetwork *CBNetwork) runIP(args ...string) {
	CBLogger.Debug("Start.........")

	CBLogger.Trace(args)

	cmd := exec.Command("/sbin/ip", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		CBLogger.Fatal("Error running /sbin/ip:", err)
	}

	CBLogger.Debug("End.........")
}

// IsRunning represents a status of CBNetwork
func (cbnetwork CBNetwork) IsRunning() bool {
	CBLogger.Debug("Start.........")

	CBLogger.Debug("IsRunning? ", cbnetwork.isRunning)

	CBLogger.Debug("End.........")
	return cbnetwork.isRunning
}

// StartCBNetworking represents a function to start networking by networking rules
func (cbnetwork *CBNetwork) StartCBNetworking(channel chan bool) {
	CBLogger.Debug("Start.........")

	CBLogger.Info("Run CBNetworking between VMs")

	cbnetwork.initCBNet()
	channel <- true
	cbnetwork.isRunning = true

	CBLogger.Debug("End.........")
}

//func (cbnet *CBNetwork) RunDecapsulation(channel chan bool) {
//	CBLogger.Debug("Start.........")
//
//	CBLogger.Debug("Blocked till Networking Rule setup")
//	<-channel
//
//	CBLogger.Debug("Start decapsulation")
//	// Decapsulation
//
//	// Listen to local socket
//	// Create network address to listen
//	lstnAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", cbnet.port))
//	if nil != err {
//		CBLogger.Fatal("Unable to get UDP socket:", err)
//	}
//
//	// Create connection to network address
//	lstnConn, err := net.ListenUDP("udp", lstnAddr)
//	if nil != err {
//		CBLogger.Fatal("Unable to listen on UDP socket:", err)
//	}
//	defer lstnConn.Close()
//
//	buf := make([]byte, BUFFERSIZE)
//	for {
//		// ReadFromUDP acts like ReadFrom but returns a UDPAddr.
//		n, addr, err := lstnConn.ReadFromUDP(buf)
//		if err != nil {
//			CBLogger.Debug("Error in cbnet.listenConnection.ReadFromUDP(buf): ", err)
//		}
//
//		// Parse header
//		header, _ := ipv4.ParseHeader(buf[:n])
//		CBLogger.Debugf("Received %d bytes from %v: %+v\n", n, addr, header)
//
//		// It might be necessary to handle or route packets to the specific destination
//		// based on the NetworkingRules table
//		// To be determined.
//
//		// Write to TUN interface
//		nWrite, errWrite := cbnet.CBNet.Write(buf[:n])
//		if errWrite != nil || nWrite == 0 {
//			CBLogger.Debugf("Error(%d len): %s", nWrite, errWrite)
//		}
//	}
//}
//
//func (cbnet *CBNetwork) RunEncapsulation(channel chan bool) {
//	CBLogger.Debug("Start.........")
//
//	CBLogger.Debug("Blocked till Networking Rule setup")
//	<-channel
//
//	CBLogger.Debug("Start encapsulation")
//	packet := make([]byte, BUFFERSIZE)
//	for {
//		// Read packet from CBNet interface "cbnet0"
//		//fmt.Println("=== *cbnet.CBNet: ", *cbnet.CBNet)
//		//fmt.Println("=== cbnet.CBNet: ",cbnet.CBNet)
//		plen, err := cbnet.CBNet.Read(packet)
//		if err != nil {
//			CBLogger.Error("Error Read() in encapsulation:", err)
//		}
//
//		// Parse header
//		header, err := ipv4.ParseHeader(packet[:plen])
//		CBLogger.Tracef("Sending to remote: %+v (%+v)\n", header, err)
//
//		// Search and change destination (Public IP of target VM)
//		idx := cbnet.NetworkingRules.GetIndexOfCBNetIP(header.Dst.String())
//
//		var remoteIP string
//		if idx != -1 {
//			remoteIP = cbnet.NetworkingRules.PublicIP[idx]
//		}
//
//		// Resolve remote addr
//		remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", remoteIP, cbnet.port))
//		if nil != err {
//			CBLogger.Fatal("Unable to resolve remote addr:", err)
//		}
//
//		// Send packet
//		nWriteToUDP, errWriteToUDP := cbnet.listenConnection.WriteToUDP(packet[:plen], remoteAddr)
//		if errWriteToUDP != nil || nWriteToUDP == 0 {
//			CBLogger.Fatalf("Error(%d len): %s", nWriteToUDP, errWriteToUDP)
//		}
//	}
//}

// RunTunneling represents a function to be performing tunneling between hosts (e.g., VMs).
func (cbnetwork *CBNetwork) RunTunneling(channel chan bool) {
	CBLogger.Debug("Start.........")

	CBLogger.Debug("Blocked till Networking Rule setup")
	<-channel

	CBLogger.Debug("Start decapsulation")
	// Decapsulation

	// Listen to local socket
	// Create network address to listen
	lstnAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", cbnetwork.port))
	if nil != err {
		CBLogger.Fatal("Unable to get UDP socket:", err)
	}

	// Create connection to network address
	lstnConn, err := net.ListenUDP("udp", lstnAddr)
	if nil != err {
		CBLogger.Fatal("Unable to listen on UDP socket:", err)
	}

	// Perform error handling
	defer func() {
		errClose := lstnConn.Close()
		if errClose != nil {
			CBLogger.Fatal("can't close the listen connection", errClose)
		}
	}()

	go func() {
		buf := make([]byte, BUFFERSIZE)
		for {
			// ReadFromUDP acts like ReadFrom but returns a UDPAddr.
			n, _, err := lstnConn.ReadFromUDP(buf)
			if err != nil {
				CBLogger.Error("Error in cbnetwork.listenConnection.ReadFromUDP(buf): ", err)
			}

			// Parse header
			header, _ := ipv4.ParseHeader(buf[:n])
			CBLogger.Tracef("Header received: %+v\n", header)
			//fmt.Printf("Received %d bytes from %v: %+v\n", n, addr, header)

			// It might be necessary to handle or route packets to the specific destination
			// based on the NetworkingRules table
			// To be determined.

			// Write to TUN interface
			nWrite, errWrite := cbnetwork.CBNet.Write(buf[:n])
			if errWrite != nil || nWrite == 0 {
				CBLogger.Errorf("Error(%d len): %s", nWrite, errWrite)
			}
		}
	}()

	CBLogger.Debug("Start encapsulation")
	// Encapsulation
	packet := make([]byte, BUFFERSIZE)
	for {
		// Read packet from CBNet interface "cbnet0"
		//fmt.Println("=== *cbnetwork.CBNet: ", *cbnetwork.CBNet)
		//fmt.Println("=== cbnetwork.CBNet: ",cbnetwork.CBNet)
		plen, err := cbnetwork.CBNet.Read(packet)
		if err != nil {
			CBLogger.Error("Error Read() in encapsulation:", err)
		}

		// Parse header
		header, _ := ipv4.ParseHeader(packet[:plen])
		CBLogger.Tracef("Sending to remote: %+v (%+v)\n", header, err)

		// Search and change destination (Public IP of target VM)
		idx := cbnetwork.NetworkingRules.GetIndexOfCBNetIP(header.Dst.String())

		var remoteIP string
		if idx != -1 {
			remoteIP = cbnetwork.NetworkingRules.PublicIP[idx]
		}

		// Resolve remote addr
		remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", remoteIP, cbnetwork.port))
		if nil != err {
			CBLogger.Fatal("Unable to resolve remote addr:", err)
		}

		// Send packet
		nWriteToUDP, errWriteToUDP := lstnConn.WriteToUDP(packet[:plen], remoteAddr)
		if errWriteToUDP != nil || nWriteToUDP == 0 {
			CBLogger.Errorf("Error(%d len): %s", nWriteToUDP, errWriteToUDP)
		}
	}
}

// UpdateNetworkInterfaceInfo represents a function to update the network interface information.
// To be deprecated
// Define a function to get network interfaces in a physical or virtual machine
func (cbnetwork *CBNetwork) UpdateNetworkInterfaceInfo() {
	CBLogger.Debug("Start.........")

	// Get network interfaces
	ifaces, _ := net.Interfaces()

	// Recursively get network interface information
	for _, iface := range ifaces {
		// Print a network interface name
		CBLogger.Debug("Interface name:", iface.Name)

		// Declare a NetworkInterface variable
		var networkInterface dataobjects.NetworkInterface

		// Assign Interface Interface Name
		networkInterface.Name = iface.Name

		// Get addresses
		addrs, _ := iface.Addrs()

		// Recursively get IP address
		for _, addr := range addrs {
			addrStr := addr.String()

			// Get IP Address and CIDRBlock ID
			ipAddr, networkID, err := net.ParseCIDR(addrStr)
			if err != nil {
				CBLogger.Fatal(err)
			}

			// Get version of IP (e.g., IPv4 or IPv6)
			var version string

			if ipAddr.To4() != nil {
				version = IPv4
			} else if ipAddr.To16() != nil {
				version = IPv6
			} else {
				version = "Unknown"
				CBLogger.Tracef("Unknown version (IPAddr: %s)\n", ipAddr.String())
			}

			// Print version, IPAddress, CIDRBlock ID
			//fmt.Println("	 Version: ", version)
			//fmt.Println("	 IPAddr: ", ipAddr)
			//fmt.Println("	 CIDRBlock: ", networkID)

			// Create IP data object
			ip := dataobjects.IP{
				Version:   version,
				IPAddress: ipAddr.String(),
				CIDRBlock: networkID.String(),
			}

			// AppendRule the IP data object to slice
			networkInterface.IPs = append(networkInterface.IPs, ip)
		}
		cbnetwork.NetworkInterfaces = append(cbnetwork.NetworkInterfaces, networkInterface)
	}
	CBLogger.Debug("End.........")
}

//func (cbnet CBNetwork) GetNetworkInterfaces() []dataobjects.NetworkInterface {
//	CBLogger.Debug("Start.........")
//
//	CBLogger.Trace("cbnet.NetworkInterfaces")
//	CBLogger.Trace(cbnet.NetworkInterfaces)
//
//	CBLogger.Debug("End.........")
//	return cbnet.NetworkInterfaces
//}