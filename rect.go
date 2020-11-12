package oppi

//OpRect4f blabla
type OpRect4f struct {
	X, Y, W, H float64
}

//Reset all fields' OpRect4f
func (r *OpRect4f) Reset(nPos, nSize OpVector2f) {
	r.X = nPos.X
	r.Y = nPos.Y
	r.H = nSize.X
	r.W = nSize.Y
}

//GetSize of the OpRect4f
func (r *OpRect4f) GetSize() OpVector2f {
	return OpVector2f{r.W, r.H}
}

//GetPos of the OpRect4f
func (r *OpRect4f) GetPos() OpVector2f {
	return OpVector2f{r.X, r.Y}
}

//GetDistanceFromCenter from the center of rectangle and circle
func (r *OpRect4f) GetDistanceFromCenter(c OpCircle) float64 {
	return Distance(OpVector2f{r.X + r.W/2, r.Y + r.H/2}, c.Pos)
}

//ColliderRectangle test collision with a another OpRect4f
func (r *OpRect4f) ColliderRectangle(r2 OpRect4f) bool {
	if r.X < r2.X+r.W && r.X+r.W > r2.X &&
		r.Y < r2.Y+r2.H && r.Y+r.H > r2.Y {
		return true
	}
	return false
}
