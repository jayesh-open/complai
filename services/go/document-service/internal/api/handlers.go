package api

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/document-service/internal/domain"
	"github.com/complai/complai/services/go/document-service/internal/store"
)

// KMSClient is the subset of the AWS KMS client used by Handlers.
type KMSClient interface {
	GenerateDataKey(ctx context.Context, params *kms.GenerateDataKeyInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error)
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

// S3Client is the subset of the AWS S3 client used by Handlers.
type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Handlers struct {
	store     store.Repository
	kmsClient KMSClient
	s3Client  S3Client
	bucket    string
}

func NewHandlers(s store.Repository, kmsClient KMSClient, s3Client S3Client, bucket string) *Handlers {
	return &Handlers{store: s, kmsClient: kmsClient, s3Client: s3Client, bucket: bucket}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "document-service"})
}

func (h *Handlers) Upload(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Parse multipart form (max 32 MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart form"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "file is required"})
		return
	}
	defer file.Close()

	documentType := r.FormValue("document_type")
	if documentType == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "document_type is required"})
		return
	}

	kmsKeyARN := r.FormValue("kms_key_arn")
	if kmsKeyARN == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "kms_key_arn is required"})
		return
	}

	documentNumber := r.FormValue("document_number")

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("read file failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "read file failed"})
		return
	}

	// Generate DEK using KMS
	dataKeyOutput, err := h.kmsClient.GenerateDataKey(r.Context(), &kms.GenerateDataKeyInput{
		KeyId:   &kmsKeyARN,
		KeySpec: kmsTypes.DataKeySpecAes256,
	})
	if err != nil {
		log.Error().Err(err).Msg("KMS GenerateDataKey failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "encryption key generation failed"})
		return
	}

	// Encrypt file content with plaintext DEK
	encryptedContent, err := encryptAESGCM(fileContent, dataKeyOutput.Plaintext)
	if err != nil {
		log.Error().Err(err).Msg("encryption failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "encryption failed"})
		return
	}

	// Zero out plaintext key from memory
	for i := range dataKeyOutput.Plaintext {
		dataKeyOutput.Plaintext[i] = 0
	}

	// Build S3 key
	docID := uuid.New()
	s3Key := fmt.Sprintf("tenants/%s/documents/%s/%s", tenantID.String(), docID.String(), header.Filename)

	// Upload encrypted content to S3
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	_, err = h.s3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      &h.bucket,
		Key:         &s3Key,
		Body:        bytes.NewReader(encryptedContent),
		ContentType: &contentType,
	})
	if err != nil {
		log.Error().Err(err).Msg("S3 PutObject failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "upload to storage failed"})
		return
	}

	// Store metadata in DB
	var docNumber *string
	if documentNumber != "" {
		docNumber = &documentNumber
	}

	doc := &domain.Document{
		ID:             docID,
		DocumentType:   documentType,
		DocumentNumber: docNumber,
		FileName:       header.Filename,
		MimeType:       contentType,
		FileSize:       int64(len(fileContent)),
		S3Bucket:       h.bucket,
		S3Key:          s3Key,
		EncryptedDEK:   dataKeyOutput.CiphertextBlob,
		KMSKeyARN:      &kmsKeyARN,
		EncryptionAlgo: "AES-256-GCM",
		Tags:           "[]",
		Metadata:       "{}",
	}

	if err := h.store.CreateDocument(r.Context(), tenantID, doc); err != nil {
		log.Error().Err(err).Msg("store document metadata failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "store metadata failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, doc)
}

func (h *Handlers) GetDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	documentID, err := uuid.Parse(r.PathValue("documentID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid document_id"})
		return
	}

	doc, err := h.store.GetDocument(r.Context(), tenantID, documentID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "document not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, doc)
}

func (h *Handlers) Download(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	documentID, err := uuid.Parse(r.PathValue("documentID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid document_id"})
		return
	}

	doc, err := h.store.GetDocument(r.Context(), tenantID, documentID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "document not found"})
		return
	}

	// Get encrypted file from S3
	s3Out, err := h.s3Client.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: &doc.S3Bucket,
		Key:    &doc.S3Key,
	})
	if err != nil {
		log.Error().Err(err).Msg("S3 GetObject failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "fetch from storage failed"})
		return
	}
	defer s3Out.Body.Close()

	encryptedContent, err := io.ReadAll(s3Out.Body)
	if err != nil {
		log.Error().Err(err).Msg("read S3 body failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "read storage content failed"})
		return
	}

	// Decrypt DEK via KMS
	decryptOut, err := h.kmsClient.Decrypt(r.Context(), &kms.DecryptInput{
		CiphertextBlob: doc.EncryptedDEK,
	})
	if err != nil {
		log.Error().Err(err).Msg("KMS Decrypt failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "decryption key recovery failed"})
		return
	}

	// Decrypt file content
	plaintext, err := decryptAESGCM(encryptedContent, decryptOut.Plaintext)
	if err != nil {
		log.Error().Err(err).Msg("file decryption failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "file decryption failed"})
		return
	}

	// Zero out plaintext key from memory
	for i := range decryptOut.Plaintext {
		decryptOut.Plaintext[i] = 0
	}

	w.Header().Set("Content-Type", doc.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, doc.FileName))
	w.WriteHeader(http.StatusOK)
	w.Write(plaintext)
}

func (h *Handlers) ListDocuments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	docs, err := h.store.ListDocuments(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list documents failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if docs == nil {
		docs = []domain.Document{}
	}

	httputil.JSON(w, http.StatusOK, docs)
}

func (h *Handlers) TriggerOCR(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	documentID, err := uuid.Parse(r.PathValue("documentID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid document_id"})
		return
	}

	// Verify document exists
	_, err = h.store.GetDocument(r.Context(), tenantID, documentID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "document not found"})
		return
	}

	if err := h.store.UpdateOCRStatus(r.Context(), tenantID, documentID, "queued", ""); err != nil {
		log.Error().Err(err).Msg("update OCR status failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "queued", "document_id": documentID.String()})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}

func encryptAESGCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptAESGCM(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// marshalJSON is a helper to marshal data to JSON bytes.
func marshalJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
