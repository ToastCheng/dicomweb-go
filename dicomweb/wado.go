package dicomweb

// WADORequest defines the filter option used in WADO queries.
type WADORequest struct {
	Type              WADOType
	StudyInstanceUID  string
	SeriesInstanceUID string
	SOPInstanceUID    string
	PatientName       string
	FrameID           int
	RetrieveURL       string
	Annotation        string
	Quality           int
	Viewport          string
	Window            string
}

// Validate validates if the request is valid.
func (r WADORequest) Validate() bool {
	switch r.Type {
	case StudyRaw:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID == "" && r.SOPInstanceUID == ""
	case StudyRendered:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID == "" && r.SOPInstanceUID == ""
	case SeriesRaw:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID != "" && r.SOPInstanceUID == ""
	case SeriesRendered:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID != "" && r.SOPInstanceUID == ""
	case SeriesMetadata:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID != "" && r.SOPInstanceUID == ""
	case InstanceRaw:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID != "" && r.SOPInstanceUID != ""
	case InstanceRendered:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID != "" && r.SOPInstanceUID != ""
	case InstanceMetadata:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID != "" && r.SOPInstanceUID != ""
	case Frame:
		return r.StudyInstanceUID != "" && r.SeriesInstanceUID != "" && r.SOPInstanceUID != "" && r.FrameID != 0
	case URIReference:
		return r.RetrieveURL != ""
	}
	return false
}

// WADOType defines the object to query.
type WADOType int

const (
	// StudyRaw raw study.
	StudyRaw WADOType = iota + 1
	// StudyRendered rendered study.
	StudyRendered
	// SeriesRaw raw series.
	SeriesRaw
	// SeriesRendered rendered series.
	SeriesRendered
	// SeriesMetadata series metadata.
	SeriesMetadata
	// InstanceRaw raw instance.
	InstanceRaw
	// InstanceRendered rendered instance.
	InstanceRendered
	// InstanceMetadata instance metadata.
	InstanceMetadata
	// Frame frame.
	Frame
	// URIReference URI reference.
	URIReference
)
