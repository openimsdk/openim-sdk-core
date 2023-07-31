package content_type

import (
	"strings"
)

var ext = map[string]string{
	"html": "text/html",
	"htm":  "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"json": "application/json",
	"xml":  "application/xml",
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"gif":  "image/gif",
	"bmp":  "image/bmp",
	"tif":  "image/tiff",
	"tiff": "image/tiff",
	"ico":  "image/x-icon",
	"svg":  "image/svg+xml",
	"webp": "image/webp",
	"mp4":  "video/mp4",
	"avi":  "video/x-msvideo",
	"mkv":  "video/x-matroska",
	"mov":  "video/quicktime",
	"wmv":  "video/x-ms-wmv",
	"flv":  "video/x-flv",
	"webm": "video/webm",
	"3gp":  "video/3gpp",
	"mp3":  "audio/mpeg",
	"wav":  "audio/wav",
	"ogg":  "audio/ogg",
	"aac":  "audio/aac",
	"wma":  "audio/x-ms-wma",
	"flac": "audio/flac",
	"mid":  "audio/midi",
	"midi": "audio/midi",
	"weba": "audio/webm",
	"pdf":  "application/pdf",
	"doc":  "application/msword",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"xls":  "application/vnd.ms-excel",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"ppt":  "application/vnd.ms-powerpoint",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"txt":  "text/plain",
	"csv":  "text/csv",
	"zip":  "application/zip",
	"rar":  "application/x-rar-compressed",
	"tar":  "application/x-tar",
	"gz":   "application/gzip",
	"exe":  "application/x-msdownload",
	"msi":  "application/x-msi",
	"deb":  "application/x-debian-package",
	"rpm":  "application/x-redhat-package-manager",
	"sh":   "application/x-sh",
	"bat":  "application/bat",
	"py":   "application/x-python",
	"java": "text/x-java-source",
	"c":    "text/x-csrc",
	"cpp":  "text/x-c++src",
	"h":    "text/x-chdr",
	"hpp":  "text/x-c++hdr",
	"php":  "application/x-php",
	"asp":  "application/x-asp",
	"jsp":  "application/x-jsp",
	"dll":  "application/x-msdownload",
	"jar":  "application/java-archive",
	"war":  "application/java-archive",
	"ear":  "application/java-archive",
}

func GetType(val ...string) string {
	for _, s := range val {
		if len(s) > 1 {
			if i := strings.IndexByte(s, '/'); i > 0 && i < len(s)-1 {
				return s
			}
			if strings.HasPrefix(s, ".") {
				s = s[1:]
			}
			if val, ok := ext[s]; ok {
				return val
			}
			if val, ok := ext[strings.ToLower(s)]; ok {
				return val
			}
		}
	}
	return ""
}
