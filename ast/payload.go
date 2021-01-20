package ast

type payload struct {
	values map[interface{}]interface{}
}

func newPayload() *payload {
	return &payload{}
}

func (n *payload) Value(key interface{}) interface{} {
	if n.values == nil {
		return nil
	}

	return n.values[key]
}

func (n *payload) SetValue(key, value interface{}) {
	if n.values == nil {
		n.values = map[interface{}]interface{}{}
	}

	n.values[key] = value
}
