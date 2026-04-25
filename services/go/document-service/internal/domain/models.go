package domain

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	DocumentType   string     `json:"document_type"`
	DocumentNumber *string    `json:"document_number,omitempty"`
	FileName       string     `json:"file_name"`
	MimeType       string     `json:"mime_type"`
	FileSize       int64      `json:"file_size"`
	S3Bucket       string     `json:"s3_bucket"`
	S3Key          string     `json:"s3_key"`
	EncryptedDEK   []byte     `json:"-"`
	KMSKeyARN      *string    `json:"kms_key_arn,omitempty"`
	EncryptionAlgo string     `json:"encryption_algo"`
	VirusStatus    string     `json:"virus_status"`
	OCRStatus      string     `json:"ocr_status"`
	OCRResult      *string    `json:"ocr_result,omitempty"`
	Tags           string     `json:"tags"`
	Metadata       string     `json:"metadata"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type UploadRequest struct {
	DocumentType   string `json:"document_type"`
	DocumentNumber string `json:"document_number"`
	FileName       string `json:"file_name"`
	KMSKeyARN      string `json:"kms_key_arn"`
}

type DocumentLineage struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	ParentDocumentID uuid.UUID `json:"parent_document_id"`
	ChildDocumentID  uuid.UUID `json:"child_document_id"`
	Relationship     string    `json:"relationship"`
	CreatedAt        time.Time `json:"created_at"`
}
