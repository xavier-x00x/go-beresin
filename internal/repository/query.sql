-- ==========================================
-- USERS QUERIES
-- ==========================================

-- name: CreateUser :one
INSERT INTO users (
  id, email, phone, password_hash, full_name, role, avatar_url, city, is_verified, is_active, google_id
) VALUES (
  sqlc.arg(id), sqlc.arg(email), sqlc.arg(phone), sqlc.arg(password_hash), sqlc.arg(full_name), sqlc.arg(role), sqlc.arg(avatar_url), sqlc.arg(city), sqlc.arg(is_verified), sqlc.arg(is_active), sqlc.arg(google_id)
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;


-- ==========================================
-- TALENT PROFILES QUERIES
-- ==========================================

-- name: CreateTalentProfile :one
INSERT INTO talent_profiles (
  id, user_id, bio, tagline, years_experience, service_radius_km, location, is_kyc_verified, kyc_document_url, kyc_selfie_url, subscription_tier, average_rating
) VALUES (
  sqlc.arg(id), sqlc.arg(user_id), sqlc.arg(bio), sqlc.arg(tagline), sqlc.arg(years_experience), sqlc.arg(service_radius_km), ST_SetSRID(ST_MakePoint(cast(sqlc.arg(longitude) as float8), cast(sqlc.arg(latitude) as float8)), 4326)::geography, sqlc.arg(is_kyc_verified), sqlc.arg(kyc_document_url), sqlc.arg(kyc_selfie_url), sqlc.arg(subscription_tier), sqlc.arg(average_rating)
)
RETURNING *;

-- name: GetTalentProfileByID :one
SELECT * FROM talent_profiles WHERE id = $1 LIMIT 1;

-- name: GetTalentProfileByUserID :one
SELECT * FROM talent_profiles WHERE user_id = $1 LIMIT 1;

-- name: GetTalentProfilesInRadius :many
-- ST_DWithin works directly on geography points with distance in meters
SELECT 
  tp.*, 
  u.full_name, 
  u.avatar_url, 
  ST_AsText(tp.location)::text as location_text,
  ST_Distance(tp.location, ST_SetSRID(ST_MakePoint(cast(sqlc.arg(longitude) as float8), cast(sqlc.arg(latitude) as float8)), 4326)::geography) as distance_meters
FROM talent_profiles tp
JOIN users u ON tp.user_id = u.id
WHERE ST_DWithin(
  tp.location,
  ST_SetSRID(ST_MakePoint(cast(sqlc.arg(longitude) as float8), cast(sqlc.arg(latitude) as float8)), 4326)::geography,
  cast(sqlc.arg(radius_meters) as float8)
)
ORDER BY distance_meters ASC;


-- ==========================================
-- CATEGORIES QUERIES
-- ==========================================

-- name: CreateCategory :one
INSERT INTO categories (
  id, name, slug, parent_id, icon_url, sort_order, is_active
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(slug), sqlc.arg(parent_id), sqlc.arg(icon_url), sqlc.arg(sort_order), sqlc.arg(is_active)
)
RETURNING *;

-- name: ListCategories :many
SELECT * FROM categories WHERE parent_id IS NULL AND is_active = TRUE ORDER BY sort_order ASC;

-- name: ListSubCategories :many
SELECT * FROM categories WHERE parent_id = $1 AND is_active = TRUE ORDER BY sort_order ASC;


-- ==========================================
-- SERVICE LISTINGS QUERIES
-- ==========================================

-- name: CreateServiceListing :one
INSERT INTO service_listings (
  id, talent_id, category_id, title, description, tags, cover_image_url, gallery_urls, video_urls, status
) VALUES (
  sqlc.arg(id), sqlc.arg(talent_id), sqlc.arg(category_id), sqlc.arg(title), sqlc.arg(description), sqlc.arg(tags), sqlc.arg(cover_image_url), sqlc.arg(gallery_urls), sqlc.arg(video_urls), sqlc.arg(status)
)
RETURNING *;

-- name: GetServiceListingsByTalent :many
SELECT * FROM service_listings WHERE talent_id = $1 AND status != 'deleted' ORDER BY created_at DESC;

-- name: ListActiveServiceListings :many
SELECT 
  sl.*, 
  tp.average_rating, 
  tp.total_reviews, 
  u.full_name as talent_name,
  u.city as talent_city
FROM service_listings sl
JOIN talent_profiles tp ON sl.talent_id = tp.id
JOIN users u ON tp.user_id = u.id
WHERE sl.status = 'published' AND sl.category_id = $1
ORDER BY sl.created_at DESC;


-- ==========================================
-- SERVICE PACKAGES QUERIES
-- ==========================================

-- name: CreateServicePackage :one
INSERT INTO service_packages (
  id, listing_id, name, description, price_amount, price_type, duration_hours, inclusions, max_revisions, sort_order
) VALUES (
  sqlc.arg(id), sqlc.arg(listing_id), sqlc.arg(name), sqlc.arg(description), sqlc.arg(price_amount), sqlc.arg(price_type), sqlc.arg(duration_hours), sqlc.arg(inclusions), sqlc.arg(max_revisions), sqlc.arg(sort_order)
)
RETURNING *;

-- name: GetServicePackagesByListing :many
SELECT * FROM service_packages WHERE listing_id = $1 ORDER BY sort_order ASC;


-- ==========================================
-- ORDERS QUERIES
-- ==========================================

-- name: CreateOrder :one
INSERT INTO orders (
  id, order_number, user_id, talent_id, listing_id, package_id, tender_id, bid_id, title, description, work_date_start, work_date_end, location_address, final_amount, dp_amount, remaining_amount, platform_fee, talent_receive_amount, status
) VALUES (
  sqlc.arg(id), sqlc.arg(order_number), sqlc.arg(user_id), sqlc.arg(talent_id), sqlc.arg(listing_id), sqlc.arg(package_id), sqlc.arg(tender_id), sqlc.arg(bid_id), sqlc.arg(title), sqlc.arg(description), sqlc.arg(work_date_start), sqlc.arg(work_date_end), sqlc.arg(location_address), sqlc.arg(final_amount), sqlc.arg(dp_amount), sqlc.arg(remaining_amount), sqlc.arg(platform_fee), sqlc.arg(talent_receive_amount), sqlc.arg(status)
)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: GetOrdersByUserID :many
SELECT o.*, u.full_name as talent_name
FROM orders o
JOIN talent_profiles tp ON o.talent_id = tp.id
JOIN users u ON tp.user_id = u.id
WHERE o.user_id = $1
ORDER BY o.created_at DESC;

-- name: GetOrdersByTalentID :many
SELECT o.*, u.full_name as user_name
FROM orders o
JOIN users u ON o.user_id = u.id
WHERE o.talent_id = $1
ORDER BY o.created_at DESC;

-- ==========================================
-- AUTHENTICATION & AUDIT QUERIES
-- ==========================================

-- name: CreateAuditLog :one
INSERT INTO audit_logs (
  id, user_id, action, ip_address, user_agent
) VALUES (
  sqlc.arg(id), sqlc.arg(user_id), sqlc.arg(action), sqlc.arg(ip_address), sqlc.arg(user_agent)
)
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users 
SET password_hash = $2, updated_at = NOW() 
WHERE id = $1 
RETURNING *;

-- name: UpdateUserVerificationStatus :one
UPDATE users 
SET is_verified = $2, updated_at = NOW() 
WHERE id = $1 
RETURNING *;

-- name: GetUserByGoogleID :one
SELECT * FROM users WHERE google_id = $1 LIMIT 1;

-- name: UpdateUserGoogleID :one
UPDATE users 
SET google_id = $2, updated_at = NOW() 
WHERE id = $1 
RETURNING *;
