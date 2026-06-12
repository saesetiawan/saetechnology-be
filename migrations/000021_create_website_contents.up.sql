CREATE TABLE IF NOT EXISTS website_contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(120) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'section',
    placement VARCHAR(80) NOT NULL DEFAULT 'home',
    title VARCHAR(255) NOT NULL,
    subtitle TEXT,
    body TEXT,
    image_url TEXT,
    link_url TEXT,
    link_label VARCHAR(120),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    publish_start_at TIMESTAMPTZ,
    publish_end_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_website_contents_placement
    ON website_contents(placement);

CREATE INDEX IF NOT EXISTS idx_website_contents_type
    ON website_contents(type);

CREATE INDEX IF NOT EXISTS idx_website_contents_is_active
    ON website_contents(is_active);

INSERT INTO website_contents (
    key,
    type,
    placement,
    title,
    subtitle,
    body,
    image_url,
    link_url,
    link_label,
    sort_order,
    is_active,
    metadata
) VALUES
(
    'home-hero-primary',
    'hero',
    'home',
    'Build software that fits the way your business actually works.',
    'Custom apps, websites, automation, and SaaS products from SAE Technology Solution.',
    'We help design, build, and maintain digital solutions from concept to production-ready operations.',
    '',
    '/contact',
    'Start a Project',
    10,
    TRUE,
    '{"tone":"primary"}'::jsonb
),
(
    'home-services-overview',
    'section',
    'home',
    'Custom software, custom apps, and system integrations',
    'For businesses that need cleaner workflows without forcing their process into generic tools.',
    'This section can be managed from the admin panel to highlight priority services, portfolio items, or consultation offers.',
    '',
    '/services',
    'View Services',
    20,
    TRUE,
    '{"tone":"service"}'::jsonb
),
(
    'home-saas-offer',
    'section',
    'home',
    'Ready-to-use SaaS products',
    'Choose an existing SaaS product to accelerate operations without starting from scratch.',
    'Use this section to showcase SaaS products, pricing packages, demos, or trial-focused calls to action.',
    '',
    '/saas',
    'View SaaS Products',
    30,
    TRUE,
    '{"tone":"saas"}'::jsonb
)
ON CONFLICT (key) DO UPDATE SET
    type = EXCLUDED.type,
    placement = EXCLUDED.placement,
    title = EXCLUDED.title,
    subtitle = EXCLUDED.subtitle,
    body = EXCLUDED.body,
    image_url = EXCLUDED.image_url,
    link_url = EXCLUDED.link_url,
    link_label = EXCLUDED.link_label,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    metadata = EXCLUDED.metadata,
    updated_at = NOW();
