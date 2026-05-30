package models

type CustomFronMatter struct {
	nodeID string
	key    string
	value string
}

func NewCustomFrontMatter(nodeID, key, value string) *CustomFronMatter {
	return &CustomFronMatter{
		nodeID: nodeID,
		key:    key,
		value: value,
	}
}

func (cfm *CustomFronMatter) NodeID() string {
	return cfm.nodeID
}

func (cfm *CustomFronMatter) Key() string {
	return cfm.key
}

func (cfm *CustomFronMatter) Value() string {
	return cfm.value
}
