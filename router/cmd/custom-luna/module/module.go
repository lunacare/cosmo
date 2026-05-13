package module

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/wundergraph/cosmo/router/core"
	"go.uber.org/zap"
)

func init() {
	core.RegisterModule(&FileUploadEmptyVarModule{})
}

const FileUploadEmptyVarModuleID = "com.getluna.file-upload-empty-var"

// multipartReadFormMemory is the in-memory threshold for multipart.Reader.ReadForm.
// Parts smaller than this stay in memory; larger parts spill to disk. 32 MiB matches
// the stdlib default and keeps small uploads off the filesystem.
const multipartReadFormMemory = 32 << 20

type FileUploadEmptyVarModule struct {
	Logger *zap.Logger
}

func (m *FileUploadEmptyVarModule) Provision(ctx *core.ModuleContext) error {
	m.Logger = ctx.Logger
	return nil
}

func (m *FileUploadEmptyVarModule) RouterOnRequest(ctx core.RequestContext, next http.Handler) {
	r := ctx.Request()
	w := ctx.ResponseWriter()

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		next.ServeHTTP(w, r)
		return
	}

	mediatype, params, err := mime.ParseMediaType(contentType)
	if err != nil || mediatype != "multipart/form-data" {
		next.ServeHTTP(w, r)
		return
	}

	boundary, ok := params["boundary"]
	if !ok {
		core.WriteResponseError(ctx, fmt.Errorf("no boundary found in multipart form"))
		return
	}

	form, err := multipart.NewReader(r.Body, boundary).ReadForm(multipartReadFormMemory)
	if err != nil {
		core.WriteResponseError(ctx, fmt.Errorf("error reading multipart form: %w", err))
		return
	}
	defer form.RemoveAll()

	operations, hasOperations := form.Value["operations"]
	maps, hasMap := form.Value["map"]

	var (
		rewrittenOperations []byte
		didRewrite          bool
	)

	if hasOperations && len(operations) > 0 && hasMap && len(maps) > 0 {
		rewrittenOperations = []byte(operations[0])
		jsonMap := []byte(maps[0])

		var setErr error
		mapErr := jsonparser.ObjectEach(jsonMap, func(_, paths []byte, dataType jsonparser.ValueType, _ int) error {
			if dataType != jsonparser.Array {
				return nil
			}
			_, arrErr := jsonparser.ArrayEach(paths, func(path []byte, pathType jsonparser.ValueType, _ int, cbErr error) {
				if cbErr != nil || setErr != nil || pathType != jsonparser.String {
					return
				}
				pathElements := strings.Split(string(path), ".")
				rewrittenOperations, setErr = jsonparser.Set(rewrittenOperations, []byte("null"), pathElements...)
			})
			if arrErr != nil {
				return arrErr
			}
			return setErr
		})
		if err := firstNonNil(mapErr, setErr); err != nil {
			core.WriteResponseError(ctx, fmt.Errorf("error normalizing multipart map: %w", err))
			return
		}
		didRewrite = true
	}

	var body bytes.Buffer
	parsedMultipart := multipart.NewWriter(&body)

	for key, values := range form.Value {
		if key == "operations" && didRewrite {
			if err := parsedMultipart.WriteField("operations", string(rewrittenOperations)); err != nil {
				core.WriteResponseError(ctx, fmt.Errorf("error writing operations field: %w", err))
				return
			}
			continue
		}
		for _, value := range values {
			if err := parsedMultipart.WriteField(key, value); err != nil {
				core.WriteResponseError(ctx, fmt.Errorf("error writing field %s: %w", key, err))
				return
			}
		}
	}

	for key, files := range form.File {
		for _, fileHeader := range files {
			if err := copyFormFile(parsedMultipart, key, fileHeader); err != nil {
				core.WriteResponseError(ctx, err)
				return
			}
		}
	}

	if err := parsedMultipart.Close(); err != nil {
		core.WriteResponseError(ctx, fmt.Errorf("error finalizing multipart body: %w", err))
		return
	}

	r.Header.Set("Content-Type", parsedMultipart.FormDataContentType())
	r.Body = io.NopCloser(&body)
	r.ContentLength = int64(body.Len())

	next.ServeHTTP(w, r)
}

func copyFormFile(mw *multipart.Writer, field string, fh *multipart.FileHeader) error {
	fileWriter, err := mw.CreateFormFile(field, fh.Filename)
	if err != nil {
		return fmt.Errorf("error creating form file %s: %w", fh.Filename, err)
	}
	fileReader, err := fh.Open()
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", fh.Filename, err)
	}
	defer fileReader.Close()
	if _, err := io.Copy(fileWriter, fileReader); err != nil {
		return fmt.Errorf("error copying file %s: %w", fh.Filename, err)
	}
	return nil
}

func firstNonNil(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

func (m *FileUploadEmptyVarModule) Module() core.ModuleInfo {
	return core.ModuleInfo{
		ID: FileUploadEmptyVarModuleID,
		New: func() core.Module {
			return &FileUploadEmptyVarModule{}
		},
	}
}

var (
	_ core.RouterOnRequestHandler = (*FileUploadEmptyVarModule)(nil)
	_ core.Provisioner            = (*FileUploadEmptyVarModule)(nil)
)
