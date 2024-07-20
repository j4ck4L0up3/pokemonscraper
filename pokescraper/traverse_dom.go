package pokescraper

import (
	"golang.org/x/net/html"
)

// returns true if element is already in the slice, false if not
func isInSlice(slice []string, element string) bool {
	seen := make(map[string]struct{})
	for _, value := range slice {
		seen[value] = struct{}{}
	}
	_, found := seen[element]
	return found
}

// fills array with all attribute values based on key for html.Node
func GetDOMAttrVals(node *html.Node, elem string, attrKey string, values *[]string) {
	if node.Type == html.ElementNode && node.Data == elem {
		if len(node.Attr) != 0 {
			for _, attr := range node.Attr {
				if attr.Key == attrKey {
					*values = append(*values, attr.Val)
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		GetDOMAttrVals(child, elem, attrKey, values)
	}
}

// fills array with batch of attribute values based on key for html.Node
func GetDOMAttrValsBatch(
	node *html.Node,
	elem string,
	attrKey string,
	values *[]string,
	batches ...int,
) *html.Node {
	batch := 1
	if len(batches) > 0 && batches[0] >= 1 {
		batch = batches[0]
	}

	if node.Type == html.ElementNode && node.Data == elem {
		if len(node.Attr) != 0 {
			for _, attr := range node.Attr {
				if attr.Key == attrKey {
					if len(*values) == batch {
						return node
					}
					if !isInSlice(*values, attr.Val) {
						*values = append(*values, attr.Val)
					}
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		GetDOMAttrValsBatch(child, elem, attrKey, values, batch)
	}
	return node
}

// fills array with all text from html.Nodes traversed based on element and attribute key and value
func GetDOMText(
	node *html.Node,
	elem string,
	attrKey string,
	attrVal string,
	elemText *[]string,
) {
	if node.Type == html.ElementNode && node.Data == elem {
		if len(node.Attr) != 0 {
			for _, attr := range node.Attr {
				if attr.Key == attrKey && attr.Val == attrVal {
					getDOMChildText(node, elemText)
					break
				}
				continue
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		GetDOMText(child, elem, attrKey, attrVal, elemText)
	}
}

func getDOMChildText(node *html.Node, elemText *[]string) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			*elemText = append(*elemText, child.Data)
		}
	}
}
