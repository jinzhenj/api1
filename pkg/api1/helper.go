package api1

func (c *HasComments) AddSemComment(key string, val interface{}) {
	if c.SemComments == nil {
		c.SemComments = make(map[string]interface{})
	}
	if c.SemComments[key] == nil {
		c.SemComments[key] = val
	} else if a, ok := c.SemComments[key].([]interface{}); ok {
		c.SemComments[key] = append(a, val)
	} else {
		var a []interface{}
		c.SemComments[key] = append(a, c.SemComments[key], val)
	}
}
