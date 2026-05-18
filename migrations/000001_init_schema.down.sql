-- Disable selected bid ID foreign key first to avoid cyclic reference issues
ALTER TABLE IF EXISTS tenders DROP CONSTRAINT IF EXISTS fk_tenders_selected_bid_id;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS disputes CASCADE;
DROP TABLE IF EXISTS reviews CASCADE;
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS chat_rooms CASCADE;
DROP TABLE IF EXISTS disbursements CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS contracts CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS bids CASCADE;
DROP TABLE IF EXISTS tenders CASCADE;
DROP TABLE IF EXISTS service_packages CASCADE;
DROP TABLE IF EXISTS service_listings CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS talent_profiles CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop PostGIS extension
DROP EXTENSION IF EXISTS postgis CASCADE;
