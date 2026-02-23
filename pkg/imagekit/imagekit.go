package imagekit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	ik "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

type ImageKitService interface {
	UploadFile(ctx context.Context, file io.Reader, fileName string, folder string) (*ik.FileUploadResponse, error)
}

type imageKitService struct {
	client      *ik.Client
	publicKey   string
	urlEndpoint string
}

func NewImageKitService(publicKey, privateKey, urlEndpoint string) (ImageKitService, error) {
	if strings.TrimSpace(privateKey) == "" {
		return nil, errors.New("imagekit privateKey is required (server-side upload uses private key)")
	}

	client := ik.NewClient(
		option.WithPrivateKey(privateKey),
	)

	return &imageKitService{
		client:      &client,
		publicKey:   publicKey,
		urlEndpoint: urlEndpoint,
	}, nil
}

func (s *imageKitService) UploadFile(ctx context.Context, file io.Reader, fileName string, folder string) (*ik.FileUploadResponse, error) {
	// Normalize folder (optional)
	folder = normalizeFolder(folder)

	// Pastikan ada content-type (sniff dari 512 bytes pertama)
	rdr, contentType, err := ensureContentType(file, "")
	if err != nil {
		return nil, err
	}

	payload := ik.NewFile(rdr, fileName, contentType)

	// Generate unique filename
	fileName = fmt.Sprintf("%d-%s", time.Now().Unix(), fileName)

	params := ik.FileUploadParams{
		File:     payload,
		FileName: fileName,

		// Default behavior yang biasanya enak untuk backend:
		UseUniqueFileName: ik.Bool(true), // biar ga ketiban kalau nama sama
	}

	if folder != "" {
		params.Folder = ik.String(folder)
	}

	resp, err := s.client.Files.Upload(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to imagekit: %w", err)
	}

	return resp, nil
}

func normalizeFolder(folder string) string {
	folder = strings.TrimSpace(folder)
	if folder == "" {
		return ""
	}
	// Biar konsisten, kasih leading slash
	if !strings.HasPrefix(folder, "/") {
		folder = "/" + folder
	}
	// Optional: buang trailing slash (ImageKit biasanya OK, tapi rapihin aja)
	folder = strings.TrimRight(folder, "/")
	return folder
}

// ensureContentType: kalau kosong, sniff 512 bytes pertama dan kembalikan reader yang "utuh lagi".
func ensureContentType(r io.Reader, current string) (io.Reader, string, error) {
	if strings.TrimSpace(current) != "" {
		return r, current, nil
	}

	head, err := io.ReadAll(io.LimitReader(r, 512))
	if err != nil {
		return nil, "", fmt.Errorf("failed to read for content-type sniff: %w", err)
	}

	ct := http.DetectContentType(head)
	// Balikin reader: head + sisanya
	return io.MultiReader(bytes.NewReader(head), r), ct, nil
}
