-- ==========================================
-- USERS QUERIES
-- ==========================================

-- name: CreateUser :one
INSERT INTO users (
  email, phone, password_hash, full_name, role, avatar_url, city, is_verified, is_active, google_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
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
  user_id, bio, tagline, years_experience, service_radius_km, location, is_kyc_verified, kyc_document_url, kyc_selfie_url, subscription_tier, average_rating
) VALUES (
  $1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint(cast(sqlc.arg(longitude) as float8), cast(sqlc.arg(latitude) as float8)), 4326)::geography, $6, $7, $8, $9, $10
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
  name, slug, parent_id, icon_url, sort_order, is_active
) VALUES (
  $1, $2, $3, $4, $5, $6
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
  talent_id, category_id, title, description, tags, cover_image_url, gallery_urls, video_urls, status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
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
  listing_id, name, description, price_amount, price_type, duration_hours, inclusions, max_revisions, sort_order
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetServicePackagesByListing :many
SELECT * FROM service_packages WHERE listing_id = $1 ORDER BY sort_order ASC;


-- ==========================================
-- ORDERS QUERIES
-- ==========================================

-- name: CreateOrder :one
INSERT INTO orders (
  order_number, user_id, talent_id, listing_id, package_id, tender_id, bid_id, title, description, work_date_start, work_date_end, location_address, final_amount, dp_amount, remaining_amount, platform_fee, talent_receive_amount, status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
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
  user_id, action, ip_address, user_agent
) VALUES (
  $1, $2, $3, $4
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
