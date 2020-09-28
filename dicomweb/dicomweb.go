package dicomweb

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/philippfranke/multipart-related/related"
)

// PrettyPrint pretty print JSON object.
func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
}

// Client defines the client for connecting to dicom server.
// For the naming of the member function such as Query, Retrieve, etc., see
// https://www.dicomstandard.org/wp-content/uploads/2018/04/DICOMweb-Cheatsheet.pdf
// for more detail.
type Client struct {
	httpClient    *http.Client
	qidoEndpoint  string
	wadoEndpoint  string
	stowEndpoint  string
	authorization string
	boundary      string
}

type ClientOption struct {
	QIDOEndpoint string
	WADOEndpoint string
	STOWEndpoint string
}

// WithAuthentication configures the client.
func (c *Client) WithAuthentication(auth string) *Client {
	data := []byte(auth)
	authStr := "Basic " + base64.StdEncoding.EncodeToString(data)
	c.authorization = authStr
	return c
}

// WithInsecure create a http client that skip verifying, do not use it in production.
func (c *Client) WithInsecure() *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	c.httpClient = client
	return c
}

// NewClient creates a new client.
func NewClient(option ClientOption) *Client {
	return &Client{
		httpClient:   &http.Client{},
		qidoEndpoint: option.QIDOEndpoint,
		wadoEndpoint: option.WADOEndpoint,
		stowEndpoint: option.STOWEndpoint,
		boundary:     "dicomwebgoWxkTrZ",
	}
}

// Query based on QIDO, query a list of either matched studies, series or instances.
func (c *Client) Query(req QIDORequest) ([]QIDOResponse, error) {
	url := c.qidoEndpoint
	switch req.Type {
	case Study:
		url += "/studies"
	case Series:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series"
	case Instance:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series/" + req.SeriesInstanceUID
		url += "/instances"
	default:
		return nil, errors.New("failed to query: need to specify query type")
	}

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	mp := map[string]string{}
	databytes, _ := json.Marshal(req)
	json.Unmarshal(databytes, &mp)
	for k, v := range mp {
		if k == "Type" {
			continue
		}
		q.Add(k, v)
	}

	r.URL.RawQuery = q.Encode()

	if c.authorization != "" {
		r.Header.Set("Authorization", c.authorization)
	}
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, errors.New(resp.Status)
	}

	result := []QIDOResponse{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// Retrieve based on WADO, retrieve the DICOM image of given id.
func (c *Client) Retrieve(req WADORequest) ([][]byte, error) {
	if ok := req.Validate(); !ok {
		return nil, errors.New("parameters does not match the given type")
	}

	url := c.wadoEndpoint

	switch req.Type {
	case StudyRaw:
		url += "/studies/" + req.StudyInstanceUID
	case StudyRendered:
		url += "/studies/" + req.StudyInstanceUID
		url += "/rendered"
	case SeriesRaw:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series/" + req.SeriesInstanceUID
	case SeriesRendered:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series/" + req.SeriesInstanceUID
		url += "/rendered"
	case SeriesMetadata:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series/" + req.SeriesInstanceUID
		url += "/metadata"
	case InstanceRaw:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series/" + req.SeriesInstanceUID
		url += "/instances/" + req.SOPInstanceUID
	case InstanceRendered:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series/" + req.SeriesInstanceUID
		url += "/instances/" + req.SOPInstanceUID
		url += "/rendered"
	case InstanceMetadata:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series/" + req.SeriesInstanceUID
		url += "/instances/" + req.SOPInstanceUID
		url += "/metadata"
	case Frame:
		url += "/studies/" + req.StudyInstanceUID
		url += "/series/" + req.SeriesInstanceUID
		url += "/instances/" + req.SOPInstanceUID
		url += "/frames/" + strconv.Itoa(req.FrameID)
	case URIReference:
		url = req.RetrieveURL
	}

	r, _ := http.NewRequest("GET", url, nil)
	if c.authorization != "" {
		r.Header.Set("Authorization", c.authorization)
	}
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(b))
	}

	parts := [][]byte{}
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		if params["start"] == "" {
			mr := multipart.NewReader(resp.Body, params["boundary"])
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					return parts, nil
				} else if err != nil {
					log.Fatalf("failed to read next multipart: %v", err)
					return nil, err
				}

				data, err := ioutil.ReadAll(p)
				if err != nil {
					log.Fatalf("failed to read multipart response: %v", err)
					return nil, err
				}
				parts = append(parts, data)
			}
		} else {
			r := related.NewReader(resp.Body, params)
			obj, err := r.ReadObject()
			if err != nil {
				return nil, err
			}
			for _, part := range obj.Values {
				data, err := ioutil.ReadAll(part)
				if err != nil {
					return nil, err
				}
				parts = append(parts, data)
			}
		}
	} else {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		parts = append(parts, data)
	}

	return parts, nil
}

// Store based on STOW, store the DICOM study to PACS server.
func (c *Client) Store(req STOWRequest) (interface{}, error) {
	url := c.stowEndpoint + "/studies/"

	if req.StudyInstanceUID != "" {
		url += req.StudyInstanceUID
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.SetBoundary(c.boundary)
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", "application/dicom")

	for _, p := range req.Parts {
		w, err := writer.CreatePart(header)
		if err != nil {
			return nil, err
		}
		if _, err = w.Write(p); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "multipart/related; type=application/dicom; boundary="+c.boundary)
	if c.authorization != "" {
		r.Header.Set("Authorization", c.authorization)
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	var result interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}
