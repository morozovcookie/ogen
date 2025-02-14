// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/morozovcookie/ogen/uri"
)

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Path
	if rawPath := r.URL.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
		}
	}
	if prefix := s.cfg.Prefix; len(prefix) > 0 {
		if strings.HasPrefix(elem, prefix) {
			// Cut prefix from the path.
			elem = strings.TrimPrefix(elem, prefix)
		} else {
			// Prefix doesn't match.
			s.notFound(w, r)
			return
		}
	}
	if len(elem) == 0 {
		s.notFound(w, r)
		return
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"
			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'n': // Prefix: "nullableStrings"
				if l := len("nullableStrings"); len(elem) >= l && elem[0:l] == "nullableStrings" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "POST":
						s.handleNullableStringsRequest([0]string{}, w, r)
					default:
						s.notAllowed(w, r, "POST")
					}

					return
				}
			case 'o': // Prefix: "objectsWithConflicting"
				if l := len("objectsWithConflicting"); len(elem) >= l && elem[0:l] == "objectsWithConflicting" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'A': // Prefix: "ArrayProperty"
					if l := len("ArrayProperty"); len(elem) >= l && elem[0:l] == "ArrayProperty" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleObjectsWithConflictingArrayPropertyRequest([0]string{}, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}
				case 'P': // Prefix: "Properties"
					if l := len("Properties"); len(elem) >= l && elem[0:l] == "Properties" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleObjectsWithConflictingPropertiesRequest([0]string{}, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}
				}
			case 'r': // Prefix: "referencedAllof"
				if l := len("referencedAllof"); len(elem) >= l && elem[0:l] == "referencedAllof" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch r.Method {
					case "POST":
						s.handleReferencedAllofRequest([0]string{}, w, r)
					default:
						s.notAllowed(w, r, "POST")
					}

					return
				}
				switch elem[0] {
				case 'O': // Prefix: "Optional"
					if l := len("Optional"); len(elem) >= l && elem[0:l] == "Optional" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleReferencedAllofOptionalRequest([0]string{}, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}
				}
			case 's': // Prefix: "s"
				if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'i': // Prefix: "imple"
					if l := len("imple"); len(elem) >= l && elem[0:l] == "imple" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'I': // Prefix: "Integer"
						if l := len("Integer"); len(elem) >= l && elem[0:l] == "Integer" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleSimpleIntegerRequest([0]string{}, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}
					case 'O': // Prefix: "Objects"
						if l := len("Objects"); len(elem) >= l && elem[0:l] == "Objects" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleSimpleObjectsRequest([0]string{}, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}
					}
				case 't': // Prefix: "tringsNotype"
					if l := len("tringsNotype"); len(elem) >= l && elem[0:l] == "tringsNotype" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleStringsNotypeRequest([0]string{}, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}
				}
			}
		}
	}
	s.notFound(w, r)
}

// Route is route object.
type Route struct {
	name        string
	operationID string
	pathPattern string
	count       int
	args        [0]string
}

// Name returns ogen operation name.
//
// It is guaranteed to be unique and not empty.
func (r Route) Name() string {
	return r.name
}

// OperationID returns OpenAPI operationId.
func (r Route) OperationID() string {
	return r.operationID
}

// PathPattern returns OpenAPI path.
func (r Route) PathPattern() string {
	return r.pathPattern
}

// Args returns parsed arguments.
func (r Route) Args() []string {
	return r.args[:r.count]
}

// FindRoute finds Route for given method and path.
//
// Note: this method does not unescape path or handle reserved characters in path properly. Use FindPath instead.
func (s *Server) FindRoute(method, path string) (Route, bool) {
	return s.FindPath(method, &url.URL{Path: path})
}

// FindPath finds Route for given method and URL.
func (s *Server) FindPath(method string, u *url.URL) (r Route, _ bool) {
	var (
		elem = u.Path
		args = r.args
	)
	if rawPath := u.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
		}
		defer func() {
			for i, arg := range r.args[:r.count] {
				if unescaped, err := url.PathUnescape(arg); err == nil {
					r.args[i] = unescaped
				}
			}
		}()
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"
			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'n': // Prefix: "nullableStrings"
				if l := len("nullableStrings"); len(elem) >= l && elem[0:l] == "nullableStrings" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch method {
					case "POST":
						// Leaf: NullableStrings
						r.name = "NullableStrings"
						r.operationID = "nullableStrings"
						r.pathPattern = "/nullableStrings"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}
			case 'o': // Prefix: "objectsWithConflicting"
				if l := len("objectsWithConflicting"); len(elem) >= l && elem[0:l] == "objectsWithConflicting" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'A': // Prefix: "ArrayProperty"
					if l := len("ArrayProperty"); len(elem) >= l && elem[0:l] == "ArrayProperty" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: ObjectsWithConflictingArrayProperty
							r.name = "ObjectsWithConflictingArrayProperty"
							r.operationID = "objectsWithConflictingArrayProperty"
							r.pathPattern = "/objectsWithConflictingArrayProperty"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				case 'P': // Prefix: "Properties"
					if l := len("Properties"); len(elem) >= l && elem[0:l] == "Properties" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: ObjectsWithConflictingProperties
							r.name = "ObjectsWithConflictingProperties"
							r.operationID = "objectsWithConflictingProperties"
							r.pathPattern = "/objectsWithConflictingProperties"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				}
			case 'r': // Prefix: "referencedAllof"
				if l := len("referencedAllof"); len(elem) >= l && elem[0:l] == "referencedAllof" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch method {
					case "POST":
						r.name = "ReferencedAllof"
						r.operationID = "referencedAllof"
						r.pathPattern = "/referencedAllof"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}
				switch elem[0] {
				case 'O': // Prefix: "Optional"
					if l := len("Optional"); len(elem) >= l && elem[0:l] == "Optional" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: ReferencedAllofOptional
							r.name = "ReferencedAllofOptional"
							r.operationID = "referencedAllofOptional"
							r.pathPattern = "/referencedAllofOptional"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				}
			case 's': // Prefix: "s"
				if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'i': // Prefix: "imple"
					if l := len("imple"); len(elem) >= l && elem[0:l] == "imple" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'I': // Prefix: "Integer"
						if l := len("Integer"); len(elem) >= l && elem[0:l] == "Integer" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: SimpleInteger
								r.name = "SimpleInteger"
								r.operationID = "simpleInteger"
								r.pathPattern = "/simpleInteger"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
					case 'O': // Prefix: "Objects"
						if l := len("Objects"); len(elem) >= l && elem[0:l] == "Objects" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: SimpleObjects
								r.name = "SimpleObjects"
								r.operationID = "simpleObjects"
								r.pathPattern = "/simpleObjects"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
					}
				case 't': // Prefix: "tringsNotype"
					if l := len("tringsNotype"); len(elem) >= l && elem[0:l] == "tringsNotype" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: StringsNotype
							r.name = "StringsNotype"
							r.operationID = "stringsNotype"
							r.pathPattern = "/stringsNotype"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				}
			}
		}
	}
	return r, false
}
