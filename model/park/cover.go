package park

type Cover struct {
	Id     int    `json:"id"`
	Place  string `json:"place"`
	Waring int    `json:"waring"`
}

func (c *Cover) String() string {
	return "ID:" + string(c.Id) + "\n" +
		"Place:" + c.Place + "\n" +
		"Waring:" + string(c.Waring) + "\n"
}

type CoverOV struct {
	Place string `json:"place"`
}

type CoverResOV struct {
	WaringList []Cover `json:"waring_list"`
	Detail     []Cover `json:"detail"`
}
