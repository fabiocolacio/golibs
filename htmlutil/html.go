package htmlutil

import(
    "golang.org/x/net/html"
)

// GetElementsByClass returns a slice of pointers to all element nodes
// that have the specified class
func GetElementsByClass(root *html.Node, class string) []*html.Node {
    ch := make(chan *html.Node, 10)
    var elements []*html.Node

    predicate := func (node *html.Node) bool {
        if node.Type == html.ElementNode {
            for _, attr := range node.Attr {
                if attr.Key == "class" && attr.Val == class {
                    return true
                }
            }
        }
        return false
    }

    go DepthFirstAccumulator(root, predicate, ch)
    for element := range ch {
        elements = append(elements, element)
    }

    return elements
}

// GetElementById searches for an html element with the given id attribute,
// and returns a pointer to that element's node. The node parameter should be
// the root of an html tree to search.
func GetElementById(root *html.Node, id string) *html.Node {
    predicate := func (node *html.Node) bool {
        if node.Type == html.ElementNode {
            for _, attr := range node.Attr {
                if attr.Key == "id" && attr.Val == id {
                    return true
                }
            }
        }
        return false
    }

    return DepthFirstSearch(root, predicate)
}

// DepthFirstAccumulator is the same as DepthFirstSearch, except it finds all matches in the tree,
// instead of stopping the search at the first match. Each time a node that fulfills the predicate,
// it is sent to ch. ch is closed after the enitre tree has been searched.
func DepthFirstAccumulator(root *html.Node, predicate func(*html.Node) bool, ch chan *html.Node) {
    if root == nil {
        close(ch)
        return
    }

    discovered := make(map[*html.Node]bool)

    nextNode := func(node *html.Node) *html.Node {
        if node == nil {
            return nil
        }

        discovered[node] = true

        if node.FirstChild != nil && !discovered[node.FirstChild]{
            return node.FirstChild
        }

        if node.NextSibling != nil && !discovered[node.NextSibling] {
            return node.NextSibling
        }

        return node.Parent
    }

    for node := root; node != nil; node = nextNode(node) {
        if !discovered[node] && predicate(node) {
            ch <- node
        }
    }

    close(ch)
}

// DepthFirstSearch performs a depth-first search on the tree who's root node is given
// by root. The search ends when it finds a Node which satisfies the predicate function.
// If no suitable candidate was found, DepthFirstSearch returns nil.
func DepthFirstSearch(root *html.Node, predicate func(*html.Node) bool) *html.Node {
    if root == nil {
        return nil
    }

    for node := root; node != nil; node = node.NextSibling {
        if predicate(node) {
            return node
        } else if res := DepthFirstSearch(node.FirstChild, predicate); res != nil {
            return res 
        }
    }

    return nil
}

