-- +goose Up
BEGIN;

INSERT INTO role_templates (name, display_name, description, permissions) VALUES
('admin', 'Admin', 'Full platform access including user management and billing', '[
  {"resource":"gst_returns","action":"view"},{"resource":"gst_returns","action":"edit"},{"resource":"gst_returns","action":"file"},
  {"resource":"gstr_9_9c","action":"view"},{"resource":"gstr_9_9c","action":"edit"},{"resource":"gstr_9_9c","action":"file"},
  {"resource":"e_invoicing","action":"view"},{"resource":"e_invoicing","action":"generate"},{"resource":"e_invoicing","action":"cancel"},
  {"resource":"e_way_bill","action":"view"},{"resource":"e_way_bill","action":"generate"},{"resource":"e_way_bill","action":"cancel"},
  {"resource":"itc_reconciliation","action":"view"},{"resource":"itc_reconciliation","action":"edit"},
  {"resource":"vendor_compliance","action":"view"},{"resource":"vendor_compliance","action":"edit"},
  {"resource":"tds","action":"view"},{"resource":"tds","action":"calculate"},{"resource":"tds","action":"file"},{"resource":"tds","action":"issue_cert"},
  {"resource":"itr","action":"view"},{"resource":"itr","action":"calculate"},{"resource":"itr","action":"file"},{"resource":"itr","action":"approve"},
  {"resource":"compliance_calendar","action":"view"},
  {"resource":"users_roles","action":"view"},{"resource":"users_roles","action":"manage"},
  {"resource":"connected_apps","action":"view"},{"resource":"connected_apps","action":"manage"},
  {"resource":"billing","action":"view"},{"resource":"billing","action":"manage"}
]'),

('tax_manager', 'Tax Manager', 'Manages all tax filings and compliance workflows', '[
  {"resource":"gst_returns","action":"view"},{"resource":"gst_returns","action":"edit"},{"resource":"gst_returns","action":"file"},
  {"resource":"gstr_9_9c","action":"view"},{"resource":"gstr_9_9c","action":"edit"},{"resource":"gstr_9_9c","action":"file"},
  {"resource":"e_invoicing","action":"view"},{"resource":"e_invoicing","action":"generate"},{"resource":"e_invoicing","action":"cancel"},
  {"resource":"e_way_bill","action":"view"},{"resource":"e_way_bill","action":"generate"},{"resource":"e_way_bill","action":"cancel"},
  {"resource":"itc_reconciliation","action":"view"},{"resource":"itc_reconciliation","action":"edit"},
  {"resource":"vendor_compliance","action":"view"},{"resource":"vendor_compliance","action":"edit"},
  {"resource":"tds","action":"view"},{"resource":"tds","action":"calculate"},{"resource":"tds","action":"file"},{"resource":"tds","action":"issue_cert"},
  {"resource":"itr","action":"view"},{"resource":"itr","action":"calculate"},{"resource":"itr","action":"file"},{"resource":"itr","action":"approve"},
  {"resource":"compliance_calendar","action":"view"},
  {"resource":"connected_apps","action":"view"}
]'),

('ap_manager', 'AP Manager', 'Manages accounts payable compliance — GST input, vendor compliance, TDS', '[
  {"resource":"gst_returns","action":"view"},{"resource":"gst_returns","action":"edit"},{"resource":"gst_returns","action":"file"},
  {"resource":"e_invoicing","action":"view"},{"resource":"e_invoicing","action":"generate"},
  {"resource":"e_way_bill","action":"view"},{"resource":"e_way_bill","action":"generate"},
  {"resource":"itc_reconciliation","action":"view"},{"resource":"itc_reconciliation","action":"edit"},
  {"resource":"vendor_compliance","action":"view"},{"resource":"vendor_compliance","action":"edit"},
  {"resource":"tds","action":"view"},{"resource":"tds","action":"calculate"},{"resource":"tds","action":"file"},{"resource":"tds","action":"issue_cert"},
  {"resource":"compliance_calendar","action":"view"},
  {"resource":"connected_apps","action":"view"}
]'),

('ap_executive', 'AP Executive', 'Handles day-to-day accounts payable data entry and views', '[
  {"resource":"gst_returns","action":"view"},
  {"resource":"e_invoicing","action":"view"},{"resource":"e_invoicing","action":"generate"},
  {"resource":"e_way_bill","action":"view"},{"resource":"e_way_bill","action":"generate"},
  {"resource":"itc_reconciliation","action":"view"},
  {"resource":"vendor_compliance","action":"view"},
  {"resource":"tds","action":"view"},{"resource":"tds","action":"calculate"},
  {"resource":"compliance_calendar","action":"view"}
]'),

('ar_manager', 'AR Manager', 'Manages accounts receivable compliance — GST output, e-invoicing', '[
  {"resource":"gst_returns","action":"view"},{"resource":"gst_returns","action":"edit"},{"resource":"gst_returns","action":"file"},
  {"resource":"gstr_9_9c","action":"view"},{"resource":"gstr_9_9c","action":"edit"},{"resource":"gstr_9_9c","action":"file"},
  {"resource":"e_invoicing","action":"view"},{"resource":"e_invoicing","action":"generate"},{"resource":"e_invoicing","action":"cancel"},
  {"resource":"e_way_bill","action":"view"},{"resource":"e_way_bill","action":"generate"},{"resource":"e_way_bill","action":"cancel"},
  {"resource":"itc_reconciliation","action":"view"},{"resource":"itc_reconciliation","action":"edit"},
  {"resource":"vendor_compliance","action":"view"},
  {"resource":"tds","action":"view"},{"resource":"tds","action":"calculate"},
  {"resource":"compliance_calendar","action":"view"},
  {"resource":"connected_apps","action":"view"}
]'),

('ar_executive', 'AR Executive', 'Handles day-to-day accounts receivable data entry and views', '[
  {"resource":"gst_returns","action":"view"},
  {"resource":"e_invoicing","action":"view"},
  {"resource":"e_way_bill","action":"view"},
  {"resource":"vendor_compliance","action":"view"},
  {"resource":"tds","action":"view"},
  {"resource":"compliance_calendar","action":"view"}
]'),

('auditor', 'Auditor', 'Read-only access to all compliance modules for audit purposes', '[
  {"resource":"gst_returns","action":"view"},
  {"resource":"gstr_9_9c","action":"view"},
  {"resource":"e_invoicing","action":"view"},
  {"resource":"e_way_bill","action":"view"},
  {"resource":"itc_reconciliation","action":"view"},
  {"resource":"vendor_compliance","action":"view"},
  {"resource":"tds","action":"view"},
  {"resource":"itr","action":"view"},
  {"resource":"compliance_calendar","action":"view"}
]');

COMMIT;

-- +goose Down
DELETE FROM role_templates WHERE name IN ('admin', 'tax_manager', 'ap_manager', 'ap_executive', 'ar_manager', 'ar_executive', 'auditor');
