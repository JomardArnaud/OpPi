package oppi

//OpCircle blabla
type OpCircle struct {
	Radius float64
	Pos    OpVector2f
}

//CollideCircle test the collision with a another OpCircle
func (c *OpCircle) CollideCircle(c2 OpCircle) bool {
	return Distance(c.Pos, c2.Pos) <= c.Radius+c2.Radius
}

//CollideRect test the collision with a OpRect4f
func (c *OpCircle) CollideRect(r OpRect4f) bool {
	circleBuffer := c.Pos
	if c.Pos.X < r.X {
		circleBuffer.X = r.X
	} else if c.Pos.X > r.X+r.W {
		circleBuffer.X = r.X + r.W
	}
	if c.Pos.Y < r.Y {
		circleBuffer.Y = r.Y
	} else if c.Pos.Y > r.Y+r.H {
		circleBuffer.Y = r.Y + r.H
	}
	if Distance(c.Pos, circleBuffer) <= c.Radius {
		return true
	}
	return false
}
