-- Enable PostGIS extension for geo-search
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. USERS
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE,
  phone VARCHAR(20) UNIQUE NOT NULL,
  password_hash VARCHAR(255),
  full_name VARCHAR(100) NOT NULL,
  role VARCHAR(50) NOT NULL CHECK (role IN ('user', 'talent', 'admin')),
  avatar_url TEXT,
  city VARCHAR(100),
  is_verified BOOLEAN DEFAULT FALSE,
  is_active BOOLEAN DEFAULT TRUE,
  google_id VARCHAR(255),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. TALENT PROFILES
CREATE TABLE talent_profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) UNIQUE NOT NULL,
  bio TEXT,
  tagline VARCHAR(100),
  years_experience INT,
  service_radius_km INT DEFAULT 20,
  -- PostGIS type instead of separate DECIMAL latitude/longitude
  location geography(POINT, 4326),
  is_kyc_verified BOOLEAN DEFAULT FALSE,
  kyc_document_url TEXT,
  kyc_selfie_url TEXT,
  kyc_reviewed_at TIMESTAMPTZ,
  kyc_reviewed_by UUID REFERENCES users(id),
  subscription_tier VARCHAR(50) DEFAULT 'free' CHECK (subscription_tier IN ('free', 'pro', 'business')),
  subscription_expires_at TIMESTAMPTZ,
  average_rating DECIMAL(3,2) DEFAULT 0,
  total_reviews INT DEFAULT 0,
  total_completed_jobs INT DEFAULT 0,
  response_time_hours INT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. SERVICE CATEGORIES
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) UNIQUE NOT NULL,
  parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
  icon_url TEXT,
  sort_order INT DEFAULT 0,
  is_active BOOLEAN DEFAULT TRUE
);

-- 4. SERVICE LISTINGS
CREATE TABLE service_listings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  talent_id UUID REFERENCES talent_profiles(id) ON DELETE CASCADE NOT NULL,
  category_id UUID REFERENCES categories(id) ON DELETE RESTRICT NOT NULL,
  title VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  tags TEXT[], -- array of tags
  cover_image_url TEXT,
  gallery_urls TEXT[],
  video_urls TEXT[],
  status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'paused', 'deleted')),
  view_count INT DEFAULT 0,
  booking_count INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. SERVICE PACKAGES (per listing)
CREATE TABLE service_packages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  listing_id UUID REFERENCES service_listings(id) ON DELETE CASCADE NOT NULL,
  name VARCHAR(100) NOT NULL, -- "Basic", "Standard", "Premium"
  description TEXT,
  price_amount BIGINT NOT NULL, -- dalam rupiah
  price_type VARCHAR(50) DEFAULT 'fixed' CHECK (price_type IN ('fixed', 'starting_from')),
  duration_hours INT, -- durasi layanan dalam jam
  inclusions TEXT[], -- apa yang termasuk
  max_revisions INT DEFAULT 0,
  sort_order INT DEFAULT 0
);

-- 6. TENDERS
CREATE TABLE tenders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  category_id UUID REFERENCES categories(id) ON DELETE RESTRICT NOT NULL,
  title VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  reference_media_urls TEXT[],
  location_address TEXT,
  location_lat DECIMAL(10,8),
  location_lng DECIMAL(11,8),
  work_date_start DATE,
  work_date_end DATE,
  budget_min BIGINT,
  budget_max BIGINT,
  budget_is_negotiable BOOLEAN DEFAULT TRUE,
  is_urgent BOOLEAN DEFAULT FALSE,
  min_experience_years INT DEFAULT 0,
  min_rating DECIMAL(3,2) DEFAULT 0,
  require_verified BOOLEAN DEFAULT FALSE,
  prefer_team BOOLEAN DEFAULT FALSE,
  max_bids INT,
  expires_at TIMESTAMPTZ,
  status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'closed', 'completed', 'cancelled')),
  selected_bid_id UUID, -- FK ke bids (setelah user pilih)
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. BIDS (penawaran pada tender)
CREATE TABLE bids (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tender_id UUID REFERENCES tenders(id) ON DELETE CASCADE NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) ON DELETE CASCADE NOT NULL,
  cover_letter TEXT NOT NULL,
  bid_amount BIGINT NOT NULL,
  estimated_duration_hours INT,
  dp_percentage INT DEFAULT 50,
  max_revisions INT DEFAULT 0,
  portfolio_ids UUID[], -- referensi ke portfolio items
  status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'shortlisted', 'accepted', 'rejected', 'withdrawn')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(tender_id, talent_id)
);

-- Tambah FK dari tenders ke bids setelah bids terbuat
ALTER TABLE tenders ADD CONSTRAINT fk_tenders_selected_bid_id FOREIGN KEY (selected_bid_id) REFERENCES bids(id) ON DELETE SET NULL;

-- 8. BOOKINGS / ORDERS
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_number VARCHAR(20) UNIQUE NOT NULL, -- BRS-20240801-0001
  user_id UUID REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) ON DELETE RESTRICT NOT NULL,
  listing_id UUID REFERENCES service_listings(id) ON DELETE SET NULL,
  package_id UUID REFERENCES service_packages(id) ON DELETE SET NULL,
  tender_id UUID REFERENCES tenders(id) ON DELETE SET NULL,
  bid_id UUID REFERENCES bids(id) ON DELETE SET NULL,
  title VARCHAR(200) NOT NULL,
  description TEXT,
  work_date_start DATE,
  work_date_end DATE,
  location_address TEXT,
  final_amount BIGINT NOT NULL,
  dp_amount BIGINT NOT NULL,
  remaining_amount BIGINT NOT NULL,
  platform_fee BIGINT NOT NULL,
  talent_receive_amount BIGINT NOT NULL,
  status VARCHAR(50) DEFAULT 'pending_confirmation' CHECK (status IN (
    'pending_confirmation',
    'confirmed',
    'dp_pending',
    'active',
    'completed_pending',
    'completed',
    'cancelled',
    'dispute'
  )),
  progress_percentage INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 9. CONTRACTS
CREATE TABLE contracts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  contract_number VARCHAR(20) UNIQUE NOT NULL,
  order_id UUID REFERENCES orders(id) ON DELETE CASCADE UNIQUE NOT NULL,
  user_id UUID REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) ON DELETE RESTRICT NOT NULL,
  work_title VARCHAR(200) NOT NULL,
  work_description TEXT NOT NULL,
  work_start_date DATE NOT NULL,
  work_end_date DATE NOT NULL,
  location TEXT,
  total_amount BIGINT NOT NULL,
  dp_percentage INT NOT NULL,
  dp_amount BIGINT NOT NULL,
  max_revisions INT DEFAULT 0,
  cancellation_terms TEXT,
  user_signed_at TIMESTAMPTZ,
  talent_signed_at TIMESTAMPTZ,
  user_signature_hash VARCHAR(255),
  talent_signature_hash VARCHAR(255),
  document_hash VARCHAR(255),
  pdf_url TEXT,
  status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'awaiting_user', 'awaiting_talent', 'active', 'completed', 'cancelled', 'dispute')),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 10. PAYMENTS
CREATE TABLE payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) ON DELETE RESTRICT NOT NULL,
  payment_type VARCHAR(50) NOT NULL CHECK (payment_type IN ('dp', 'remaining', 'full')),
  amount BIGINT NOT NULL,
  method VARCHAR(50) NOT NULL CHECK (method IN ('bank_transfer', 'ewallet', 'qris', 'credit_card', 'va')),
  provider VARCHAR(50), -- "midtrans", "gopay", "ovo", dll
  external_transaction_id VARCHAR(255),
  status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'expired', 'refunded')),
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 11. DISBURSEMENTS (pembayaran ke talent)
CREATE TABLE disbursements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) ON DELETE RESTRICT NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) ON DELETE RESTRICT NOT NULL,
  amount BIGINT NOT NULL,
  bank_code VARCHAR(10),
  bank_account VARCHAR(30),
  account_name VARCHAR(100),
  external_reference VARCHAR(255),
  status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'success', 'failed')),
  disbursed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 12. CHATS
CREATE TABLE chat_rooms (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type VARCHAR(50) NOT NULL CHECK (type IN ('order', 'bid_negotiation', 'inquiry')),
  order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
  tender_id UUID REFERENCES tenders(id) ON DELETE SET NULL,
  bid_id UUID REFERENCES bids(id) ON DELETE SET NULL,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) ON DELETE CASCADE NOT NULL,
  last_message_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 13. CHAT MESSAGES
CREATE TABLE chat_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id UUID REFERENCES chat_rooms(id) ON DELETE CASCADE NOT NULL,
  sender_id UUID REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
  message_type VARCHAR(50) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'video', 'audio', 'file', 'quotation', 'offer', 'system')),
  content TEXT,
  media_url TEXT,
  metadata JSONB, -- untuk quotation/offer: {amount, terms, etc}
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 14. REVIEWS
CREATE TABLE reviews (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) ON DELETE RESTRICT NOT NULL,
  reviewer_id UUID REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
  reviewee_id UUID REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
  review_type VARCHAR(50) NOT NULL CHECK (review_type IN ('user_to_talent', 'talent_to_user')),
  rating_overall DECIMAL(3,2) NOT NULL,
  rating_quality DECIMAL(3,2),
  rating_timeliness DECIMAL(3,2),
  rating_communication DECIMAL(3,2),
  rating_friendliness DECIMAL(3,2),
  rating_value DECIMAL(3,2),
  comment TEXT,
  photo_urls TEXT[],
  would_recommend BOOLEAN,
  talent_response TEXT,
  talent_responded_at TIMESTAMPTZ,
  is_flagged BOOLEAN DEFAULT FALSE,
  is_published BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 15. DISPUTES
CREATE TABLE disputes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) ON DELETE RESTRICT NOT NULL,
  filed_by UUID REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
  reason_category VARCHAR(100) NOT NULL,
  description TEXT NOT NULL,
  evidence_urls TEXT[],
  respondent_response TEXT,
  respondent_evidence_urls TEXT[],
  admin_notes TEXT,
  decision VARCHAR(50) DEFAULT 'pending' CHECK (decision IN ('user_wins', 'talent_wins', 'split', 'pending')),
  split_user_percentage INT,
  resolved_at TIMESTAMPTZ,
  resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
  status VARCHAR(50) DEFAULT 'open' CHECK (status IN ('open', 'under_review', 'resolved', 'appealed')),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 16. NOTIFICATIONS
CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  type VARCHAR(100) NOT NULL, -- "new_order", "payment_received", dll
  title VARCHAR(200) NOT NULL,
  body TEXT NOT NULL,
  action_url TEXT,
  metadata JSONB,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);


-- =========================================================================
-- DATABASE INDEXES
-- =========================================================================

-- Index Service Listings (explore/filter)
CREATE INDEX idx_service_listings_status ON service_listings(status);
CREATE INDEX idx_service_listings_category ON service_listings(category_id);

-- Index Orders (dashboard query)
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_talent_id ON orders(talent_id);
CREATE INDEX idx_orders_status ON orders(status);

-- Index Chat Messages (chat pagination)
CREATE INDEX idx_chat_messages_room_id ON chat_messages(room_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);

-- PostGIS Geo-Index on talent_profiles location
CREATE INDEX idx_talent_profiles_location ON talent_profiles USING GIST (location);
