package dicomweb

// WADORequest defines the filter option used in WADO queries.
type WADORequest struct {
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
