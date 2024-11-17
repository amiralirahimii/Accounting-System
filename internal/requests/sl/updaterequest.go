package sl

type UpdateRequest struct {
	ID      int
	Code    string
	Title   string
	HasDL   bool
	Version int
}
