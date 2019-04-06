package p1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math"
	"reflect"
	"strings"
)

type Flag_value struct {
	encoded_prefix []uint8
	value          string
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type MerklePatriciaTrie struct {
	db   map[string]Node
	Root string
}

/*
Get method
input: key
output: value, error
*/
func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	return mpt.recursive_Get_helper(string_to_intArr(key), mpt.Root)
}

/*
recursive Get helper method
input: cur_key, cur_address
output: value, error
*/
func (mpt *MerklePatriciaTrie) recursive_Get_helper(cur_key []uint8, cur_address string) (string, error) {

	value := ""
	var curNode = mpt.db[cur_address]

	if curNode.node_type == 2 {
		cur_prefix := compact_decode(curNode.flag_value.encoded_prefix)
		if is_ext_node(curNode.flag_value.encoded_prefix) {
			is_same, commonLength := how_many_array_matched(cur_key, cur_prefix)
			if is_same || commonLength > 0 {
				return mpt.recursive_Get_helper(cur_key[len(cur_prefix):], curNode.flag_value.value)
			} else {
				return value, errors.New("key is not in MPT")
			}
		} else if !is_ext_node(curNode.flag_value.encoded_prefix) {
			is_same, _ := how_many_array_matched(cur_prefix, cur_key)
			if is_same {
				value = curNode.flag_value.value
				return value, nil
			} else {
				return value, errors.New("key is not in MPT")
			}
		}
	} else if curNode.node_type == 1 {
		if len(cur_key) == 0 {
			return curNode.branch_value[16], nil
		}
		for i := 0; i < 16; i++ {
			if uint8(i) == cur_key[0] {
				return mpt.recursive_Get_helper(cur_key[1:], curNode.branch_value[i])
			}
		}
	}

	return value, nil
}

/*
get prefix code
input: prefix code array
output: prefix code
*/
func get_prefix_code(prefix_code_arr []uint8) (prefix_code uint8) {
	return prefix_code_arr[0] / 16
}

/*
convert string to intArr
input: string
output: intArr
*/
func string_to_intArr(str string) (char_code []uint8) {

	for e := range str {
		char_code = append(char_code, str[e]/16)
		char_code = append(char_code, str[e]%16)
	}
	return char_code
}

/*
how many unit in two arrays matched
input: a_arr, b_arr
output: bool, common Length
*/
func how_many_array_matched(a_arr []uint8, b_arr []uint8) (is_same bool, commonLength int) {

	if len(a_arr) == 0 && len(b_arr) == 0 {
		return true, 0
	} else if len(a_arr) == 0 || len(b_arr) == 0 {
		return false, 0
	} else {
		var i int
		for i = 0; i < int(math.Min(float64(len(a_arr)), float64(len(b_arr)))); i++ {
			if a_arr[i] != b_arr[i] {
				break
			}
		}
		if i == len(a_arr) && i == len(b_arr) {
			return true, i
		} else {
			if len(a_arr)-i > 0 && len(b_arr) == i {
				commonLength = i
			} else {
				commonLength = -i
			}
			return false, commonLength
		}
	}

}

/*
Insert method
input: key, value
output:
*/
func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	cur_address := mpt.Root
	var curNode Node
	if cur_address == "" {
		curNode.node_type = 2
		curNode.flag_value.value = new_value
		curNode.flag_value.encoded_prefix = compact_encode(append(string_to_intArr(key), 16))
		mpt.db[curNode.hash_node()] = curNode
		mpt.Root = curNode.hash_node()
		return
	}
	curNode = mpt.db[cur_address]
	has_found_position, has_solved, new_address := mpt.recursive_Insert_helper(string_to_intArr(key), cur_address, new_value)
	if !has_solved && has_found_position {
		var dummy_ext_node Node
		dummy_ext_node.node_type = 1
		dummy_ext_node.branch_value[16] = cur_address
		dummy_bran_node_address := dummy_ext_node.hash_node()
		mpt.db[dummy_bran_node_address] = dummy_ext_node
		new_address := mpt.split_to_insert(string_to_intArr(key), dummy_bran_node_address, 16, new_value)
		delete(mpt.db, dummy_bran_node_address)
		mpt.Root = new_address
	} else {
		mpt.Root = new_address
	}

}

/*
recursive Insert helper method
input: cur_key, cur_address, new_value
output: new hash address
*/
func (mpt *MerklePatriciaTrie) recursive_Insert_helper(cur_key []uint8, cur_address string, new_value string) (has_found_position bool, has_solved bool, new_address string) {

	var curNode = mpt.db[cur_address]

	if curNode.node_type == 2 {
		cur_prefix := compact_decode(curNode.flag_value.encoded_prefix)
		if is_ext_node(curNode.flag_value.encoded_prefix) {
			is_same, commonLength := how_many_array_matched(cur_key, cur_prefix)
			if is_same {
				//already has this key in a ext
				nextNode := mpt.db[curNode.flag_value.value]
				nextNode.branch_value[16] = new_value
				curNode.flag_value.value = nextNode.hash_node()
				mpt.db[nextNode.hash_node()] = nextNode
				delete(mpt.db, cur_address)
				mpt.db[curNode.hash_node()] = curNode
				return true, true, curNode.hash_node()
			} else if commonLength > 0 {
				has_found_position, has_solved, new_address = mpt.recursive_Insert_helper(cur_key[commonLength:], curNode.flag_value.value, new_value)
				if !has_solved {
					fmt.Errorf("unexpected unsolved task from ext")
				} else if has_found_position {
					curNode.flag_value.value = new_address
				}
				delete(mpt.db, cur_address)
				mpt.db[curNode.hash_node()] = curNode
				return true, true, curNode.hash_node()
			} else {
				return true, false, cur_address
			}
		} else if !is_ext_node(curNode.flag_value.encoded_prefix) {
			is_same, _ := how_many_array_matched(cur_key, cur_prefix)
			if is_same {
				//already has this key in a leaf
				curNode.flag_value.value = new_value
				delete(mpt.db, cur_address)
				mpt.db[curNode.hash_node()] = curNode
				return true, true, curNode.hash_node()
			} else {
				return true, false, cur_address
			}
		}
	} else if curNode.node_type == 1 {
		if len(cur_key) == 0 {
			curNode.branch_value[16] = new_value
			delete(mpt.db, cur_address)
			mpt.db[curNode.hash_node()] = curNode
			return true, true, curNode.hash_node()
		}
		for i := 0; i < 16; i++ {
			if uint8(i) == cur_key[0] {
				if curNode.branch_value[i] == "" {
					var newLeafNode Node
					newLeafNode.flag_value.value = new_value
					newLeafNode.node_type = 2
					newLeafNode.flag_value.encoded_prefix = compact_encode(append(cur_key[1:], 16))
					curNode.branch_value[i] = newLeafNode.hash_node()
					mpt.db[newLeafNode.hash_node()] = newLeafNode
					delete(mpt.db, cur_address)
					mpt.db[curNode.hash_node()] = curNode
					return true, true, curNode.hash_node()
				}
				has_found_position, has_solved, new_address = mpt.recursive_Insert_helper(cur_key[1:], curNode.branch_value[i], new_value)

				if !has_solved && has_found_position {
					new_address = mpt.split_to_insert(cur_key[1:], cur_address, i, new_value)
				}

				if has_found_position {
					curNode.branch_value[i] = new_address
				}

				delete(mpt.db, cur_address)
				mpt.db[curNode.hash_node()] = curNode
				return true, true, curNode.hash_node()
			}
		}
	}

	delete(mpt.db, cur_address)
	mpt.db[curNode.hash_node()] = curNode
	return false, false, curNode.hash_node()
}

/*
split node to insert
input: next key, cur address, position, new_value
output: new hash address
*/
func (mpt *MerklePatriciaTrie) split_to_insert(next_key []uint8, cur_address string, position int, new_value string) string {
	var curNode = mpt.db[cur_address]
	next_address := curNode.branch_value[position]
	var nextNode = mpt.db[next_address]
	cur_prefix := compact_decode(nextNode.flag_value.encoded_prefix)
	if nextNode.node_type == 2 {
		if is_ext_node(nextNode.flag_value.encoded_prefix) {
			// is ext
			_, commonLength := how_many_array_matched(next_key, cur_prefix)
			if commonLength <= 0 {
				commonLength = -commonLength

				var newBranNode Node
				newBranNode.node_type = 1

				if commonLength == len(next_key) {
					newBranNode.branch_value[16] = new_value
				} else {
					var newLeafNode Node
					newLeafNode.flag_value.value = new_value
					newLeafNode.node_type = 2
					newLeafNode.flag_value.encoded_prefix = compact_encode(append(next_key[1+commonLength:], 16))
					newBranNode.branch_value[next_key[0+commonLength]] = newLeafNode.hash_node()
					mpt.db[newLeafNode.hash_node()] = newLeafNode
				}

				if len(cur_prefix)-commonLength != 1 {
					nextNode.flag_value.encoded_prefix = compact_encode(cur_prefix[1+commonLength:])
					delete(mpt.db, next_address)
					newBranNode.branch_value[cur_prefix[0+commonLength]] = nextNode.hash_node()
					mpt.db[nextNode.hash_node()] = nextNode
				} else {
					delete(mpt.db, next_address)
					newBranNode.branch_value[cur_prefix[0+commonLength]] = nextNode.flag_value.value
				}

				mpt.db[newBranNode.hash_node()] = newBranNode

				if commonLength > 0 {
					var newExtNode Node
					newExtNode.flag_value.value = newBranNode.hash_node()

					newExtNode.node_type = 2
					newExtNode.flag_value.encoded_prefix = compact_encode(next_key[:commonLength])
					mpt.db[newExtNode.hash_node()] = newExtNode

					curNode.branch_value[position] = newExtNode.hash_node()

					return newExtNode.hash_node()
				} else {
					curNode.branch_value[position] = newBranNode.hash_node()
					return newBranNode.hash_node()
				}

			} else {
				fmt.Errorf("unexpected positive commonLength from ext")
				return nextNode.hash_node()
			}

		} else if !is_ext_node(nextNode.flag_value.encoded_prefix) {
			//is leaf
			_, commonLength := how_many_array_matched(next_key, cur_prefix)
			commonLength = int(math.Abs(float64(commonLength)))

			var newBranNode Node
			newBranNode.node_type = 1

			if commonLength == len(cur_prefix) {
				newBranNode.branch_value[16] = nextNode.flag_value.value
				delete(mpt.db, next_address)
			} else {
				nextNode.flag_value.encoded_prefix = compact_encode(append(cur_prefix[1+commonLength:], 16))
				delete(mpt.db, next_address)
				newBranNode.branch_value[cur_prefix[0+commonLength]] = nextNode.hash_node()
				mpt.db[nextNode.hash_node()] = nextNode
			}

			if commonLength == len(next_key) {
				newBranNode.branch_value[16] = new_value
			} else {
				var newLeafNode Node
				newLeafNode.flag_value.value = new_value
				newLeafNode.node_type = 2
				newLeafNode.flag_value.encoded_prefix = compact_encode(append(next_key[1+commonLength:], 16))
				newBranNode.branch_value[next_key[0+commonLength]] = newLeafNode.hash_node()
				mpt.db[newLeafNode.hash_node()] = newLeafNode
			}

			if commonLength > 0 {
				var newExtNode Node
				newExtNode.flag_value.value = newBranNode.hash_node()
				mpt.db[newBranNode.hash_node()] = newBranNode
				newExtNode.node_type = 2
				newExtNode.flag_value.encoded_prefix = compact_encode(next_key[:commonLength])
				mpt.db[newExtNode.hash_node()] = newExtNode
				curNode.branch_value[position] = newExtNode.hash_node()
				return newExtNode.hash_node()
			} else {
				mpt.db[newBranNode.hash_node()] = newBranNode
				return newBranNode.hash_node()
			}

		}
	} else {
		fmt.Errorf("unexpected branch node appeared")
	}
	return nextNode.hash_node()
}

/*
Delete method
input: key
output: string, error
*/
func (mpt *MerklePatriciaTrie) Delete(key string) (string, error) {
	_, error := mpt.Get(key)
	if error == errors.New("path_not_found") {
		return "path_not_found", errors.New("path_not_found")
	} else {
		mpt.Root = mpt.recursive_Delete_helper(string_to_intArr(key), mpt.Root)

		return "", nil
	}

}

/*
recursive Delete helper method
input: cur_key, cur_address
output: new hash address
*/
func (mpt *MerklePatriciaTrie) recursive_Delete_helper(cur_key []uint8, cur_address string) string {

	var curNode = mpt.db[cur_address]

	if curNode.node_type == 2 {
		cur_prefix := compact_decode(curNode.flag_value.encoded_prefix)
		if is_ext_node(curNode.flag_value.encoded_prefix) {
			is_same, commonLength := how_many_array_matched(cur_key, cur_prefix)
			if is_same || commonLength > 0 {

				next_address := mpt.recursive_Delete_helper(cur_key[len(cur_prefix):], curNode.flag_value.value)

				nextNode := mpt.db[next_address]

				if nextNode.node_type == 1 {
					curNode.flag_value.value = next_address
					delete(mpt.db, cur_address)
					mpt.db[curNode.hash_node()] = curNode
					return curNode.hash_node()
				} else if nextNode.node_type == 2 && is_ext_node(nextNode.flag_value.encoded_prefix) {
					delete(mpt.db, next_address)
					new_prefix := append(cur_prefix, compact_decode(nextNode.flag_value.encoded_prefix)...)
					curNode.flag_value.encoded_prefix = compact_encode(new_prefix)
					curNode.flag_value.value = nextNode.flag_value.value

					mpt.db[curNode.hash_node()] = curNode
					return curNode.hash_node()
				} else if nextNode.node_type == 2 && !is_ext_node(nextNode.flag_value.encoded_prefix) {
					delete(mpt.db, next_address)
					new_prefix := append(cur_prefix, compact_decode(nextNode.flag_value.encoded_prefix)...)
					nextNode.flag_value.encoded_prefix = compact_encode(append(new_prefix, 16))

					mpt.db[nextNode.hash_node()] = nextNode
					delete(mpt.db, cur_address)
					return nextNode.hash_node()
				}

			} else {
				return cur_address
			}
		} else if !is_ext_node(curNode.flag_value.encoded_prefix) {
			is_same, _ := how_many_array_matched(cur_prefix, cur_key)
			if is_same {
				delete(mpt.db, cur_address)
				return ""
			} else {
				return cur_address
			}
		}
	} else if curNode.node_type == 1 {

		if len(cur_key) == 0 {
			curNode.branch_value[16] = ""
		} else {
			for i := 0; i < 16; i++ {
				if uint8(i) == cur_key[0] {
					curNode.branch_value[i] = mpt.recursive_Delete_helper(cur_key[1:], curNode.branch_value[i])
				}
			}
		}
		num, position := how_many_children_under_Branch(curNode)
		if num == 1 {
			return mpt.deal_with_one_child(cur_address, position)
		} else {
			delete(mpt.db, cur_address)
			mpt.db[curNode.hash_node()] = curNode
			return curNode.hash_node()
		}

	}

	return cur_address
}

/*
deal with one child
input: address, position
output: new hash address
*/
func (mpt *MerklePatriciaTrie) deal_with_one_child(cur_address string, position int) (new_address string) {

	curNode := mpt.db[cur_address]
	delete(mpt.db, cur_address)
	if position == 16 {
		var newLeafNode Node
		newLeafNode.node_type = 2
		newLeafNode.flag_value.value = curNode.branch_value[position]
		newLeafNode.flag_value.encoded_prefix = compact_encode([]uint8{16})
		mpt.db[newLeafNode.hash_node()] = newLeafNode
		return newLeafNode.hash_node()
	}
	next_address := curNode.branch_value[position]
	nextNode := mpt.db[next_address]
	if nextNode.node_type == 1 {
		var newExtNode Node
		newExtNode.flag_value.value = next_address
		newExtNode.node_type = 2

		newExtNode.flag_value.encoded_prefix = compact_encode([]uint8{uint8(position)})
		mpt.db[newExtNode.hash_node()] = newExtNode
		return newExtNode.hash_node()
	} else if nextNode.node_type == 2 && is_ext_node(nextNode.flag_value.encoded_prefix) {
		delete(mpt.db, next_address)
		new_prefix := append([]uint8{uint8(position)}, compact_decode(nextNode.flag_value.encoded_prefix)...)
		nextNode.flag_value.encoded_prefix = compact_encode(new_prefix)
		mpt.db[nextNode.hash_node()] = nextNode
		return nextNode.hash_node()
	} else if nextNode.node_type == 2 && !is_ext_node(nextNode.flag_value.encoded_prefix) {
		delete(mpt.db, next_address)
		new_prefix := append([]uint8{uint8(position)}, compact_decode(nextNode.flag_value.encoded_prefix)...)
		nextNode.flag_value.encoded_prefix = compact_encode(append(new_prefix, 16))
		mpt.db[nextNode.hash_node()] = nextNode
		return nextNode.hash_node()
	} else {
		mpt.db[curNode.hash_node()] = curNode
		return curNode.hash_node()
	}
}

/*
how_many_children_under_Branch
input: node
output: number, position
*/
func how_many_children_under_Branch(node Node) (num int, position int) {

	for i := 0; i < 17; i++ {
		if "" != node.branch_value[i] {
			num++
			position = i
		}
	}
	return num, position
}

/*
compact_encode
input: prefix array
output: encoded array
*/
func compact_encode(hex_array []uint8) []uint8 {
	if hex_array[len(hex_array)-1] == 16 {
		hex_array = hex_array[:len(hex_array)-1]
		if len(hex_array)%2 == 0 {
			hex_array = append([]uint8{2, 0}, hex_array...)
		} else {
			hex_array = append([]uint8{3}, hex_array...)
		}
	} else {
		if len(hex_array)%2 == 0 {
			hex_array = append([]uint8{0, 0}, hex_array...)
		} else {
			hex_array = append([]uint8{1}, hex_array...)
		}
	}

	encoded_arr := []uint8{}
	for i := 0; i < len(hex_array); i += 2 {
		encoded_arr = append(encoded_arr, hex_array[i]*16+hex_array[i+1])
	}
	return encoded_arr
}

/*
compact_decode
input: encoded array
output: prefix array
*/
// If Leaf, ignore 16 at the end
func compact_decode(encoded_arr []uint8) []uint8 {

	hex_array := []uint8{}
	for i := 0; i < len(encoded_arr); i++ {
		hex_array = append(hex_array, encoded_arr[i]/16)
		hex_array = append(hex_array, encoded_arr[i]%16)
	}
	switch hex_array[0] {
	case 0:
		hex_array = hex_array[2:]
	case 1:
		hex_array = hex_array[1:]
	case 2:
		hex_array = hex_array[2:]
	case 3:
		hex_array = hex_array[1:]
	}
	return hex_array
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.db = make(map[string]Node)
	mpt.Root = ""
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.Root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart"+cur_hash+"HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart"+cur_hash+"HashEnd", fmt.Sprintf("Hash%v", i), -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}
