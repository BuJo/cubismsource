package main

import "encoding/xml"
import "flag"
import "os"
import "fmt"
import "io/ioutil"

type JbossJvmStatus struct {
	Free  int `xml:"free,attr"`
	Total int `xml:"total,attr"`
	Max   int `xml:"max,attr"`
}
type JbossThreadInfo struct {
	MaxThreads         int `xml:"maxThreads,attr"`
	MinSpareThreads    int `xml:"minSpareThreads,attr"`
	MaxSpareThreads    int `xml:"maxSpareThreads,attr"`
	CurrentThreadCount int `xml:"currentThreadCount,attr"`
	CurrentThreadsBusy int `xml:"currentThreadsBusy,attr"`
}
type JbossRequestInfo struct {
	MaxTime        int `xml:"maxTime,attr"`
	ProcessingTime int `xml:"processingTime,attr"`
	RequestCount   int `xml:"requestCount,attr"`
	ErrorCount     int `xml:"errorCount,attr"`
	BytesReceived  int `xml:"bytesReceived,attr"`
	BytesSent      int `xml:"bytesSent,attr"`
}
type JbossWorker struct {
	Stage                 string `xml:"stage,attr"`
	RequestProcessingTime int    `xml:"requestProcessingTime,attr"`
	RequestBytesSent      int    `xml:"requestBytesSent,attr"`
	RequestBytesRecieved  int    `xml:"requestBytesRecieved,attr"`
	RemoteAddr            string `xml:"remoteAddr,attr"`
	VirtualHost           string `xml:"virtualHost,attr"`
	Method                string `xml:"method,attr"`
	CurrentUri            string `xml:"currentUri,attr"`
	CurrentQueryString    string `xml:"currentQueryString,attr"`
	Protocol              string `xml:"protocol,attr"`
}
type JbossConnector struct {
	Name        string           `xml:"name"`
	ThreadInfo  JbossThreadInfo  `xml:"threadInfo"`
	RequestInfo JbossRequestInfo `xml:"requestInfo"`
	Workers     []JbossWorker    `xml:"workers"`
}
type JbossStatus struct {
	XMLName   xml.Name       `xml:"status"`
	JvmStatus JbossJvmStatus `xml:"jvm>memory"`
	Connector JbossConnector `xml:"connector"`
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Input file missing")
		os.Exit(1)
	}

	data, readError := ioutil.ReadFile(args[0])
	if readError != nil || len(data) < 1 {
		fmt.Printf("error: %v", readError)
		os.Exit(1)
	}

	v := JbossStatus{}

	err := xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("xml error: %v", err)
		os.Exit(1)
	}

	fmt.Printf("JVM: Used: %.2f MB\n", (float64)(v.JvmStatus.Total-v.JvmStatus.Free)/1024/1024)
	//fmt.Printf("%#v\n", v)
	//fmt.Printf("%v\n", string(data[0:100]))
}
