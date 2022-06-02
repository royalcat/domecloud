package ptree

type PKTree[K comparable, I any] struct {
	root *pknode[K, I]
}

func (t *PKTree[K, I]) Set(key []K, item I) {

}

type pknode[K comparable, I any] struct {
	name     K
	children map[K]*pknode[K, I]
	value    *I
}

func (n *pknode[K, I]) set(nextkeys []K, item I) {
	if len(nextkeys) == 1 {
		name := nextkeys[0]
		n.children[name] = newValueNode(name, &item)
	} else {
		nextname := nextkeys[0]
		if n, ok := n.children[nextname]; !ok || n == nil {
			n.children[nextname] = newEmptyNode[K, I](nextname)

		}
		n.children[nextname].set(nextkeys[1:], item)
	}
}

func newValueNode[K comparable, I any](name K, value *I) *pknode[K, I] {
	return &pknode[K, I]{
		name:     name,
		children: map[K]*pknode[K, I]{},
		value:    value,
	}
}

func newEmptyNode[K comparable, I any](name K) *pknode[K, I] {
	return &pknode[K, I]{
		name:     name,
		children: map[K]*pknode[K, I]{},
	}
}
