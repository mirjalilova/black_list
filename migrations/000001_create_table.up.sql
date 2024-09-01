CREATE TYPE action AS ENUM ('added', 'removed');

CREATE TABLE IF NOT EXISTS hr (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    approved_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS employees (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    position VARCHAR(50),
    hr_id UUID NOT NULL REFERENCES hr(id),
    is_blocked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS black_list (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    employee_id UUID UNIQUE REFERENCES employees(id),
    reason TEXT NOT NULL,
    blacklisted_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    added_by UUID REFERENCES hr(id),
    action action NOT NULL,
    employee_id UUID REFERENCES employees(id),
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);