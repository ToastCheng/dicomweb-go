# DICOMweb Go

[![license](https://img.shields.io/badge/license-MIT-blue)](https://github.com/toastcheng/dicomweb-go/blob/master/LICENSE.md)
[![GoDoc](https://img.shields.io/badge/go-doc-blue)](https://pkg.go.dev/github.com/toastcheng/dicomweb-go/dicomweb)
[![Go Report Card](https://goreportcard.com/badge/github.com/toastcheng/dicomweb-go)](https://goreportcard.com/report/github.com/toastcheng/dicomweb-go)
[![Coverage Status](https://coveralls.io/repos/github/ToastCheng/dicomweb-go/badge.svg)](https://coveralls.io/github/ToastCheng/dicomweb-go)
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Ftoastcheng%2Fdicomweb-go%2Fbadge&style=flat-square)](https://actions-badge.atrox.dev/toastcheng/dicomweb-go/goto)


## Introduction
A DICOMweb client for Golang.

There are plenty of packages that allow you to read DICOM files in Go whereas not much for communicating with DICOM server. 

Currently there are DICOM servers such as dcm4chee, Orthanc, etc., that support read/write DICOM by HTTP protocol, known as [DICOMweb](https://www.dicomstandard.org/dicomweb).

This package provides a simple DICOMweb client that allows you to query DICOM info (QIDO), retrieve DICOM files (WADO), and store DICOM files (STOW).


## Documentation
* pkg.go.dev : https://pkg.go.dev/github.com/toastcheng/dicomweb-go/dicomweb
* Dicomweb : https://www.dicomstandard.org/dicomweb

## Getting Started
### Installation
```
go get github.com/toastcheng/dicomweb-go/dicomweb
```

### Requirements
* Go 1.12+

### Quick Examples

note: for demonstration, the endpoint is set to a `dcm4chee` server hosted by `dcmjs.org`. Change it to your DICOM server instead.
#### Query all study
```go
client := dicomweb.NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

qido := dicomweb.QIDORequest{
    Type: dicomweb.Study,
}
resp, err := client.Query(qido)
if err != nil {
    fmt.Errorf("faild to query: %v", err)
}
```

#### Query all series under specific study
```go

client := dicomweb.NewClient(dicomweb.ClientOption{
    QIDOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
    WADOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
    STOWEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
})

studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
qido := dicomweb.QIDORequest{
    Type:              dicomweb.Series,
    StudyInstanceUID:  studyInstanceUID,

}
resp, err := client.Query(qido)
if err != nil {
    fmt.Errorf("faild to query: %v", err)
}
fmt.Println(resp)
```

##### Retrieve the DICOM file
```go
client := dicomweb.NewClient(dicomweb.ClientOption{
    QIDOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
    WADOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
    STOWEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
})
studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
seriesInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.2"
instanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.3"

wado := dicomweb.WADORequest{
    Type:              dicomweb.InstanceRaw,
    StudyInstanceUID:  studyInstanceUID,
    SeriesInstanceUID: seriesInstanceUID,
    SOPInstanceUID:    instanceUID,
    FrameID:           1,
}
parts, err := client.Retrieve(wado)
if err != nil {
    fmt.Errorf("faild to query: %v", err)
}

for i, p := range parts {
    // save it into file like this:
    err := ioutil.WriteFile("/tmp/test_"+strconv.Itoa(i)+".dcm", p, 0666)
    if err != nil {
        fmt.Errorf("faild to retrieve: %v", err)
    }
}
```

##### Store the DICOM file

```go
client := dicomweb.NewClient(dicomweb.ClientOption{
    STOWEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
})

parts := [][]byte{}
// read your data like this:
for i := 0; i < 10; i++ {
    fname := fmt.Sprintf("data_%d.dcm", i)
    b, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Fatal(err)
    }
    parts = append(parts, b)
}

stow := dicomweb.STOWRequest{
    StudyInstanceUID: "1.2.840.113820.0.20200429.174041.3",
    Parts:            parts,
}
resp, err := c.Store(stow)
if err != nil {
    fmt.Errorf("faild to query: %v", err)
}
fmt.Println(resp)
```

## Contributing

This project is still in development, any contributions, issues and feature requests are welcome!
Please check out the [issues page](https://github.com/toastcheng/dicomweb-go/issues).

## License

`dicomweb-go` is available under the [MIT](https://github.com/toastcheng/dicomweb-go/blob/master/LICENSE.md) license.
