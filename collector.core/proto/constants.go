package proto

type Constant string

//S returns string
func (c *Constant) S() string {
	return string(*c)
}
