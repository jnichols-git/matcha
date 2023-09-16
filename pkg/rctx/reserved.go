package rctx

const reserved_fullpath = paramKey("fullpath")

type reservedParams struct {
	fullpath string
}

// get gets a reserved rctx param.
// Returns the value, and true if key is reserved.
func (rps *reservedParams) get(key paramKey) (string, bool) {
	switch string(key) {
	case "fullpath":
		return rps.fullpath, true
	default:
		return "", false
	}
}

// set sets a reserved rctx param.
// Returns true if the key is reserved.
func (rps *reservedParams) set(in *Context, key paramKey, value string) (bool, error) {
	switch key {
	case reserved_fullpath:
		if rps.fullpath == "" {
			if orig := in.parent.Value(reserved_fullpath); orig != nil && orig.(string) != "" {
				value = orig.(string)
			}
			rps.fullpath = value
		}
		return true, nil
	default:
		return false, nil
	}
}

// reset resets the reserved params.
// Some reserved params have special behavior when set multiple times; this sets back to
// default values so that behavior can be replicated on pooled context.
func (rps *reservedParams) reset() {
	rps.fullpath = ""
}
