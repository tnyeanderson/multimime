package multimime

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

type MultipartFindFunc func(*multipart.Part) bool

type Part struct {
	multipart.Part
	Content []byte
}

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

func GetAllParts(r io.Reader) ([]Part, error) {
	return GetParts(r, func(_ *multipart.Part) bool { return true })
}

func GetTextParts(r io.Reader) ([]Part, error) {
	return GetParts(r, IsTextPart)
}

func GetPlainTextParts(r io.Reader) ([]Part, error) {
	return GetParts(r, IsPlainTextPart)
}

func GetHtmlParts(r io.Reader) ([]Part, error) {
	return GetParts(r, IsHtmlPart)
}

func GetAttachments(r io.Reader) ([]Part, error) {
	return GetParts(r, IsAttachment)
}

func IsPlainTextPart(part *multipart.Part) bool {
	mediaType := GetPartType(part)
	return strings.HasPrefix(mediaType, "text/plain")
}

func IsHtmlPart(part *multipart.Part) bool {
	mediaType := GetPartType(part)
	return strings.HasPrefix(mediaType, "text/html")
}

func IsTextPart(part *multipart.Part) bool {
	mediaType := GetPartType(part)
	return strings.HasPrefix(mediaType, "text/")
}

func IsAttachment(part *multipart.Part) bool {
	disposition := GetPartDisposition(part)
	return strings.HasPrefix(disposition, "attachment")
}

func GetPartDisposition(part *multipart.Part) string {
	header := part.Header.Get("Content-Disposition")
	disposition, _, err := mime.ParseMediaType(header)
	if err != nil {
		return ""
	}
	return disposition
}

func GetPartType(part *multipart.Part) string {
	header := part.Header.Get("Content-Type")
	contentType, _, err := mime.ParseMediaType(header)
	if err != nil {
		return ""
	}
	return contentType
}

func GetMessageType(message *mail.Message) (mediatype string, params map[string]string, err error) {
	contentType := message.Header.Get("Content-Type")
	return mime.ParseMediaType(contentType)
}

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

func GetPartContent(part multipart.Part) (content []byte) {
	content, err := io.ReadAll(&part)
	if err != nil {
		fmt.Println("Failed to read part body")
	}
	return content
}

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

func isMultipart(mediaType string) bool {
	return strings.HasPrefix(mediaType, "multipart/")
}

//func readPart(p *multipart.Part) {
//	mediaType, params := getMediaType(p.Header)
//
//	if isMultipart(mediaType) {
//		boundary := params["boundary"]
//
//		fmt.Printf("Reading multipart with boundary %q\n", boundary)
//
//		pr := multipart.NewReader(p, boundary)
//		readMultipart(pr)
//		return
//	} else if mediaType == "message/rfc822" {
//		fmt.Println("Part is attached email")
//		readMsg(p)
//	}
//
//	fmt.Printf("Not multipart: %s\n", p.Header.Get("Content-Type"))
//}
