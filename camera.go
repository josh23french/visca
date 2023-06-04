package visca

// Camera represents a camera
type Camera struct {
	Name string
	conn Connection
}

// NewCamera creates a new Camera
func NewCamera(name, connString string) (*Camera, error) {
	conn, err := NewConnectionFromString(connString)
	if err != nil {
		return nil, err
	}
	return &Camera{
		Name: name,
		conn: conn,
	}, nil
}
