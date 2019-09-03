package codec

import (
	"github.com/fananchong/v-micro/codec"
	"github.com/fananchong/v-micro/common/metadata"
)

// GetHeader get header
func GetHeader(hdr string, md map[string]string) string {
	if hd := md[hdr]; len(hd) > 0 {
		return hd
	}
	return ""
}

// GetHeaders get headers
func GetHeaders(m *codec.Message) {
	m.Method = GetHeader(metadata.METHOD, m.Header)
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
	set(metadata.METHOD, r.Method)

	if v, ok := r.Header[metadata.CONTENTTYPE]; ok && v == "" {
		delete(r.Header, metadata.CONTENTTYPE)
	}
}
