package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/document-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock store
// ---------------------------------------------------------------------------

type mockStore struct {
	createDocumentFn  func(ctx context.Context, tenantID uuid.UUID, d *domain.Document) error
	getDocumentFn     func(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Document, error)
	listDocumentsFn   func(ctx context.Context, tenantID uuid.UUID) ([]domain.Document, error)
	updateOCRStatusFn func(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, result string) error
	deleteDocumentFn  func(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error
}

func (m *mockStore) CreateDocument(ctx context.Context, tenantID uuid.UUID, d *domain.Document) error {
	if m.createDocumentFn != nil {
		return m.createDocumentFn(ctx, tenantID, d)
	}
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	d.TenantID = tenantID
	d.VirusStatus = "pending"
	d.OCRStatus = "none"
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) GetDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Document, error) {
	if m.getDocumentFn != nil {
		return m.getDocumentFn(ctx, tenantID, id)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListDocuments(ctx context.Context, tenantID uuid.UUID) ([]domain.Document, error) {
	if m.listDocumentsFn != nil {
		return m.listDocumentsFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) UpdateOCRStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, result string) error {
	if m.updateOCRStatusFn != nil {
		return m.updateOCRStatusFn(ctx, tenantID, id, status, result)
	}
	return nil
}

func (m *mockStore) DeleteDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error {
	if m.deleteDocumentFn != nil {
		return m.deleteDocumentFn(ctx, tenantID, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Mock KMS
// ---------------------------------------------------------------------------

type mockKMS struct {
	GenerateDataKeyFn func(ctx context.Context, params *kms.GenerateDataKeyInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error)
	DecryptFn         func(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

func (m *mockKMS) GenerateDataKey(ctx context.Context, params *kms.GenerateDataKeyInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error) {
	if m.GenerateDataKeyFn != nil {
		return m.GenerateDataKeyFn(ctx, params, optFns...)
	}
	// Generate a real 32-byte key for AES-256
	plaintextKey := make([]byte, 32)
	rand.Read(plaintextKey)
	cipherBlob := make([]byte, 64)
	rand.Read(cipherBlob)
	return &kms.GenerateDataKeyOutput{
		Plaintext:      plaintextKey,
		CiphertextBlob: cipherBlob,
		KeyId:          params.KeyId,
	}, nil
}

func (m *mockKMS) Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	if m.DecryptFn != nil {
		return m.DecryptFn(ctx, params, optFns...)
	}
	return nil, errors.New("decrypt not configured")
}

// ---------------------------------------------------------------------------
// Mock S3
// ---------------------------------------------------------------------------

type mockS3 struct {
	PutObjectFn func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObjectFn func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func (m *mockS3) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if m.PutObjectFn != nil {
		return m.PutObjectFn(ctx, params, optFns...)
	}
	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if m.GetObjectFn != nil {
		return m.GetObjectFn(ctx, params, optFns...)
	}
	return nil, errors.New("not configured")
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func buildMultipartRequest(t *testing.T, tenantID uuid.UUID, fileContent []byte, fileName, documentType, kmsKeyARN, documentNumber string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", fileName)
	require.NoError(t, err)
	_, err = part.Write(fileContent)
	require.NoError(t, err)

	require.NoError(t, writer.WriteField("document_type", documentType))
	require.NoError(t, writer.WriteField("kms_key_arn", kmsKeyARN))
	if documentNumber != "" {
		require.NoError(t, writer.WriteField("document_number", documentNumber))
	}

	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Tenant-Id", tenantID.String())
	return req
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil, nil, "test-bucket")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "document-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: Upload with DEK encryption round-trip
// ---------------------------------------------------------------------------

func TestUpload_Success(t *testing.T) {
	tenantID := uuid.New()
	kmsKeyARN := "arn:aws:kms:ap-south-1:000000000000:key/" + uuid.New().String()
	originalContent := []byte("hello world - this is a test document")

	// Known 32-byte key for testing
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i)
	}
	testCipherBlob := []byte("encrypted-dek-blob")

	var capturedS3Body []byte
	var capturedS3Key string
	var capturedEncryptedDEK []byte

	mk := &mockKMS{
		GenerateDataKeyFn: func(ctx context.Context, params *kms.GenerateDataKeyInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error) {
			assert.Equal(t, kmsKeyARN, *params.KeyId)
			return &kms.GenerateDataKeyOutput{
				Plaintext:      append([]byte{}, testKey...),
				CiphertextBlob: testCipherBlob,
				KeyId:          &kmsKeyARN,
			}, nil
		},
	}

	ms3 := &mockS3{
		PutObjectFn: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			body, err := io.ReadAll(params.Body)
			require.NoError(t, err)
			capturedS3Body = body
			capturedS3Key = *params.Key
			assert.Equal(t, "test-bucket", *params.Bucket)
			return &s3.PutObjectOutput{}, nil
		},
	}

	mst := &mockStore{
		createDocumentFn: func(ctx context.Context, tid uuid.UUID, d *domain.Document) error {
			assert.Equal(t, tenantID, tid)
			capturedEncryptedDEK = d.EncryptedDEK
			d.TenantID = tid
			d.VirusStatus = "pending"
			d.OCRStatus = "none"
			d.CreatedAt = time.Now()
			d.UpdatedAt = time.Now()
			return nil
		},
	}

	h := NewHandlers(mst, mk, ms3, "test-bucket")

	req := buildMultipartRequest(t, tenantID, originalContent, "invoice.pdf", "invoice", kmsKeyARN, "INV-001")
	rec := httptest.NewRecorder()

	h.Upload(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	// Verify S3 received encrypted content (not the original)
	assert.NotEqual(t, originalContent, capturedS3Body)
	assert.NotEmpty(t, capturedS3Body)
	assert.Contains(t, capturedS3Key, tenantID.String())

	// Verify the encrypted DEK stored in DB matches what KMS returned
	assert.Equal(t, testCipherBlob, capturedEncryptedDEK)

	// Verify the encrypted content can be decrypted with the test key
	decrypted, err := decryptAESGCM(capturedS3Body, testKey)
	require.NoError(t, err)
	assert.Equal(t, originalContent, decrypted)

	// Verify response
	var doc domain.Document
	parseDataResponse(t, rec.Body.Bytes(), &doc)
	assert.Equal(t, "invoice", doc.DocumentType)
	assert.Equal(t, "invoice.pdf", doc.FileName)
	assert.Equal(t, int64(len(originalContent)), doc.FileSize)
	assert.Equal(t, "AES-256-GCM", doc.EncryptionAlgo)
}

func TestUpload_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{}, &mockKMS{}, &mockS3{}, "test-bucket")
	req := httptest.NewRequest(http.MethodPost, "/v1/documents/upload", nil)
	rec := httptest.NewRecorder()

	h.Upload(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpload_MissingFile(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, &mockKMS{}, &mockS3{}, "test-bucket")

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("document_type", "invoice")
	writer.WriteField("kms_key_arn", "arn:aws:kms:ap-south-1:000:key/test")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.Upload(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "file is required", data["error"])
}

func TestUpload_MissingDocumentType(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, &mockKMS{}, &mockS3{}, "test-bucket")

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "test.pdf")
	part.Write([]byte("content"))
	writer.WriteField("kms_key_arn", "arn:aws:kms:ap-south-1:000:key/test")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.Upload(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "document_type is required", data["error"])
}

func TestUpload_MissingKMSKeyARN(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, &mockKMS{}, &mockS3{}, "test-bucket")

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "test.pdf")
	part.Write([]byte("content"))
	writer.WriteField("document_type", "invoice")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.Upload(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "kms_key_arn is required", data["error"])
}

func TestUpload_KMSGenerateDataKeyError(t *testing.T) {
	tenantID := uuid.New()
	mk := &mockKMS{
		GenerateDataKeyFn: func(ctx context.Context, params *kms.GenerateDataKeyInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error) {
			return nil, errors.New("kms unavailable")
		},
	}
	h := NewHandlers(&mockStore{}, mk, &mockS3{}, "test-bucket")

	req := buildMultipartRequest(t, tenantID, []byte("content"), "test.pdf", "invoice", "arn:aws:kms:test", "")
	rec := httptest.NewRecorder()

	h.Upload(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "encryption key generation failed", data["error"])
}

func TestUpload_S3PutObjectError(t *testing.T) {
	tenantID := uuid.New()
	mk := &mockKMS{}
	ms3 := &mockS3{
		PutObjectFn: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			return nil, errors.New("s3 down")
		},
	}
	h := NewHandlers(&mockStore{}, mk, ms3, "test-bucket")

	req := buildMultipartRequest(t, tenantID, []byte("content"), "test.pdf", "invoice", "arn:aws:kms:test", "")
	rec := httptest.NewRecorder()

	h.Upload(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "upload to storage failed", data["error"])
}

func TestUpload_StoreError(t *testing.T) {
	tenantID := uuid.New()
	mst := &mockStore{
		createDocumentFn: func(ctx context.Context, tid uuid.UUID, d *domain.Document) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(mst, &mockKMS{}, &mockS3{}, "test-bucket")

	req := buildMultipartRequest(t, tenantID, []byte("content"), "test.pdf", "invoice", "arn:aws:kms:test", "")
	rec := httptest.NewRecorder()

	h.Upload(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "store metadata failed", data["error"])
}

// ---------------------------------------------------------------------------
// Tests: Download with DEK decryption round-trip
// ---------------------------------------------------------------------------

func TestDownload_Success(t *testing.T) {
	tenantID := uuid.New()
	documentID := uuid.New()
	originalContent := []byte("this is confidential document content")

	// Use a known key
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i + 10)
	}
	testCipherBlob := []byte("wrapped-dek")

	// Encrypt content with the test key
	encryptedContent, err := encryptAESGCM(originalContent, testKey)
	require.NoError(t, err)

	kmsKeyARN := "arn:aws:kms:ap-south-1:000:key/test"
	mst := &mockStore{
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, documentID, id)
			return &domain.Document{
				ID:             documentID,
				TenantID:       tenantID,
				DocumentType:   "invoice",
				FileName:       "invoice.pdf",
				MimeType:       "application/pdf",
				FileSize:       int64(len(originalContent)),
				S3Bucket:       "test-bucket",
				S3Key:          "tenants/" + tenantID.String() + "/documents/" + documentID.String() + "/invoice.pdf",
				EncryptedDEK:   testCipherBlob,
				KMSKeyARN:      &kmsKeyARN,
				EncryptionAlgo: "AES-256-GCM",
				VirusStatus:    "pending",
				OCRStatus:      "none",
				Tags:           "[]",
				Metadata:       "{}",
			}, nil
		},
	}

	mk := &mockKMS{
		DecryptFn: func(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
			assert.Equal(t, testCipherBlob, params.CiphertextBlob)
			return &kms.DecryptOutput{
				Plaintext: append([]byte{}, testKey...),
			}, nil
		},
	}

	ms3 := &mockS3{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body:          io.NopCloser(bytes.NewReader(encryptedContent)),
				ContentLength: ptrInt64(int64(len(encryptedContent))),
			}, nil
		},
	}

	h := NewHandlers(mst, mk, ms3, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+documentID.String()+"/download", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", documentID.String())
	rec := httptest.NewRecorder()

	h.Download(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/pdf", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "invoice.pdf")
	assert.Equal(t, originalContent, rec.Body.Bytes())
}

func TestDownload_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{}, &mockKMS{}, &mockS3{}, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+uuid.New().String()+"/download", nil)
	req.SetPathValue("documentID", uuid.New().String())
	rec := httptest.NewRecorder()

	h.Download(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDownload_InvalidDocumentID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, &mockKMS{}, &mockS3{}, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/bad/download", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", "bad")
	rec := httptest.NewRecorder()

	h.Download(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDownload_DocumentNotFound(t *testing.T) {
	tenantID := uuid.New()
	docID := uuid.New()
	h := NewHandlers(&mockStore{}, &mockKMS{}, &mockS3{}, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+docID.String()+"/download", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", docID.String())
	rec := httptest.NewRecorder()

	h.Download(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDownload_S3GetObjectError(t *testing.T) {
	tenantID := uuid.New()
	documentID := uuid.New()
	kmsKeyARN := "arn:aws:kms:test"

	mst := &mockStore{
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			return &domain.Document{
				ID: documentID, TenantID: tenantID, S3Bucket: "test-bucket", S3Key: "key",
				EncryptedDEK: []byte("dek"), KMSKeyARN: &kmsKeyARN, MimeType: "application/pdf",
				FileName: "test.pdf", EncryptionAlgo: "AES-256-GCM", Tags: "[]", Metadata: "{}",
			}, nil
		},
	}

	ms3 := &mockS3{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return nil, errors.New("s3 error")
		},
	}

	h := NewHandlers(mst, &mockKMS{}, ms3, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+documentID.String()+"/download", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", documentID.String())
	rec := httptest.NewRecorder()

	h.Download(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "fetch from storage failed", data["error"])
}

func TestDownload_KMSDecryptError(t *testing.T) {
	tenantID := uuid.New()
	documentID := uuid.New()
	kmsKeyARN := "arn:aws:kms:test"

	mst := &mockStore{
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			return &domain.Document{
				ID: documentID, TenantID: tenantID, S3Bucket: "test-bucket", S3Key: "key",
				EncryptedDEK: []byte("dek"), KMSKeyARN: &kmsKeyARN, MimeType: "application/pdf",
				FileName: "test.pdf", EncryptionAlgo: "AES-256-GCM", Tags: "[]", Metadata: "{}",
			}, nil
		},
	}

	ms3 := &mockS3{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader([]byte("encrypted"))),
			}, nil
		},
	}

	mk := &mockKMS{
		DecryptFn: func(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
			return nil, errors.New("kms decrypt failed")
		},
	}

	h := NewHandlers(mst, mk, ms3, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+documentID.String()+"/download", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", documentID.String())
	rec := httptest.NewRecorder()

	h.Download(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "decryption key recovery failed", data["error"])
}

func TestDownload_DecryptionError(t *testing.T) {
	tenantID := uuid.New()
	documentID := uuid.New()
	kmsKeyARN := "arn:aws:kms:test"

	mst := &mockStore{
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			return &domain.Document{
				ID: documentID, TenantID: tenantID, S3Bucket: "test-bucket", S3Key: "key",
				EncryptedDEK: []byte("dek"), KMSKeyARN: &kmsKeyARN, MimeType: "application/pdf",
				FileName: "test.pdf", EncryptionAlgo: "AES-256-GCM", Tags: "[]", Metadata: "{}",
			}, nil
		},
	}

	// Return garbage that cannot be decrypted
	ms3 := &mockS3{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader([]byte("this-is-not-valid-encrypted-content-needs-to-be-longer-than-nonce"))),
			}, nil
		},
	}

	wrongKey := make([]byte, 32)
	for i := range wrongKey {
		wrongKey[i] = byte(i + 99)
	}

	mk := &mockKMS{
		DecryptFn: func(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
			return &kms.DecryptOutput{
				Plaintext: wrongKey,
			}, nil
		},
	}

	h := NewHandlers(mst, mk, ms3, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+documentID.String()+"/download", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", documentID.String())
	rec := httptest.NewRecorder()

	h.Download(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "file decryption failed", data["error"])
}

// ---------------------------------------------------------------------------
// Tests: GetDocument
// ---------------------------------------------------------------------------

func TestGetDocument_Success(t *testing.T) {
	tenantID := uuid.New()
	documentID := uuid.New()
	kmsKeyARN := "arn:aws:kms:test"

	mst := &mockStore{
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, documentID, id)
			return &domain.Document{
				ID: documentID, TenantID: tenantID, DocumentType: "invoice",
				FileName: "test.pdf", MimeType: "application/pdf", FileSize: 1024,
				S3Bucket: "bucket", S3Key: "key", EncryptionAlgo: "AES-256-GCM",
				VirusStatus: "pending", OCRStatus: "none", Tags: "[]", Metadata: "{}",
				KMSKeyARN: &kmsKeyARN,
			}, nil
		},
	}

	h := NewHandlers(mst, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+documentID.String(), nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", documentID.String())
	rec := httptest.NewRecorder()

	h.GetDocument(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var doc domain.Document
	parseDataResponse(t, rec.Body.Bytes(), &doc)
	assert.Equal(t, documentID, doc.ID)
	assert.Equal(t, "invoice", doc.DocumentType)
	assert.Equal(t, "test.pdf", doc.FileName)
}

func TestGetDocument_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/bad", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", "bad")
	rec := httptest.NewRecorder()

	h.GetDocument(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid document_id", data["error"])
}

func TestGetDocument_NotFound(t *testing.T) {
	tenantID := uuid.New()
	docID := uuid.New()
	mst := &mockStore{
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			return nil, errors.New("not found")
		},
	}

	h := NewHandlers(mst, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+docID.String(), nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", docID.String())
	rec := httptest.NewRecorder()

	h.GetDocument(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetDocument_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents/"+uuid.New().String(), nil)
	req.SetPathValue("documentID", uuid.New().String())
	rec := httptest.NewRecorder()

	h.GetDocument(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListDocuments
// ---------------------------------------------------------------------------

func TestListDocuments_Success(t *testing.T) {
	tenantID := uuid.New()
	docID := uuid.New()

	mst := &mockStore{
		listDocumentsFn: func(ctx context.Context, tid uuid.UUID) ([]domain.Document, error) {
			assert.Equal(t, tenantID, tid)
			return []domain.Document{
				{ID: docID, TenantID: tenantID, DocumentType: "invoice", FileName: "test.pdf",
					MimeType: "application/pdf", FileSize: 1024, S3Bucket: "bucket", S3Key: "key",
					EncryptionAlgo: "AES-256-GCM", VirusStatus: "pending", OCRStatus: "none",
					Tags: "[]", Metadata: "{}"},
			}, nil
		},
	}

	h := NewHandlers(mst, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListDocuments(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var docs []domain.Document
	parseDataResponse(t, rec.Body.Bytes(), &docs)
	require.Len(t, docs, 1)
	assert.Equal(t, "invoice", docs[0].DocumentType)
}

func TestListDocuments_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents", nil)
	rec := httptest.NewRecorder()

	h.ListDocuments(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListDocuments_StoreError(t *testing.T) {
	tenantID := uuid.New()
	mst := &mockStore{
		listDocumentsFn: func(ctx context.Context, tid uuid.UUID) ([]domain.Document, error) {
			return nil, errors.New("db error")
		},
	}

	h := NewHandlers(mst, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListDocuments(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListDocuments_NilResultReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	mst := &mockStore{
		listDocumentsFn: func(ctx context.Context, tid uuid.UUID) ([]domain.Document, error) {
			return nil, nil
		},
	}

	h := NewHandlers(mst, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodGet, "/v1/documents", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListDocuments(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var docs []domain.Document
	parseDataResponse(t, rec.Body.Bytes(), &docs)
	assert.NotNil(t, docs)
	assert.Len(t, docs, 0)
}

// ---------------------------------------------------------------------------
// Tests: TriggerOCR
// ---------------------------------------------------------------------------

func TestTriggerOCR_Success(t *testing.T) {
	tenantID := uuid.New()
	documentID := uuid.New()

	var capturedStatus string
	mst := &mockStore{
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			return &domain.Document{ID: documentID, TenantID: tenantID}, nil
		},
		updateOCRStatusFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID, status string, result string) error {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, documentID, id)
			capturedStatus = status
			assert.Equal(t, "", result)
			return nil
		},
	}

	h := NewHandlers(mst, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/"+documentID.String()+"/ocr", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", documentID.String())
	rec := httptest.NewRecorder()

	h.TriggerOCR(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "queued", capturedStatus)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "queued", data["status"])
	assert.Equal(t, documentID.String(), data["document_id"])
}

func TestTriggerOCR_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/"+uuid.New().String()+"/ocr", nil)
	req.SetPathValue("documentID", uuid.New().String())
	rec := httptest.NewRecorder()

	h.TriggerOCR(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTriggerOCR_InvalidDocumentID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/bad/ocr", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", "bad")
	rec := httptest.NewRecorder()

	h.TriggerOCR(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTriggerOCR_DocumentNotFound(t *testing.T) {
	tenantID := uuid.New()
	docID := uuid.New()

	h := NewHandlers(&mockStore{}, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/"+docID.String()+"/ocr", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", docID.String())
	rec := httptest.NewRecorder()

	h.TriggerOCR(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTriggerOCR_StoreUpdateError(t *testing.T) {
	tenantID := uuid.New()
	documentID := uuid.New()

	mst := &mockStore{
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			return &domain.Document{ID: documentID, TenantID: tenantID}, nil
		},
		updateOCRStatusFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID, status string, result string) error {
			return errors.New("db error")
		},
	}

	h := NewHandlers(mst, nil, nil, "test-bucket")

	req := httptest.NewRequest(http.MethodPost, "/v1/documents/"+documentID.String()+"/ocr", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("documentID", documentID.String())
	rec := httptest.NewRecorder()

	h.TriggerOCR(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: tenantIDFromRequest
// ---------------------------------------------------------------------------

func TestTenantIDFromRequest_Valid(t *testing.T) {
	expected := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", expected.String())

	got, err := tenantIDFromRequest(req)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestTenantIDFromRequest_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing X-Tenant-Id")
}

func TestTenantIDFromRequest_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", "not-uuid")

	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Tests: AES-256-GCM encrypt/decrypt helpers
// ---------------------------------------------------------------------------

func TestEncryptDecryptAESGCM_RoundTrip(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	plaintext := []byte("sensitive data for round-trip test")

	encrypted, err := encryptAESGCM(plaintext, key)
	require.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := decryptAESGCM(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDecryptAESGCM_TooShort(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	_, err := decryptAESGCM([]byte("short"), key)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext too short")
}

func TestEncryptAESGCM_InvalidKeySize(t *testing.T) {
	_, err := encryptAESGCM([]byte("data"), []byte("short-key"))
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Tests: NewRouter
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	ms := &mockStore{}
	r := NewRouter(ms, nil, nil, "test-bucket")
	require.NotNil(t, r)

	// Verify health endpoint is reachable
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify ping heartbeat
	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Upload + Download full encryption round-trip
// ---------------------------------------------------------------------------

func TestUploadDownload_FullRoundTrip(t *testing.T) {
	tenantID := uuid.New()
	kmsKeyARN := "arn:aws:kms:ap-south-1:000000000000:key/" + uuid.New().String()
	originalContent := []byte("CONFIDENTIAL: This is the full round-trip test content with special chars: @#$%^&*()")

	// Known key for both encrypt and decrypt
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i * 3)
	}
	testCipherBlob := []byte("encrypted-dek-for-round-trip")

	// Storage to capture S3 and DB data
	var storedS3Content []byte
	var storedDocument *domain.Document

	mk := &mockKMS{
		GenerateDataKeyFn: func(ctx context.Context, params *kms.GenerateDataKeyInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error) {
			return &kms.GenerateDataKeyOutput{
				Plaintext:      append([]byte{}, testKey...),
				CiphertextBlob: testCipherBlob,
				KeyId:          &kmsKeyARN,
			}, nil
		},
		DecryptFn: func(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
			assert.Equal(t, testCipherBlob, params.CiphertextBlob)
			return &kms.DecryptOutput{
				Plaintext: append([]byte{}, testKey...),
			}, nil
		},
	}

	ms3 := &mockS3{
		PutObjectFn: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			body, _ := io.ReadAll(params.Body)
			storedS3Content = body
			return &s3.PutObjectOutput{}, nil
		},
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader(storedS3Content)),
			}, nil
		},
	}

	mst := &mockStore{
		createDocumentFn: func(ctx context.Context, tid uuid.UUID, d *domain.Document) error {
			d.TenantID = tid
			d.VirusStatus = "pending"
			d.OCRStatus = "none"
			d.CreatedAt = time.Now()
			d.UpdatedAt = time.Now()
			storedDocument = d
			return nil
		},
		getDocumentFn: func(ctx context.Context, tid uuid.UUID, id uuid.UUID) (*domain.Document, error) {
			if storedDocument != nil {
				return storedDocument, nil
			}
			return nil, errors.New("not found")
		},
	}

	h := NewHandlers(mst, mk, ms3, "test-bucket")

	// Step 1: Upload
	uploadReq := buildMultipartRequest(t, tenantID, originalContent, "secret.pdf", "contract", kmsKeyARN, "CTR-001")
	uploadRec := httptest.NewRecorder()
	h.Upload(uploadRec, uploadReq)
	assert.Equal(t, http.StatusCreated, uploadRec.Code)

	var uploadedDoc domain.Document
	parseDataResponse(t, uploadRec.Body.Bytes(), &uploadedDoc)

	// Verify S3 content is encrypted (different from original)
	assert.NotEqual(t, originalContent, storedS3Content)

	// Step 2: Download
	downloadReq := httptest.NewRequest(http.MethodGet, "/v1/documents/"+uploadedDoc.ID.String()+"/download", nil)
	downloadReq.Header.Set("X-Tenant-Id", tenantID.String())
	downloadReq.SetPathValue("documentID", uploadedDoc.ID.String())
	downloadRec := httptest.NewRecorder()
	h.Download(downloadRec, downloadReq)

	assert.Equal(t, http.StatusOK, downloadRec.Code)
	assert.Equal(t, originalContent, downloadRec.Body.Bytes())
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func ptrInt64(v int64) *int64 { return &v }

// Ensure s3Types is used (avoid unused import)
var _ = s3Types.ObjectStorageClassStandard
