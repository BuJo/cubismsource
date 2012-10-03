package jbossinfo

import "io"
import "io/ioutil"
import "encoding/xml"

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
	XMLName    xml.Name         `xml:"status"`
	JvmStatus  JbossJvmStatus   `xml:"jvm>memory"`
	Connectors []JbossConnector `xml:"connector"`
}

func ParseJbossInfoXML(r io.Reader) (*JbossStatus, error) {
	data, readError := ioutil.ReadAll(r)
	if readError != nil || len(data) < 1 {
		return nil, readError
	}

	v := JbossStatus{}
	err := xml.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func NewStatus() *JbossStatus {
	status := new(JbossStatus)
	status.XMLName = xml.Name{"", "status"}
	status.JvmStatus = JbossJvmStatus{}
	status.Connectors = make([]JbossConnector, 1)
	status.Connectors[0].Name = ""
	status.Connectors[0].ThreadInfo = JbossThreadInfo{}

	return status
}

func InfoXML(s *JbossStatus) ([]byte, error) {
	out, err := xml.MarshalIndent(s, " ", "  ")

	return out, err
}
