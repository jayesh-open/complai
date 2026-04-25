output "dns_record_ids" {
  description = "Map of DNS record names to their Cloudflare record IDs"
  value = {
    root          = cloudflare_record.root.id
    wildcard      = cloudflare_record.wildcard.id
    api           = cloudflare_record.api.id
    vendor_portal = cloudflare_record.vendor_portal.id
  }
}
