package dicomweb

// STOWRequest defines the filter option used in STOW queries.
type STOWRequest struct {
	StudyInstanceUID string
	Parts            [][]byte
}
