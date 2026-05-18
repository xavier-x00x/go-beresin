-- Migration up: Add audit_logs table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(255) NOT NULL, -- e.g., 'login', 'register', 'reset_password'
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for performance querying on user audit trail
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
