package visca

// Controller controls cameras
type Controller struct {
	cameras []Camera
}

// NewController creates a new controller
func NewController() *Controller {
	return &Controller{
		cameras: make([]Camera, 0),
	}
}

// AddCamera adds a camera to the controller
func (c *Controller) AddCamera(camera *Camera) {
	c.cameras = append(c.cameras, *camera)
}
