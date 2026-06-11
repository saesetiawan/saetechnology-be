ALTER TABLE website_settings
    ADD COLUMN IF NOT EXISTS primary_image_url TEXT,
    ADD COLUMN IF NOT EXISTS secondary_image_url TEXT,
    ADD COLUMN IF NOT EXISTS background_image_url TEXT,
    ADD COLUMN IF NOT EXISTS secondary_color VARCHAR(20) DEFAULT '#0ea5e9',
    ADD COLUMN IF NOT EXISTS background_color VARCHAR(20) DEFAULT '#f8fafc',
    ADD COLUMN IF NOT EXISTS surface_color VARCHAR(20) DEFAULT '#ffffff',
    ADD COLUMN IF NOT EXISTS text_color VARCHAR(20) DEFAULT '#0f172a',
    ADD COLUMN IF NOT EXISTS muted_text_color VARCHAR(20) DEFAULT '#64748b',
    ADD COLUMN IF NOT EXISTS border_color VARCHAR(20) DEFAULT '#e2e8f0';

UPDATE website_settings
SET
    secondary_color = COALESCE(secondary_color, accent_color, '#0ea5e9'),
    background_color = COALESCE(background_color, '#f8fafc'),
    surface_color = COALESCE(surface_color, '#ffffff'),
    text_color = COALESCE(text_color, '#0f172a'),
    muted_text_color = COALESCE(muted_text_color, '#64748b'),
    border_color = COALESCE(border_color, '#e2e8f0');
