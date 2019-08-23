package codec

import "github.com/fananchong/v-micro/codec"

// GetHeader get header
func GetHeader(hdr string, md map[string]string) string {
	if hd := md[hdr]; len(hd) > 0 {
		return hd
	}
	return ""
}

// GetHeaders get headers
func GetHeaders(m *codec.Message) {
	m.Method = GetHeader("Micro-Method", m.Header)
	m.Service = GetHeader("Micro-Service", m.Header)
}

// SetHeaders set headers
func SetHeaders(m, r *codec.Message) {
	set := func(hdr, v string) {
		if len(v) == 0 {
			return
		}
		m.Header[hdr] = v
	}

	// set headers
	set("Micro-Service", r.Service)
	set("Micro-Method", r.Method)
}
