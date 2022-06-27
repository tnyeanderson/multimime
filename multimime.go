package multimime

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

// MultipartFindFunc is used to filter parts
//
// Returns true if the part should be included, and false if it should be
// filtered out
type MultipartFindFunc func(*multipart.Part) bool

// Part is a multipart.Part and its Content
type Part struct {
	multipart.Part
	Content []byte
}

// GetParts converts an io.Reader to a multipart.Reader and calls FindParts
func GetParts(r io.Reader, findFunc MultipartFindFunc) ([]Part, error) {
	mr, err := GetMultipartReader(r)
	if err != nil {
		return nil, err
	}
	parts, err := FindParts(mr, findFunc)
	if err != nil {
		return nil, err
	}
	return parts, nil
}

// GetAllParts returns a list of all Parts
func GetAllParts(r io.Reader) ([]Part, error) {
	return GetParts(r, func(_ *multipart.Part) bool { return true })
}

// GetTextParts returns a list of Parts that are text/*
func GetTextParts(r io.Reader) ([]Part, error) {
	return GetParts(r, IsTextPart)
}

// GetPlainTextParts returns a list of Parts that are text/plain
func GetPlainTextParts(r io.Reader) ([]Part, error) {
	return GetParts(r, IsPlainTextPart)
}

// GetHtmlParts returns a list of Parts that are text/html
func GetHtmlParts(r io.Reader) ([]Part, error) {
	return GetParts(r, IsHtmlPart)
}

// GetAttachments returns a list of Parts that are attachments
func GetAttachments(r io.Reader) ([]Part, error) {
	return GetParts(r, IsAttachment)
}

func GetInlineText(r io.Reader) (text string, err error) {
	parts, err := GetParts(r, IsInlineTextPart)
	if err != nil {
		return
	}
	return CombineParts(parts)
}

// CombineParts returns the body of each part, concatenated and separated by new lines
func CombineParts(parts []Part) (text string, err error) {
	for _, part := range parts {
		text = fmt.Sprintf("%s\n%s", text, part.Content)
	}
	return
}

// IsPlainTextPart returns true for text/plain media types
//
// Implements MultipartFindFunc
func IsPlainTextPart(part *multipart.Part) bool {
	mediaType := GetPartType(part)
	return strings.HasPrefix(mediaType, "text/plain")
}

// IsHtmlPart returns true for text/html media types
//
// Implements MultipartFindFunc
func IsHtmlPart(part *multipart.Part) bool {
	mediaType := GetPartType(part)
	return strings.HasPrefix(mediaType, "text/html")
}

// IsTextPart returns true for text media types
//
// Implements MultipartFindFunc
func IsTextPart(part *multipart.Part) bool {
	mediaType := GetPartType(part)
	return strings.HasPrefix(mediaType, "text/")
}

// IsInlineTextPart returns true for text parts that are not attachments
//
// Implements MultipartFindFunc
func IsInlineTextPart(part *multipart.Part) bool {
	result := IsTextPart(part) && !IsAttachment(part)
	return result
}

// IsAttachment returns true if the part disposition starts with attachment
//
// Implements MultipartFindFunc
func IsAttachment(part *multipart.Part) bool {
	disposition := GetPartDisposition(part)
	return strings.HasPrefix(disposition, "attachment")
}

// GetPartDisposition returns the Content-Disposition media type
//
// If the header is not set, an empty string is returned
func GetPartDisposition(part *multipart.Part) string {
	header := part.Header.Get("Content-Disposition")
	disposition, _, err := mime.ParseMediaType(header)
	if err != nil {
		return ""
	}
	return disposition
}

// GetMessageType returns the media type from a multipart.Part
func GetPartType(part *multipart.Part) string {
	header := part.Header.Get("Content-Type")
	contentType, _, err := mime.ParseMediaType(header)
	if err != nil {
		return ""
	}
	return contentType
}

// GetMessageType returns the media type and params from a mail.Message
func GetMessageType(message *mail.Message) (mediatype string, params map[string]string, err error) {
	contentType := message.Header.Get("Content-Type")
	return mime.ParseMediaType(contentType)
}

// GetMultipartReader returns a multipart.Reader from an io.Reader
//
// Returns an error if the io.Reader does not contain an email or if it is not
// a multipart email
func GetMultipartReader(r io.Reader) (*multipart.Reader, error) {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}
	mediaType, params, err := GetMessageType(m)
	if err != nil {
		return nil, err
	}
	if !isMultipart(mediaType) {
		return nil, fmt.Errorf("Not multipart: %s\n", mediaType)
	}
	boundary := params["boundary"]
	mr := multipart.NewReader(m.Body, boundary)
	return mr, nil
}

// GetPartContent returns the content (body) of a multipart.Part
func GetPartContent(part multipart.Part) (content []byte) {
	content, err := io.ReadAll(&part)
	if err != nil {
		fmt.Println("Failed to read part body")
	}
	return content
}

// FindParts returns a list of Parts for which findFunc returns true
func FindParts(r *multipart.Reader, findFunc MultipartFindFunc) (parts []Part, err error) {
	for {
		p, err := r.NextPart()
		if err == io.EOF {
			// Reached the end successfully
			return parts, nil
		}
		if err != nil {
			return parts, err
		}
		if findFunc(p) {
			part := &Part{
				Part:    *p,
				Content: GetPartContent(*p),
			}
			parts = append(parts, *part)
		}
	}
}

// isMultipart checks if a media type is multipart/*
func isMultipart(mediaType string) bool {
	return strings.HasPrefix(mediaType, "multipart/")
}
