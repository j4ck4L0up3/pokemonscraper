package pokescraper

import (
	"golang.org/x/net/html"
)

// fills array with all attribute values based on key for html.Node
func TraverseDOMAttr(node *html.Node, elem string, attrKey string, values *[]string) {
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
		TraverseDOMAttr(child, elem, attrKey, values)
	}
}

// fills array with batch of attribute values based on key for html.Node
func TraverseDOMAttrBatch(node *html.Node, elem string, attrKey string, values *[]string, batches ...int) {
	batch := 1
	if len(batches) > 0 && batches[0] >= 1 {
		batch = batches[0]
	}

	if node.Type == html.ElementNode && node.Data == elem {
		if len(node.Attr) != 0 {
			for _, attr := range node.Attr {
				if attr.Key == attrKey {
					if len(*values) == batch {
						return
					}

					*values = append(*values, attr.Val)
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		TraverseDOMAttrBatch(child, elem, attrKey, values, batch)
	}
}

// fills array with all text from html.Nodes traversed based on element and attribute key and value
func TraverseDOMText(node *html.Node, elem string, attrKey string, attrVal string, elemText *[]string) {
	if node.Type == html.ElementNode && node.Data == elem {
		if len(node.Attr) != 0 {
			for _, attr := range node.Attr {
				if attr.Key == attrKey && attr.Val == attrVal {
					for c := node.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.TextNode {
							*elemText = append(*elemText, c.Data)
						}
					}
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		TraverseDOMText(child, elem, attrKey, attrVal, elemText)
	}
}