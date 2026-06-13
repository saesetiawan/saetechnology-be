CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(140) UNIQUE NOT NULL,
    name VARCHAR(180) NOT NULL,
    tagline TEXT,
    summary TEXT,
    description TEXT,
    category VARCHAR(120),
    status VARCHAR(40) NOT NULL DEFAULT 'draft',
    price_label VARCHAR(120),
    image_url TEXT,
    demo_url TEXT,
    cta_label VARCHAR(120),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_products_status
    ON products(status);

CREATE INDEX IF NOT EXISTS idx_products_sort_order
    ON products(sort_order);

INSERT INTO products (
    slug,
    name,
    tagline,
    summary,
    description,
    category,
    status,
    price_label,
    image_url,
    demo_url,
    cta_label,
    sort_order,
    is_featured,
    metadata,
    created_at,
    updated_at
) VALUES (
    'saecommerce',
    'SAECommerce',
    'A ready-to-customize commerce system for modern sales operations.',
    'SAECommerce helps teams manage products, orders, customers, payments, and operational workflows in one scalable platform.',
    'Built for businesses that need a flexible ecommerce foundation without being locked into generic workflows. SAECommerce can be extended for catalogs, checkout flows, internal admin panels, payment integrations, fulfillment tracking, reporting, and customer operations.',
    'SaaS Commerce Platform',
    'published',
    'Custom implementation',
    '/images/saas-dashboard-preview.png',
    '#contact',
    'Discuss SAECommerce',
    10,
    TRUE,
    '{"features":["Product and catalog management","Order and payment workflow","Customer and admin dashboards","Custom checkout and integration flow","Reporting-ready backend APIs"],"use_cases":["B2B commerce portal","Internal ordering system","Multi-branch product sales","Custom marketplace foundation"],"tech_stack":["Next.js","NestJS","Golang","PostgreSQL","Redis","S3"],"faqs":[{"question":"Can SAECommerce be customized?","answer":"Yes. It is intended as a flexible foundation that can be adapted to business rules, checkout flows, integrations, and operations."},{"question":"Is this suitable for SaaS?","answer":"Yes. The architecture can support subscription logic, role-based dashboards, tenant-aware data, and backend integrations."}]}'::jsonb,
    NOW(),
    NOW()
)
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    tagline = EXCLUDED.tagline,
    summary = EXCLUDED.summary,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    status = EXCLUDED.status,
    price_label = EXCLUDED.price_label,
    image_url = EXCLUDED.image_url,
    demo_url = EXCLUDED.demo_url,
    cta_label = EXCLUDED.cta_label,
    sort_order = EXCLUDED.sort_order,
    is_featured = EXCLUDED.is_featured,
    metadata = EXCLUDED.metadata,
    deleted_at = NULL,
    updated_at = NOW();
