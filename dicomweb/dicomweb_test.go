package dicomweb

import (
	"log"
	"testing"
)

func TestQIDOQueryStudy(t *testing.T) {
	c := NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170125112931.11"
	qido := QIDORequest{
		Type:             Study,
		StudyInstanceUID: studyInstanceUID,
	}
	resp, err := c.Query(qido)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp) > 0 && resp[0].StudyInstanceUID.Value[0].(string) != studyInstanceUID {
		log.Fatalf("expect: %s, get: %s", studyInstanceUID, resp[0].StudyInstanceUID)
	}
}

func TestQIDOQuerySeries(t *testing.T) {
	c := NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
	seriesInstanceUID := "2.25.720409440530442732085780991589110433975"
	qido := QIDORequest{
		Type:              Series,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
	}
	resp, err := c.Query(qido)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp) > 0 && resp[0].StudyInstanceUID.Value[0].(string) != studyInstanceUID {
		log.Fatalf("expect: %s, get: %s", studyInstanceUID, resp[0].StudyInstanceUID.Value[0].(string))
	}
	if len(resp) > 0 && resp[0].SeriesInstanceUID.Value[0].(string) != seriesInstanceUID {
		log.Fatalf("expect: %s, get: %s", seriesInstanceUID, resp[0].SeriesInstanceUID.Value[0].(string))
	}
}

func TestQuerySeries(t *testing.T) {
	c := NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
	seriesInstanceUID := "2.25.687032174858108535882385160051760343725"
	instanceUID := "773645909590137995838355818619864160367"
	qido := QIDORequest{
		Type:              Instance,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
		SOPInstanceUID:    instanceUID,
	}
	resp, err := c.Query(qido)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp) > 0 && resp[0].StudyInstanceUID.Value[0].(string) != studyInstanceUID {
		log.Fatalf("expect: %s, get: %s", studyInstanceUID, resp[0].StudyInstanceUID.Value[0].(string))
	}
	if len(resp) > 0 && resp[0].SeriesInstanceUID.Value[0].(string) != seriesInstanceUID {
		log.Fatalf("expect: %s, get: %s", seriesInstanceUID, resp[0].SeriesInstanceUID.Value[0].(string))
	}
	if len(resp) > 0 && resp[0].SOPInstanceUID.Value[0].(string) != instanceUID {
		log.Fatalf("expect: %s, get: %s", instanceUID, resp[0].StudyInstanceUID.Value[0].(string))
	}
}

func TestWADORetrieve(t *testing.T) {
	c := NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
	seriesInstanceUID := "2.25.687032174858108535882385160051760343725"
	instanceUID := "773645909590137995838355818619864160367"

	wado := WADORequest{
		StudyID:    studyInstanceUID,
		SeriesID:   seriesInstanceUID,
		InstanceID: instanceUID,
	}
	parts, err := c.Retrieve(wado)
	if err != nil {
		log.Fatal(err)
	}
	for i, p := range parts {
		if p == nil {
			log.Fatalf("data is empty on #%d part", i)
		}
		// save it into file like this:
		// ioutil.WriteFile("tmp/test_"+strconv.Itoa(i)+".dcm", p, 0666)
	}
}
