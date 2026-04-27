package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/ewb-service/internal/domain"
)

type PGStore struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PGStore {
	return &PGStore{pool: pool}
}

func setTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

func (s *PGStore) CreateEWB(ctx context.Context, tenantID uuid.UUID, ewb *domain.EWayBill) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}

	ewb.ID = uuid.New()
	ewb.TenantID = tenantID
	ewb.RequestID = uuid.New()
	ewb.Status = domain.EWBStatusPending
	now := time.Now()
	ewb.CreatedAt = now
	ewb.UpdatedAt = now

	_, err = tx.Exec(ctx, `INSERT INTO ewb (
		id, tenant_id, doc_type, doc_number, doc_date,
		supplier_gstin, supplier_name, buyer_gstin, buyer_name,
		supply_type, sub_supply_type, transport_mode,
		vehicle_number, vehicle_type, transporter_id,
		from_place, from_state, from_pincode,
		to_place, to_state, to_pincode,
		distance_km, taxable_value, cgst_amount, sgst_amount, igst_amount, cess_amount, total_value,
		status, request_id, source_system, source_id, created_at, updated_at
	) VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34
	)`,
		ewb.ID, ewb.TenantID, ewb.DocType, ewb.DocNumber, ewb.DocDate,
		ewb.SupplierGSTIN, ewb.SupplierName, ewb.BuyerGSTIN, ewb.BuyerName,
		ewb.SupplyType, ewb.SubSupplyType, ewb.TransportMode,
		ewb.VehicleNumber, ewb.VehicleType, ewb.TransporterID,
		ewb.FromPlace, ewb.FromState, ewb.FromPincode,
		ewb.ToPlace, ewb.ToState, ewb.ToPincode,
		ewb.DistanceKM, ewb.TaxableValue, ewb.CGSTAmount, ewb.SGSTAmount, ewb.IGSTAmount, ewb.CessAmount, ewb.TotalValue,
		ewb.Status, ewb.RequestID, ewb.SourceSystem, ewb.SourceID, ewb.CreatedAt, ewb.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) GetEWB(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EWayBill, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, err
	}

	var ewb domain.EWayBill
	err = tx.QueryRow(ctx, `SELECT
		id, tenant_id, ewb_number, doc_type, doc_number, doc_date,
		supplier_gstin, supplier_name, buyer_gstin, buyer_name,
		supply_type, sub_supply_type, transport_mode,
		vehicle_number, vehicle_type, transporter_id,
		from_place, from_state, from_pincode,
		to_place, to_state, to_pincode,
		distance_km, taxable_value, cgst_amount, sgst_amount, igst_amount, cess_amount, total_value,
		status, valid_from, valid_until, generated_at, cancelled_at, cancel_reason,
		request_id, source_system, source_id, created_at, updated_at
	FROM ewb WHERE id = $1`, id).Scan(
		&ewb.ID, &ewb.TenantID, &ewb.EWBNumber, &ewb.DocType, &ewb.DocNumber, &ewb.DocDate,
		&ewb.SupplierGSTIN, &ewb.SupplierName, &ewb.BuyerGSTIN, &ewb.BuyerName,
		&ewb.SupplyType, &ewb.SubSupplyType, &ewb.TransportMode,
		&ewb.VehicleNumber, &ewb.VehicleType, &ewb.TransporterID,
		&ewb.FromPlace, &ewb.FromState, &ewb.FromPincode,
		&ewb.ToPlace, &ewb.ToState, &ewb.ToPincode,
		&ewb.DistanceKM, &ewb.TaxableValue, &ewb.CGSTAmount, &ewb.SGSTAmount, &ewb.IGSTAmount, &ewb.CessAmount, &ewb.TotalValue,
		&ewb.Status, &ewb.ValidFrom, &ewb.ValidUntil, &ewb.GeneratedAt, &ewb.CancelledAt, &ewb.CancelReason,
		&ewb.RequestID, &ewb.SourceSystem, &ewb.SourceID, &ewb.CreatedAt, &ewb.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ewb, tx.Commit(ctx)
}

func (s *PGStore) GetEWBByNumber(ctx context.Context, tenantID uuid.UUID, ewbNumber string) (*domain.EWayBill, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, err
	}

	var ewb domain.EWayBill
	err = tx.QueryRow(ctx, `SELECT
		id, tenant_id, ewb_number, doc_type, doc_number, doc_date,
		supplier_gstin, supplier_name, buyer_gstin, buyer_name,
		supply_type, transport_mode, vehicle_number, vehicle_type,
		distance_km, taxable_value, total_value, status,
		valid_from, valid_until, generated_at, cancelled_at, cancel_reason,
		request_id, source_system, created_at, updated_at
	FROM ewb WHERE ewb_number = $1`, ewbNumber).Scan(
		&ewb.ID, &ewb.TenantID, &ewb.EWBNumber, &ewb.DocType, &ewb.DocNumber, &ewb.DocDate,
		&ewb.SupplierGSTIN, &ewb.SupplierName, &ewb.BuyerGSTIN, &ewb.BuyerName,
		&ewb.SupplyType, &ewb.TransportMode, &ewb.VehicleNumber, &ewb.VehicleType,
		&ewb.DistanceKM, &ewb.TaxableValue, &ewb.TotalValue, &ewb.Status,
		&ewb.ValidFrom, &ewb.ValidUntil, &ewb.GeneratedAt, &ewb.CancelledAt, &ewb.CancelReason,
		&ewb.RequestID, &ewb.SourceSystem, &ewb.CreatedAt, &ewb.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ewb, tx.Commit(ctx)
}

func (s *PGStore) ListEWBs(ctx context.Context, tenantID uuid.UUID, req *domain.ListEWBRequest) ([]domain.EWayBill, int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, 0, err
	}

	var total int
	countQ := `SELECT COUNT(*) FROM ewb WHERE supplier_gstin = $1`
	args := []interface{}{req.GSTIN}
	if req.Status != "" {
		countQ += ` AND status = $2`
		args = append(args, req.Status)
	}
	if err := tx.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := `SELECT id, tenant_id, ewb_number, doc_type, doc_number, doc_date,
		supplier_gstin, buyer_gstin, vehicle_number, distance_km, total_value,
		status, valid_until, generated_at, created_at
	FROM ewb WHERE supplier_gstin = $1`
	if req.Status != "" {
		q += ` AND status = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = append(args, req.PageSize, req.PageOffset)
	} else {
		q += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, req.PageSize, req.PageOffset)
	}

	rows, err := tx.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var ewbs []domain.EWayBill
	for rows.Next() {
		var e domain.EWayBill
		if err := rows.Scan(
			&e.ID, &e.TenantID, &e.EWBNumber, &e.DocType, &e.DocNumber, &e.DocDate,
			&e.SupplierGSTIN, &e.BuyerGSTIN, &e.VehicleNumber, &e.DistanceKM, &e.TotalValue,
			&e.Status, &e.ValidUntil, &e.GeneratedAt, &e.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		ewbs = append(ewbs, e)
	}
	return ewbs, total, tx.Commit(ctx)
}

func (s *PGStore) UpdateEWBGenerated(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, ewbNumber string, validFrom, validUntil time.Time) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	now := time.Now()
	_, err = tx.Exec(ctx, `UPDATE ewb SET ewb_number=$1, status=$2, valid_from=$3, valid_until=$4, generated_at=$5, updated_at=$6 WHERE id=$7`,
		ewbNumber, domain.EWBStatusActive, validFrom, validUntil, now, now, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) UpdateEWBCancelled(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	now := time.Now()
	_, err = tx.Exec(ctx, `UPDATE ewb SET status=$1, cancelled_at=$2, cancel_reason=$3, updated_at=$4 WHERE id=$5`,
		domain.EWBStatusCancelled, now, reason, now, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) UpdateEWBStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status domain.EWBStatus) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE ewb SET status=$1, updated_at=$2 WHERE id=$3`, status, time.Now(), id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) UpdateEWBVehicle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, vehicleNumber string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE ewb SET vehicle_number=$1, status=$2, updated_at=$3 WHERE id=$4`,
		vehicleNumber, domain.EWBStatusVehicleUpdated, time.Now(), id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) UpdateEWBValidity(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, validUntil time.Time) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE ewb SET valid_until=$1, status=$2, updated_at=$3 WHERE id=$4`,
		validUntil, domain.EWBStatusExtended, time.Now(), id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) SetConsolidatedEWBID(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, consolidatedID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE ewb SET consolidated_ewb_id=$1, status=$2, updated_at=$3 WHERE id=$4`,
		consolidatedID, domain.EWBStatusConsolidated, time.Now(), id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) CreateItems(ctx context.Context, tenantID uuid.UUID, items []domain.EWBItem) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	for i, item := range items {
		item.ID = uuid.New()
		item.TenantID = tenantID
		item.ItemNumber = i + 1
		_, err = tx.Exec(ctx, `INSERT INTO ewb_items (
			id, ewb_id, tenant_id, item_number, product_name, product_desc,
			hsn_code, quantity, unit, taxable_value, cgst_rate, sgst_rate, igst_rate, cess_rate
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
			item.ID, item.EWBID, item.TenantID, item.ItemNumber,
			item.ProductName, item.ProductDesc, item.HSNCode,
			item.Quantity, item.Unit, item.TaxableValue,
			item.CGSTRate, item.SGSTRate, item.IGSTRate, item.CessRate,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *PGStore) GetItems(ctx context.Context, tenantID uuid.UUID, ewbID uuid.UUID) ([]domain.EWBItem, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, `SELECT id, ewb_id, tenant_id, item_number, product_name, hsn_code, quantity, unit, taxable_value FROM ewb_items WHERE ewb_id=$1 ORDER BY item_number`, ewbID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.EWBItem
	for rows.Next() {
		var item domain.EWBItem
		if err := rows.Scan(&item.ID, &item.EWBID, &item.TenantID, &item.ItemNumber, &item.ProductName, &item.HSNCode, &item.Quantity, &item.Unit, &item.TaxableValue); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, tx.Commit(ctx)
}

func (s *PGStore) CreateVehicleUpdate(ctx context.Context, tenantID uuid.UUID, update *domain.VehicleUpdate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	update.ID = uuid.New()
	update.TenantID = tenantID
	_, err = tx.Exec(ctx, `INSERT INTO ewb_vehicle_updates (id, ewb_id, tenant_id, vehicle_number, from_place, from_state, transport_mode, reason, remark)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		update.ID, update.EWBID, update.TenantID, update.VehicleNumber,
		update.FromPlace, update.FromState, update.TransportMode, update.Reason, update.Remark,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) GetVehicleUpdates(ctx context.Context, tenantID uuid.UUID, ewbID uuid.UUID) ([]domain.VehicleUpdate, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, `SELECT id, ewb_id, tenant_id, vehicle_number, from_place, from_state, transport_mode, reason, remark, updated_at
		FROM ewb_vehicle_updates WHERE ewb_id=$1 ORDER BY updated_at`, ewbID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var updates []domain.VehicleUpdate
	for rows.Next() {
		var u domain.VehicleUpdate
		if err := rows.Scan(&u.ID, &u.EWBID, &u.TenantID, &u.VehicleNumber, &u.FromPlace, &u.FromState, &u.TransportMode, &u.Reason, &u.Remark, &u.UpdatedAt); err != nil {
			return nil, err
		}
		updates = append(updates, u)
	}
	return updates, tx.Commit(ctx)
}

func (s *PGStore) CreateConsolidation(ctx context.Context, tenantID uuid.UUID, consolidation *domain.Consolidation) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return err
	}
	consolidation.ID = uuid.New()
	consolidation.TenantID = tenantID
	now := time.Now()
	consolidation.CreatedAt = now
	consolidation.UpdatedAt = now
	consolidation.GeneratedAt = &now
	_, err = tx.Exec(ctx, `INSERT INTO ewb_consolidations (
		id, tenant_id, consolidated_ewb_number, vehicle_number,
		from_place, from_state, to_place, to_state, transport_mode, status, generated_at, created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		consolidation.ID, consolidation.TenantID, consolidation.ConsolidatedEWBNumber,
		consolidation.VehicleNumber, consolidation.FromPlace, consolidation.FromState,
		consolidation.ToPlace, consolidation.ToState, consolidation.TransportMode,
		consolidation.Status, consolidation.GeneratedAt, consolidation.CreatedAt, consolidation.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

var _ Repository = (*PGStore)(nil)
