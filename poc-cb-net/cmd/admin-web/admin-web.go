package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cbnet "github.com/cloud-barista/cb-larva/poc-cb-net/internal/cb-network"
	dataobjects "github.com/cloud-barista/cb-larva/poc-cb-net/internal/cb-network/data-objects"
	etcdkey "github.com/cloud-barista/cb-larva/poc-cb-net/internal/etcd-key"
	file "github.com/cloud-barista/cb-larva/poc-cb-net/internal/file"
	cblog "github.com/cloud-barista/cb-log"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var dscp *cbnet.DynamicSubnetConfigurator

// CBLogger represents a logger to show execution processes according to the logging level.
var CBLogger *logrus.Logger
var config dataobjects.Config

func init() {
	fmt.Println("Start......... init() of controller.go")
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	fmt.Printf("exePath: %v\n", exePath)

	// Load cb-log config from the current directory (usually for the production)
	logConfPath := filepath.Join(exePath, "configs", "log_conf.yaml")
	fmt.Printf("logConfPath: %v\n", logConfPath)
	if !file.Exists(logConfPath) {
		// Load cb-log config from the project directory (usually for development)
		path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
		if err != nil {
			panic(err)
		}
		projectPath := strings.TrimSpace(string(path))
		logConfPath = filepath.Join(projectPath, "poc-cb-net", "configs", "log_conf.yaml")
	}
	CBLogger = cblog.GetLoggerWithConfigPath("cb-network", logConfPath)
	CBLogger.Debugf("Load %v", logConfPath)

	// Load cb-network config from the current directory (usually for the production)
	configPath := filepath.Join(exePath, "configs", "config.yaml")
	fmt.Printf("configPath: %v\n", configPath)
	if !file.Exists(configPath) {
		// Load cb-network config from the project directory (usually for the development)
		path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
		if err != nil {
			panic(err)
		}
		projectPath := strings.TrimSpace(string(path))
		configPath = filepath.Join(projectPath, "poc-cb-net", "configs", "config.yaml")
	}
	config, _ = dataobjects.LoadConfig(configPath)
	CBLogger.Debugf("Load %v", configPath)
	fmt.Println("End......... init() of controller.go")
}

var (
	upgrader = websocket.Upgrader{}
)

var connectionPool = struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}{
	connections: make(map[*websocket.Conn]struct{}),
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// WebsocketHandler represents a handler to watch and send networking rules to AdminWeb frontend.
func WebsocketHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	connectionPool.Lock()
	connectionPool.connections[ws] = struct{}{}

	defer func(connection *websocket.Conn) {
		connectionPool.Lock()
		delete(connectionPool.connections, connection)
		connectionPool.Unlock()
	}(ws)

	connectionPool.Unlock()

	// etcd Section
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.ETCD.Endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		CBLogger.Fatal(err)
	}

	defer func() {
		errClose := etcdClient.Close()
		if errClose != nil {
			CBLogger.Fatal("Can't close the etcd client", errClose)
		}
	}()

	CBLogger.Infoln("The etcdClient is connected.")

	// Get the existing both the networking rule and the configuration information of the CLADNet
	errInitData := getExistingNetworkInfo(ws, etcdClient)
	if errInitData != nil {
		CBLogger.Errorf("getExistingNetworkInfo() error: %v", errInitData)
	}

	for {
		// Read
		_, msgRead, err := ws.ReadMessage()
		if err != nil {
			CBLogger.Error(err)
			return err
		}
		CBLogger.Tracef("Message Read: %s", msgRead)

		var message dataobjects.WebsocketMessageFrame
		errUnmarshalDataFrame := json.Unmarshal(msgRead, &message)
		if errUnmarshalDataFrame != nil {
			CBLogger.Error(errUnmarshalDataFrame)
		}

		switch message.Type {
		case "configuration-information":
			configurationInformationHandler(ws, etcdClient, []byte(message.Text))

		case "test-specification":
			testSpecificationHandler(ws, etcdClient, []byte(message.Text))

		default:

		}

	}
}

func configurationInformationHandler(ws *websocket.Conn, etcdClient *clientv3.Client, responseText []byte) {
	CBLogger.Debug("Start.........")
	// Unmarshal the configuration information of Cloud Adaptive Network (CLADNet)
	// :IPv4 CIDR block, Description
	var cladNetConfInfo dataobjects.CLADNetConfigurationInformation
	errUnmarshal := json.Unmarshal(responseText, &cladNetConfInfo)
	if errUnmarshal != nil {
		CBLogger.Error(errUnmarshal)
	}

	CBLogger.Tracef("Configuration information: %v", cladNetConfInfo)

	// Generate a unique CLADNet ID by the xid package
	guid := xid.New()
	CBLogger.Tracef("A unique CLADNet ID: %v", guid)
	cladNetConfInfo.CLADNetID = guid.String()

	// Currently assign the 1st IP address for Gateway IP (Not used till now)
	ipv4Address, _, errParseCIDR := net.ParseCIDR(cladNetConfInfo.CIDRBlock)
	if errParseCIDR != nil {
		CBLogger.Fatal(errParseCIDR)
	}

	CBLogger.Tracef("IPv4Address: %v", ipv4Address)
	ip := ipv4Address.To4()
	gatewayIP := incrementIP(ip, 1)
	cladNetConfInfo.GatewayIP = gatewayIP.String()
	CBLogger.Tracef("GatewayIP: %v", cladNetConfInfo.GatewayIP)

	// Put the configuration information of the CLADNet to the etcd
	keyConfigurationInformationOfCLADNet := fmt.Sprint(etcdkey.ConfigurationInformation + "/" + cladNetConfInfo.CLADNetID)
	strCLADNetConfInfo, _ := json.Marshal(cladNetConfInfo)
	_, err := etcdClient.Put(context.Background(), keyConfigurationInformationOfCLADNet, string(strCLADNetConfInfo))
	if err != nil {
		CBLogger.Fatal(err)
	}

	// Get the configuration information of the CLADNet
	respMultiConfInfo, errMultiConfInfo := etcdClient.Get(context.Background(), etcdkey.ConfigurationInformation, clientv3.WithPrefix())
	if errMultiConfInfo != nil {
		CBLogger.Fatal(errMultiConfInfo)
	}

	var CLADNetConfigurationInformationList []string
	for _, confInfo := range respMultiConfInfo.Kvs {
		CLADNetConfigurationInformationList = append(CLADNetConfigurationInformationList, string(confInfo.Value))
	}

	CBLogger.Tracef("CLADNetConfigurationInformationList: %v", CLADNetConfigurationInformationList)

	// Build response JSON
	var buf bytes.Buffer
	text := strings.Join(CLADNetConfigurationInformationList, ",")
	buf.WriteString("[")
	buf.WriteString(text)
	buf.WriteString("]")

	// Response to the front-end
	errResp := sendResponseText(ws, "CLADNetList", buf.String())
	if errResp != nil {
		CBLogger.Error(errResp)
	}
	CBLogger.Debug("End.........")
}

func testSpecificationHandler(ws *websocket.Conn, etcdClient *clientv3.Client, responseText []byte) {
	CBLogger.Debug("Start.........")
	var testSpecification dataobjects.TestSpecification
	errUnmarshalEvalSpec := json.Unmarshal(responseText, &testSpecification)
	if errUnmarshalEvalSpec != nil {
		CBLogger.Error(errUnmarshalEvalSpec)
	}

	CBLogger.Tracef("Evaluation specification: %v", testSpecification)

	// Put the evaluation specification of the CLADNet to the etcd
	keyStatusTestSpecificationOfCLADNet := fmt.Sprint(etcdkey.StatusTestSpecification + "/" + testSpecification.CLADNetID)
	strStatusTestSpecification, _ := json.Marshal(testSpecification)
	//spec := message.Text
	_, err := etcdClient.Put(context.Background(), keyStatusTestSpecificationOfCLADNet, string(strStatusTestSpecification))
	if err != nil {
		CBLogger.Error(err)
	}
	CBLogger.Debug("End.........")
}

func getExistingNetworkInfo(ws *websocket.Conn, etcdClient *clientv3.Client) error {

	// Get the networking rule
	CBLogger.Debugf("Get - %v", etcdkey.NetworkingRule)
	resp, etcdErr := etcdClient.Get(context.Background(), etcdkey.NetworkingRule, clientv3.WithPrefix())
	CBLogger.Tracef("etcdResp: %v", resp)
	if etcdErr != nil {
		CBLogger.Error(etcdErr)
	}

	if len(resp.Kvs) != 0 {
		for _, kv := range resp.Kvs {
			CBLogger.Tracef("CLADNet ID: %v", kv.Key)
			CBLogger.Tracef("The networking rule of the CLADNet: %v", kv.Value)
			CBLogger.Debug("Send a networking rule of CLADNet to AdminWeb frontend")

			// Send the networking rule to the front-end
			errResp := sendResponseText(ws, "NetworkingRule", string(kv.Value))
			if errResp != nil {
				CBLogger.Error(errResp)
			}
		}
	} else {
		CBLogger.Debug("No networking rule of CLADNet exists")
	}

	// Get the configuration information of the CLADNet
	CBLogger.Debugf("Get - %v", etcdkey.ConfigurationInformation)
	respMultiConfInfo, err := etcdClient.Get(context.Background(), etcdkey.ConfigurationInformation, clientv3.WithPrefix())
	if err != nil {
		CBLogger.Fatal(err)
		return err
	}

	if len(respMultiConfInfo.Kvs) != 0 {
		var CLADNetConfigurationInformationList []string
		for _, confInfo := range respMultiConfInfo.Kvs {
			CLADNetConfigurationInformationList = append(CLADNetConfigurationInformationList, string(confInfo.Value))
		}

		CBLogger.Tracef("CLADNetConfigurationInformationList: %v", CLADNetConfigurationInformationList)

		// Build response JSON
		var buf bytes.Buffer
		text := strings.Join(CLADNetConfigurationInformationList, ",")
		buf.WriteString("[")
		buf.WriteString(text)
		buf.WriteString("]")

		// Response to the front-end
		errResp := sendResponseText(ws, "CLADNetList", buf.String())
		if errResp != nil {
			CBLogger.Error(errResp)
			return errResp
		}
	}
	return nil
}

func incrementIP(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}

func buildResponseBytes(responseType string, responseText string) []byte {
	CBLogger.Debug("Start.........")
	var response dataobjects.WebsocketMessageFrame
	response.Type = responseType
	response.Text = responseText

	CBLogger.Tracef("ResponseStr: %v", response)
	responseBytes, _ := json.Marshal(response)
	CBLogger.Debug("End.........")
	return responseBytes
}

func sendResponseText(ws *websocket.Conn, responseType string, responseText string) error {
	CBLogger.Debug("Start.........")
	responseBytes := buildResponseBytes(responseType, responseText)

	// Response to the front-end
	errWriteJSON := ws.WriteMessage(websocket.TextMessage, responseBytes)
	if errWriteJSON != nil {
		CBLogger.Error(errWriteJSON)
		return errWriteJSON
	}
	CBLogger.Debug("End.........")
	return nil
}

func sendMessageToAllPool(message []byte) error {
	connectionPool.RLock()
	defer connectionPool.RUnlock()
	for connection := range connectionPool.connections {
		err := connection.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunEchoServer represents a function to run echo server.
func RunEchoServer(wg *sync.WaitGroup, config dataobjects.Config) {
	defer wg.Done()

	// Set web assets path to the current directory (usually for the production)
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	CBLogger.Tracef("exePath: %v", exePath)
	webPath := filepath.Join(exePath, "web")

	indexPath := filepath.Join(webPath, "public", "index.html")
	CBLogger.Tracef("indexPath: %v", indexPath)
	if !file.Exists(indexPath) {
		// Set web assets path to the project directory (usually for the development)
		path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
		if err != nil {
			panic(err)
		}
		projectPath := strings.TrimSpace(string(path))
		webPath = filepath.Join(projectPath, "poc-cb-net", "web")
	}

	CBLogger.Debug("Start.........")
	e := echo.New()

	e.Static("/", webPath+"/assets")
	e.Static("/js", webPath+"/assets/js")
	e.Static("/css", webPath+"/assets/css")
	e.Static("/introspect", webPath+"/assets/introspect")

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob(webPath + "/public/*.html")),
	}
	e.Renderer = renderer

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"websocket_host": "http://" + config.AdminWeb.Host + ":" + config.AdminWeb.Port,
		})
	})

	// Render
	e.GET("/ws", WebsocketHandler)

	e.Logger.Fatal(e.Start(":" + config.AdminWeb.Port))
	CBLogger.Debug("End.........")
}

func watchNetworkingRule(wg *sync.WaitGroup, etcdClient *clientv3.Client) {
	defer wg.Done()

	// Watch "/registry/cloud-adaptive-network/networking-rule"
	CBLogger.Debugf("Start to watch \"%v\"", etcdkey.NetworkingRule)
	watchChan1 := etcdClient.Watch(context.Background(), etcdkey.NetworkingRule, clientv3.WithPrefix())
	for watchResponse := range watchChan1 {
		for _, event := range watchResponse.Events {
			CBLogger.Tracef("Watch - %s %q : %q", event.Type, event.Kv.Key, event.Kv.Value)
			slicedKeys := strings.Split(string(event.Kv.Key), "/")
			parsedHostID := slicedKeys[len(slicedKeys)-1]
			CBLogger.Tracef("ParsedHostID: %v", parsedHostID)

			networkingRule := event.Kv.Value
			CBLogger.Tracef("A networking rule of CLADNet: %v", networkingRule)

			// Build the response bytes of the networking rule
			responseBytes := buildResponseBytes("NetworkingRule", string(networkingRule))

			// Send the networking rule to the front-end
			CBLogger.Debug("Send the networking rule to AdminWeb frontend")
			sendErr := sendMessageToAllPool(responseBytes)
			if sendErr != nil {
				CBLogger.Error(sendErr)
			}
			CBLogger.Tracef("Message Written: %s", event.Kv.Value)
		}
	}
	CBLogger.Debugf("End to watch \"%v\"", etcdkey.NetworkingRule)
}

// func watchConfigurationInformation(wg *sync.WaitGroup, etcdClient *clientv3.Client) {
// 	defer wg.Done()

// 	// It doesn't work for the time being
// 	// Watch "/registry/cloud-adaptive-network/configuration-information"
// 	CBLogger.Debugf("Start to watch \"%v\"", etcdkey.ConfigurationInformation)
// 	watchChan1 := etcdClient.Watch(context.Background(), etcdkey.ConfigurationInformation, clientv3.WithPrefix())
// 	for watchResponse := range watchChan1 {
// 		for _, event := range watchResponse.Events {
// 			CBLogger.Tracef("Watch - %s %q : %q", event.Type, event.Kv.Key, event.Kv.Value)
// 			//slicedKeys := strings.Split(string(event.Kv.Key), "/")
// 			//for _, value := range slicedKeys {
// 			//	fmt.Println(value)
// 			//}
// 		}
// 	}
// 	CBLogger.Debugf("End to watch \"%v\"", etcdkey.ConfigurationInformation)
// }

func watchStatusInformation(wg *sync.WaitGroup, etcdClient *clientv3.Client) {
	defer wg.Done()

	// Watch "/registry/cloud-adaptive-network/status/information/{cladnet-id}/{host-id}"
	CBLogger.Debugf("Start to watch \"%v\"", etcdkey.StatusInformation)
	watchChan1 := etcdClient.Watch(context.Background(), etcdkey.StatusInformation, clientv3.WithPrefix())
	for watchResponse := range watchChan1 {
		for _, event := range watchResponse.Events {
			CBLogger.Tracef("Watch - %s %q : %q", event.Type, event.Kv.Key, event.Kv.Value)
			slicedKeys := strings.Split(string(event.Kv.Key), "/")
			parsedHostID := slicedKeys[len(slicedKeys)-1]
			CBLogger.Tracef("ParsedHostID: %v", parsedHostID)

			status := event.Kv.Value
			CBLogger.Tracef("The status: %v", status)

			// Build the response bytes of the networking rule
			responseBytes := buildResponseBytes("NetworkStatus", string(status))

			// Send the networking rule to the front-end
			CBLogger.Debug("Send the status to AdminWeb frontend")
			sendErr := sendMessageToAllPool(responseBytes)
			if sendErr != nil {
				CBLogger.Error(sendErr)
			}
			CBLogger.Tracef("Message Written: %s", event.Kv.Value)
		}
	}
	CBLogger.Debugf("End to watch \"%v\"", etcdkey.Status)
}

func main() {
	CBLogger.Debug("Start cb-network controller .........")

	// Wait for multiple goroutines to complete
	var wg sync.WaitGroup

	// Create DynamicSubnetConfigurator instance
	dscp = cbnet.NewDynamicSubnetConfigurator()
	// etcd Section
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.ETCD.Endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		CBLogger.Fatal(err)
	}

	defer func() {
		errClose := etcdClient.Close()
		if errClose != nil {
			CBLogger.Fatal("Can't close the etcd client", errClose)
		}
	}()

	CBLogger.Infoln("The etcdClient is connected.")

	wg.Add(1)
	go watchNetworkingRule(&wg, etcdClient)

	// wg.Add(1)
	// go watchConfigurationInformation(&wg, etcdClient)

	wg.Add(1)
	go watchStatusInformation(&wg, etcdClient)

	wg.Add(1)
	go RunEchoServer(&wg, config)

	// Waiting for all goroutines to finish
	CBLogger.Info("Waiting for all goroutines to finish")
	wg.Wait()

	CBLogger.Debug("End cb-network controller .........")
}
