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
	// StudyRaw
	StudyRaw WADOType = iota + 1
	// StudyRendered
	StudyRendered
	// SeriesRaw
	SeriesRaw
	// SeriesRendered
	SeriesRendered
	// SeriesMetadata
	SeriesMetadata
	// InstanceRaw
	InstanceRaw
	// InstanceRendered
	InstanceRendered
	// InstanceMetadata
	InstanceMetadata
	// Frame
	Frame
	// URIReference
	URIReference
)
