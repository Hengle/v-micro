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
	m.ID = GetHeader("Micro-Id", m.Header)
	m.Method = GetHeader("Micro-Method", m.Header)
	m.Service = GetHeader("Micro-Service", m.Header)
	m.Error = GetHeader("Micro-Error", m.Header)
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
	set("Micro-Id", r.ID)
	set("Micro-Service", r.Service)
	set("Micro-Method", r.Method)
	set("Micro-Error", r.Error)
}
