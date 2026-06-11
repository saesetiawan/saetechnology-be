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
    'Bangun software yang benar-benar mengikuti cara kerja bisnis Anda',
    'Custom apps, website, automation, dan produk SaaS dari SAE Technology Solution.',
    'Kami membantu merancang, membangun, dan merawat solusi digital mulai dari ide sampai siap dipakai operasional.',
    '',
    '/contact',
    'Konsultasi Proyek',
    10,
    TRUE,
    '{"tone":"primary"}'::jsonb
),
(
    'home-services-overview',
    'section',
    'home',
    'Custom software, custom apps, dan integrasi sistem',
    'Untuk bisnis yang butuh alur kerja lebih rapi tanpa memaksakan proses ke tool generik.',
    'Konten ini bisa Anda kelola dari panel admin untuk menampilkan layanan prioritas, portofolio, atau promo konsultasi.',
    '',
    '/services',
    'Lihat Layanan',
    20,
    TRUE,
    '{"tone":"service"}'::jsonb
),
(
    'home-saas-offer',
    'section',
    'home',
    'Berlangganan SaaS siap pakai',
    'Pilih produk SaaS yang sudah tersedia untuk mempercepat operasional tanpa mulai dari nol.',
    'Gunakan section ini untuk menampilkan produk SaaS, paket harga, demo, atau call-to-action trial.',
    '',
    '/saas',
    'Lihat Produk SaaS',
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
