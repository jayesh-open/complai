package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/document-service/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Pool() *pgxpool.Pool { return s.pool }

func setTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

func (s *Store) CreateDocument(ctx context.Context, tenantID uuid.UUID, d *domain.Document) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	d.ID = uuid.New()
	d.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO documents (id, tenant_id, document_type, document_number, file_name, mime_type, file_size,
		                        s3_bucket, s3_key, encrypted_dek, kms_key_arn, encryption_algo, tags, metadata, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		 RETURNING created_at, updated_at`,
		d.ID, d.TenantID, d.DocumentType, d.DocumentNumber, d.FileName, d.MimeType, d.FileSize,
		d.S3Bucket, d.S3Key, d.EncryptedDEK, d.KMSKeyARN, d.EncryptionAlgo, d.Tags, d.Metadata, d.CreatedBy,
	).Scan(&d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert document: %w", err)
	}

	d.VirusStatus = "pending"
	d.OCRStatus = "none"
	return tx.Commit(ctx)
}

func (s *Store) GetDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Document, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var d domain.Document
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, document_type, document_number, file_name, mime_type, file_size,
		        s3_bucket, s3_key, encrypted_dek, kms_key_arn, encryption_algo, virus_status, ocr_status,
		        ocr_result, tags, metadata, created_by, created_at, updated_at, deleted_at
		 FROM documents WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&d.ID, &d.TenantID, &d.DocumentType, &d.DocumentNumber, &d.FileName, &d.MimeType, &d.FileSize,
		&d.S3Bucket, &d.S3Key, &d.EncryptedDEK, &d.KMSKeyARN, &d.EncryptionAlgo, &d.VirusStatus, &d.OCRStatus,
		&d.OCRResult, &d.Tags, &d.Metadata, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}
	return &d, tx.Commit(ctx)
}

func (s *Store) ListDocuments(ctx context.Context, tenantID uuid.UUID) ([]domain.Document, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, document_type, document_number, file_name, mime_type, file_size,
		        s3_bucket, s3_key, encrypted_dek, kms_key_arn, encryption_algo, virus_status, ocr_status,
		        ocr_result, tags, metadata, created_by, created_at, updated_at, deleted_at
		 FROM documents WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()

	var docs []domain.Document
	for rows.Next() {
		var d domain.Document
		if err := rows.Scan(&d.ID, &d.TenantID, &d.DocumentType, &d.DocumentNumber, &d.FileName, &d.MimeType, &d.FileSize,
			&d.S3Bucket, &d.S3Key, &d.EncryptedDEK, &d.KMSKeyARN, &d.EncryptionAlgo, &d.VirusStatus, &d.OCRStatus,
			&d.OCRResult, &d.Tags, &d.Metadata, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt); err != nil {
			return nil, fmt.Errorf("scan document: %w", err)
		}
		docs = append(docs, d)
	}
	return docs, tx.Commit(ctx)
}

func (s *Store) UpdateOCRStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, result string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	var ocrResult *string
	if result != "" {
		ocrResult = &result
	}

	_, err = tx.Exec(ctx,
		`UPDATE documents SET ocr_status = $1, ocr_result = $2, updated_at = now() WHERE id = $3 AND deleted_at IS NULL`,
		status, ocrResult, id)
	if err != nil {
		return fmt.Errorf("update ocr status: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) DeleteDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE documents SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL`,
		id)
	if err != nil {
		return fmt.Errorf("soft delete document: %w", err)
	}
	return tx.Commit(ctx)
}
