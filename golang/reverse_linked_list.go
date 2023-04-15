package main

type LinkNode struct {
	Next *LinkNode
	Val  int
}

func ReverseLinkList(root *LinkNode) *LinkNode {
	curr := root
	var last *LinkNode

	for curr != nil {
		next := curr.Next
		curr.Next = last
		//  step
		last = curr
		curr = next
	}

	return last
}

func printLinkedList(root *LinkNode) {
	curr := root
	for curr != nil {
		println(curr.Val)
		curr = curr.Next
	}
}

func main() {
	// root := &LinkNode{
	// 	Val: 0,
	// 	Next: &LinkNode{
	// 		Val: 1,
	// 		Next: &LinkNode{
	// 			Val: 2,
	// 			Next: &LinkNode{
	// 				Val:  3,
	// 				Next: nil,
	// 			},
	// 		},
	// 	},
	// }

	root := &LinkNode{
		Val: 0,
		Next: &LinkNode{
			Val:  1,
			Next: nil,
		},
	}

	printLinkedList(root)
	reversed := ReverseLinkList(root)

	println("reversed")

	printLinkedList(reversed)
}
