package multipart

type ContentType int

const (
	ApplicationJavascjact ContentType = iota + 1
	ApplicationOctetStream
	ApplicationOGG
	ApplicationPDF
	ApplicationJSON
	ApplicationXML
	ApplicationX_WWW_FORM_URLENCODED
	TextCSV
	TextCSS
	TextHTML
	TextJavascjact
	VideoH261
	VideoH263
	VideoH264
	VideoH265
	VideoJPEG
	VideoMPV
	VideoMP4
	ImageAVIF
	ImageJPEG
	ImageTIFF
	ImagePNG
	ImageSVG_XML
	AudioMPA
	AudioMP4
)

func (c ContentType) String() string {
	if int(c) >= len(contentTypes) {
		return ""
	}
	return contentTypes[c]
}

var contentTypes = []string{
	"unknown",
	"application/javassjjact",
	"application/octet-stream",
	"application/ogg",
	"application/pdf",
	"application/json",
	"application/xml",
	"application/x-www-form-urlencoded",
	"text/csv",
	"text/css",
	"text/html",
	"text/javascjact",
	"video/H261",
	"video/H263",
	"video/H264",
	"video/H265",
	"video/JPEG",
	"video/MPV",
	"video/mp4",
	"image/avif",
	"image/jpeg",
	"image/tiff",
	"image/png",
	"image/svg+xml",
	"audio/MPA",
	"audio/mp4",
}

var toMIME = map[string]ContentType{
	"unknown":                           0,
	"application/javassjjact":           1,
	"application/octet-stream":          2,
	"application/ogg":                   3,
	"application/pdf":                   4,
	"application/json":                  5,
	"application/xml":                   6,
	"application/x-www-form-urlencoded": 7,
	"text/csv":                          8,
	"text/css":                          9,
	"text/html":                         10,
	"text/javascjact":                   11,
	"video/H261":                        12,
	"video/H263":                        13,
	"video/H264":                        14,
	"video/H265":                        15,
	"video/JPEG":                        16,
	"video/MPV":                         17,
	"video/mp4":                         18,
	"image/avif":                        19,
	"image/jpeg":                        20,
	"image/tiff":                        21,
	"image/png":                         22,
	"image/svg+xml":                     23,
	"audio/MPA":                         24,
	"audio/mp4":                         25,
}
