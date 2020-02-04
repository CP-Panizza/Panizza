package hash

//通过字符串生成hash值
func hash_key_fun(key string) int {
	var char []byte = []byte(key)
	var seed int = 131
	var hash int = 0
	for _, c := range char {
		hash = hash*seed + int(c)
	}
	return (hash & 0x7FFFFFFF)
}

func str_equal(keyA, keyB string) bool {
	return keyA == keyB
}


type list struct {
	key  string
	val  interface{}
	next *list
}


type Hash struct {
	list_node []*list
	hash_fun  func(key string) int
	equal_fun func(keyA, keyB string) bool
	num       int
}


func (hash *Hash) Add(key string, value interface{}) {
	node := new(list)
	node.key = key
	node.val = value
	hval := hash.hash_fun(key) % hash.num
	//每次都插入到链表的最前面
	node.next = hash.list_node[hval]
	hash.list_node[hval] = node
}


func (hash *Hash) Get(key string) (interface{}, bool) {
	hval := hash.hash_fun(key) % hash.num
	node := hash.list_node[hval]

	for node != nil && !hash.equal_fun(node.key, key) {
		node = node.next
	}

	if node != nil {
		return node.val, true
	} else {
		return nil, false
	}
}

func (hash *Hash)ForEach(f func(i interface{})){
	for _,list := range hash.list_node {
		if list != nil {
			for e := list; e != nil; e = e.next{
				f(e.val)
			}
		}
	}
}


func HashCreate(num int, hash_fun func(key string) int, equal_fun func(keyA, keyB string) bool) *Hash {
	h := new(Hash)
	h.num = num
	h.equal_fun = equal_fun
	h.hash_fun = hash_fun
	h.list_node = make([]*list, num)
	for i := 0; i < num; i++ {
		h.list_node[i] = nil
	}
	return h
}


func NewHash(num int) *Hash {
	h := new(Hash)
	h.num = num
	h.equal_fun = str_equal
	h.hash_fun = hash_key_fun
	h.list_node = make([]*list, num)
	for i := 0; i < num; i++ {
		h.list_node[i] = nil
	}
	return h
}
