package types

type HandlerDocAttr string

const (
	HandlerDocAttrSummary     HandlerDocAttr = "summary:"
	HandlerDocAttrDescription HandlerDocAttr = "description:"
)

type HandlerAttr string

const (
	HandlerAttrDoc     HandlerAttr = "@doc"
	HandlerAttrHandler HandlerAttr = "@handler"
)
