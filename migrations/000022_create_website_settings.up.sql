CREATE TABLE IF NOT EXISTS website_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_name VARCHAR(120) NOT NULL,
    tagline TEXT,
    logo_url TEXT,
    favicon_url TEXT,
    email VARCHAR(160),
    phone VARCHAR(50),
    address TEXT,
    facebook_url TEXT,
    instagram_url TEXT,
    tiktok_url TEXT,
    primary_color VARCHAR(20) DEFAULT '#ec4899',
    accent_color VARCHAR(20) DEFAULT '#06b6d4',
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_website_settings_singleton
    ON website_settings ((true));

INSERT INTO website_settings (
    id,
    site_name,
    tagline,
    primary_color,
    accent_color,
    metadata
) VALUES (
    '00000000-0000-0000-0000-000000000100',
    'SAE Technology Solution',
    'Custom software, custom apps, dan SaaS untuk bisnis yang ingin bergerak lebih cepat',
    '#0f766e',
    '#f59e0b',
    '{}'::jsonb
) ON CONFLICT DO NOTHING;
