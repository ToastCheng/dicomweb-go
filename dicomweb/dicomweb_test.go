package dicomweb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQIDOQueryCertainStudy(t *testing.T) {
	c := NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170125112931.11"
	qido := QIDORequest{
		Type:             Study,
		StudyInstanceUID: studyInstanceUID,
	}
	resp, err := c.Query(qido)
	assert.NoError(t, err)
	if len(resp) > 0 {
		assert.Equal(t, studyInstanceUID, resp[0].StudyInstanceUID.Value[0].(string))
	}
}

func TestQIDOQueryCertainSeries(t *testing.T) {
	c := NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
	seriesInstanceUID := "2.25.720409440530442732085780991589110433975"
	qido := QIDORequest{
		Type:              Series,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
	}
	resp, err := c.Query(qido)
	assert.NoError(t, err)
	if len(resp) > 0 {
		assert.Equal(t, studyInstanceUID, resp[0].StudyInstanceUID.Value[0].(string))
	}
	if len(resp) > 0 {
		assert.Equal(t, seriesInstanceUID, resp[0].SeriesInstanceUID.Value[0].(string))
	}
}

func TestQIDOQueryCertainInstance(t *testing.T) {
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
	assert.NoError(t, err)
	if len(resp) > 0 {
		assert.Equal(t, studyInstanceUID, resp[0].StudyInstanceUID.Value[0].(string))
	}
	if len(resp) > 0 {
		assert.Equal(t, seriesInstanceUID, resp[0].SeriesInstanceUID.Value[0].(string))
	}
	if len(resp) > 0 {
		assert.Equal(t, instanceUID, resp[0].SOPInstanceUID.Value[0].(string))
	}
}

func TestWADORetrieve(t *testing.T) {
	c := NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

	studyInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.1"
	seriesInstanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.2"
	instanceUID := "1.3.6.1.4.1.25403.345050719074.3824.20170126085406.3"

	wado := WADORequest{
		Type:       InstanceMetadata,
		StudyID:    studyInstanceUID,
		SeriesID:   seriesInstanceUID,
		InstanceID: instanceUID,
		FrameID:    1,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.NotNil(t, p, "data is empty on #%d part", i)
		// save it into file like this:
		// err := ioutil.WriteFile("test_"+strconv.Itoa(i)+".dcm", p, 0666)
		assert.NoError(t, err)
	}
}

func TestSTOWStore(t *testing.T) {
	c := NewClient("https://server.dcmjs.org/dcm4chee-arc/aets/DCM4CHEE/rs")

	parts := [][]byte{}
	// read your data like this:
	// file, err := os.Open("data.dcm")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// b := make([]byte, 100)
	// file.Read(b)
	// parts = append(parts, b)

	stow := STOWRequest{
		Parts: parts,
	}
	_, err := c.Store(stow)
	assert.NoError(t, err)
}
