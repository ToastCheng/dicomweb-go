package dicomweb

import (
	"errors"
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

func TestClientWithOptionFunc(t *testing.T) {
	bearer := "Bearer f2c45335-6bb1-4caf-99d1-7e0849bcad0d"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, bearer, r.Header.Get("Authorization"))
	}))
	c := NewClient(ClientOption{
		QIDOEndpoint: ts.URL,
		WADOEndpoint: ts.URL,
		STOWEndpoint: ts.URL,
		OptionFuncs: &[]OptionFunc{
			func(req *http.Request) error {
				req.Header.Set("Authorization", bearer)
				return nil
			},
		},
	})

	// just make an arbitrary request to mock server.
	qido := QIDORequest{
		Type:             Study,
		StudyInstanceUID: "study-instance-id",
	}
	c.Query(qido)
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

func TestQIDOQueryAllStudyWithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/studies?00080050=an&limit=1", r.URL.String())
	}))

	c := NewClient(ClientOption{
		QIDOEndpoint: ts.URL,
	})

	qido := QIDORequest{
		Type:            Study,
		Limit:           1,
		AccessionNumber: "an",
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

func TestQIDOQueryUnspecifyType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/studies/study-id/series/series-id/instances", r.URL.String())
	}))
	c := NewClient(ClientOption{
		QIDOEndpoint: ts.URL,
	})

	studyInstanceUID := "study-id"
	seriesInstanceUID := "series-id"
	qido := QIDORequest{
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
	}
	_, err := c.Query(qido)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("failed to query: need to specify query type"), err)
	}
}

func TestQIDOQueryInternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
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
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("500 Internal Server Error"), err)
	}
}

func TestQIDOQueryInvalidURL(t *testing.T) {
	c := NewClient(ClientOption{
		QIDOEndpoint: "%$^",
	})

	qido := QIDORequest{
		Type: Study,
	}
	_, err := c.Query(qido)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid URL escape \"%$^\"")
	}
}

func TestWADORetrieveWithAuthenticate(t *testing.T) {
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
	}).WithAuthentication("user:name")

	studyInstanceUID := "study-id"

	wado := WADORequest{
		Type:             StudyRaw,
		StudyInstanceUID: studyInstanceUID,
	}
	_, err := c.Retrieve(wado)
	assert.NoError(t, err)
}

func TestWADOQueryInvalidURL(t *testing.T) {
	c := NewClient(ClientOption{
		WADOEndpoint: "%$^",
	})

	studyInstanceUID := "study-id"

	wado := WADORequest{
		Type:             StudyRaw,
		StudyInstanceUID: studyInstanceUID,
	}
	_, err := c.Retrieve(wado)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid URL escape \"%$^\"")
	}
}

func TestWADORetrieveStudyRawWithStartAttribute(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boundary := "TOAST"
		w.Header().Set("Content-Type", fmt.Sprintf("multipart/related; type=\"application/dicom\"; start=FIRST; boundary=%s", boundary))
		fmt.Fprint(w, `--TOAST
Content-Type: application/dicom
Content-ID: FIRST

part: 0
--TOAST
Content-Type: application/dicom
Content-ID: SECOND

part: 1
--TOAST--`)
	}))
	c := NewClient(ClientOption{
		WADOEndpoint: ts.URL,
	})

	studyInstanceUID := "study-id"

	wado := WADORequest{
		Type:             StudyRaw,
		StudyInstanceUID: studyInstanceUID,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveStudyRaw(t *testing.T) {
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

	wado := WADORequest{
		Type:             StudyRaw,
		StudyInstanceUID: studyInstanceUID,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveStudyRendered(t *testing.T) {
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

	wado := WADORequest{
		Type:             StudyRendered,
		StudyInstanceUID: studyInstanceUID,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveSeriesRaw(t *testing.T) {
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

	wado := WADORequest{
		Type:              SeriesRaw,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveSeriesRendered(t *testing.T) {
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

	wado := WADORequest{
		Type:              SeriesRendered,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveSeriesMetadata(t *testing.T) {
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

	wado := WADORequest{
		Type:              SeriesMetadata,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveInstanceRaw(t *testing.T) {
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
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveInstanceRendered(t *testing.T) {
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
		Type:              InstanceRendered,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
		SOPInstanceUID:    instanceUID,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveInstanceMetadata(t *testing.T) {
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
		Type:              InstanceMetadata,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
		SOPInstanceUID:    instanceUID,
	}
	parts, err := c.Retrieve(wado)
	assert.NoError(t, err)

	for i, p := range parts {
		assert.Equal(t, fmt.Sprintf("part: %d", i), string(p))
	}
}

func TestWADORetrieveFrame(t *testing.T) {
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
		Type:              Frame,
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

func TestWADORetrieveInternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
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
	_, err := c.Retrieve(wado)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("500 Internal Server Error"), err)
	}
}

func TestWADORetrieveInternalUnsupportContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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
	_, err := c.Retrieve(wado)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected Content-Type, should be multipart/related"), err)
	}
}

func TestWADORetrieveInvalidRequest(t *testing.T) {
	c := NewClient(ClientOption{})

	wado := WADORequest{
		Type: InstanceRaw,
	}
	_, err := c.Retrieve(wado)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("parameters does not match the given type"), err)
	}
}

func TestWADOURIReference(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/uri", r.URL.String())
	}))
	c := NewClient(ClientOption{
		WADOEndpoint: ts.URL,
	})

	wado := WADORequest{
		Type:        URIReference,
		RetrieveURL: "/uri",
	}
	c.Retrieve(wado)
}

func TestWADOUnspecifiedType(t *testing.T) {
	c := NewClient(ClientOption{})

	wado := WADORequest{}
	_, err := c.Retrieve(wado)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("parameters does not match the given type"), err)
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

func TestSTOWStoreWithAuthenticate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	c := NewClient(ClientOption{
		STOWEndpoint: ts.URL,
	}).WithAuthentication("user:name")

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
