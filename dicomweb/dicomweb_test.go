package dicomweb

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientWithAuthentication(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic dXNlcjpwYXNzd29yZA==", r.Header.Get("Authorization"))
	}))
	c := NewClient(ClientOption{
		QIDOEndpoint: ts.URL,
		WADOEndpoint: ts.URL,
		STOWEndpoint: ts.URL,
	}).WithAuthentication("user:password")

	// just make an arbitrary request to mock server.
	qido := QIDORequest{
		Type:             Study,
		StudyInstanceUID: "study-instance-id",
	}
	c.Query(qido)
}

func TestClientWithInsecure(t *testing.T) {
	c := NewClient(ClientOption{}).WithInsecure()

	insecure := c.httpClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify
	assert.Equal(t, true, insecure)
}

func TestQIDOQueryAllStudy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/studies", r.URL.String())
	}))

	c := NewClient(ClientOption{
		QIDOEndpoint: ts.URL,
	})

	qido := QIDORequest{
		Type: Study,
	}
	_, err := c.Query(qido)
	assert.NoError(t, err)

}

func TestQIDOQuerySeries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/studies/study-id/series", r.URL.String())
	}))

	c := NewClient(ClientOption{
		QIDOEndpoint: ts.URL,
	})

	studyInstanceUID := "study-id"
	qido := QIDORequest{
		Type:             Series,
		StudyInstanceUID: studyInstanceUID,
	}
	_, err := c.Query(qido)
	assert.NoError(t, err)
}

func TestQIDOQueryInstance(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/studies/study-id/series/series-id/instances", r.URL.String())
	}))
	c := NewClient(ClientOption{
		QIDOEndpoint: ts.URL,
	})

	studyInstanceUID := "study-id"
	seriesInstanceUID := "series-id"
	qido := QIDORequest{
		Type:              Instance,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
	}
	_, err := c.Query(qido)
	assert.NoError(t, err)
}

func TestWADORetrieve(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boundary := "TOAST"
		w.Header().Set("Content-Type", fmt.Sprintf("multipart/related; type=\"application/dicom\"; boundary=%s", boundary))
		fmt.Fprint(w, `--TOAST
Content-Type: application/dicom

part: 0
--TOAST
Content-Type: application/dicom

part: 1
--TOAST--`)
	}))
	c := NewClient(ClientOption{
		WADOEndpoint: ts.URL,
	})

	studyInstanceUID := "study-id"
	seriesInstanceUID := "series-id"
	instanceUID := "instance-id"

	wado := WADORequest{
		Type:              InstanceRaw,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
		SOPInstanceUID:    instanceUID,
		FrameID:           1,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestSTOWStore(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		assert.Contains(t, contentType, "multipart/")

		multipartReader := multipart.NewReader(r.Body, params["boundary"])
		defer r.Body.Close()

		idx := 0
		for {
			part, err := multipartReader.NextPart()
			if err == io.EOF {
				break
			}
			assert.NoError(t, err)
			defer part.Close()

			fileBytes, err := ioutil.ReadAll(part)
			assert.NoError(t, err)

			assert.Equal(t, fmt.Sprintf("part: %d", idx), string(fileBytes))
			idx++
		}
	}))

	c := NewClient(ClientOption{
		STOWEndpoint: ts.URL,
	})

	parts := [][]byte{}
	for i := 0; i < 3; i++ {
		p := []byte(fmt.Sprintf("part: %d", i))
		parts = append(parts, p)
	}

	stow := STOWRequest{
		StudyInstanceUID: "1.2.840.113820.0.20200429.174041.3",
		Parts:            parts,
	}
	_, err := c.Store(stow)
	assert.NoError(t, err)
}
