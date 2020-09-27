package dicomweb

// WADORequest defines the filter option used in WADO queries.
type WADORequest struct {
	Type        WADOType
	StudyID     string
	SeriesID    string
	InstanceID  string
	PatientName string
	FrameID     int
	RetrieveURL string
	Annotation  string
	Quality     int
	Viewport    string
	Window      string
}

func (r WADORequest) Validate() bool {
	switch r.Type {
	case StudyRaw:
		return r.StudyID != "" && r.SeriesID == "" && r.InstanceID == ""
	case StudyRendered:
		return r.StudyID != "" && r.SeriesID == "" && r.InstanceID == ""
	case SeriesRaw:
		return r.StudyID != "" && r.SeriesID != "" && r.InstanceID == ""
	case SeriesRendered:
		return r.StudyID != "" && r.SeriesID != "" && r.InstanceID == ""
	case SeriesMetadata:
		return r.StudyID != "" && r.SeriesID != "" && r.InstanceID == ""
	case InstanceRaw:
		return r.StudyID != "" && r.SeriesID != "" && r.InstanceID != ""
	case InstanceRendered:
		return r.StudyID != "" && r.SeriesID != "" && r.InstanceID != ""
	case InstanceMetadata:
		return r.StudyID != "" && r.SeriesID != "" && r.InstanceID != ""
	case Frame:
		return r.StudyID != "" && r.SeriesID != "" && r.InstanceID != "" && r.FrameID != 0
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
