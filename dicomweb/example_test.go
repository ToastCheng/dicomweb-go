package dicomweb_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/toastcheng/dicomweb-go/dicomweb"
)

func ExampleClient_Query_allStudy() {
	c := dicomweb.NewClient(dicomweb.ClientOption{
		QIDOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
	})

	qido := dicomweb.QIDORequest{
		Type: dicomweb.Study,
	}
	resp, err := c.Query(qido)
	if err != nil {
		fmt.Errorf("faild to query: %v", err)
	}
	fmt.Println(resp)
}

func ExampleClient_Query_certainStudy() {
	c := dicomweb.NewClient(dicomweb.ClientOption{
		QIDOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
	})

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170125112931.11"
	qido := dicomweb.QIDORequest{
		Type:             dicomweb.Study,
		StudyInstanceUID: studyInstanceUID,
	}
	resp, err := c.Query(qido)
	if err != nil {
		fmt.Errorf("faild to query: %v", err)
	}
	fmt.Println(resp[0].StudyInstanceUID.Value[0].(string))
}

func ExampleClient_Query_certainSeries() {
	c := dicomweb.NewClient(dicomweb.ClientOption{
		QIDOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
	})

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
	seriesInstanceUID := "2.25.720409440530442732085780991589110433975"
	qido := dicomweb.QIDORequest{
		Type:              dicomweb.Series,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
	}
	resp, err := c.Query(qido)
	if err != nil {
		fmt.Errorf("faild to query: %v", err)
	}
	fmt.Println(resp)
}

func ExampleClient_Query_certainInstance() {
	c := dicomweb.NewClient(dicomweb.ClientOption{
		QIDOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
	})

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
	seriesInstanceUID := "2.25.687032174858108535882385160051760343725"
	instanceUID := "773645909590137995838355818619864160367"
	qido := dicomweb.QIDORequest{
		Type:              dicomweb.Instance,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
		SOPInstanceUID:    instanceUID,
	}
	resp, err := c.Query(qido)
	if err != nil {
		fmt.Errorf("faild to query: %v", err)
	}
	fmt.Println(resp)
}

func Example_Retrieve() {
	c := dicomweb.NewClient(dicomweb.ClientOption{
		WADOEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
	})

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
	seriesInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.2"
	instanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.3"

	wado := dicomweb.WADORequest{
		Type:       dicomweb.InstanceMetadata,
		StudyID:    studyInstanceUID,
		SeriesID:   seriesInstanceUID,
		InstanceID: instanceUID,
		FrameID:    1,
	}
	parts, err := c.Retrieve(wado)
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
}

func ExampleClient_Store() {
	c := dicomweb.NewClient(dicomweb.ClientOption{
		STOWEndpoint: "https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs",
	})

	parts := [][]byte{}
	// read your data like this:
	file, err := os.Open("/tmp/test_0.dcm")
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 100)
	file.Read(b)
	parts = append(parts, b)

	stow := dicomweb.STOWRequest{
		Parts: parts,
	}
	resp, err := c.Store(stow)
	if err != nil {
		fmt.Errorf("faild to query: %v", err)
	}
	fmt.Println(resp)
}
