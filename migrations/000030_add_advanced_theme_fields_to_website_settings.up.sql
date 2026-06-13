ALTER TABLE website_settings
    ADD COLUMN IF NOT EXISTS primary_contrast_color VARCHAR(20) DEFAULT '#ffffff',
    ADD COLUMN IF NOT EXISTS accent_contrast_color VARCHAR(20) DEFAULT '#ffffff',
    ADD COLUMN IF NOT EXISTS surface_contrast_color VARCHAR(20) DEFAULT '#0f172a',
    ADD COLUMN IF NOT EXISTS success_color VARCHAR(20) DEFAULT '#10b981',
    ADD COLUMN IF NOT EXISTS warning_color VARCHAR(20) DEFAULT '#f59e0b',
    ADD COLUMN IF NOT EXISTS danger_color VARCHAR(20) DEFAULT '#ef4444',
    ADD COLUMN IF NOT EXISTS info_color VARCHAR(20) DEFAULT '#3b82f6',
    ADD COLUMN IF NOT EXISTS label_color VARCHAR(20) DEFAULT '#334155',
    ADD COLUMN IF NOT EXISTS label_background_color VARCHAR(20) DEFAULT '#f1f5f9',
    ADD COLUMN IF NOT EXISTS font_family VARCHAR(160) DEFAULT 'Plus Jakarta Sans, ui-sans-serif, system-ui, sans-serif',
    ADD COLUMN IF NOT EXISTS heading_font_family VARCHAR(160) DEFAULT 'Plus Jakarta Sans, ui-sans-serif, system-ui, sans-serif',
    ADD COLUMN IF NOT EXISTS border_radius VARCHAR(20) DEFAULT '2rem',
    ADD COLUMN IF NOT EXISTS button_radius VARCHAR(20) DEFAULT '999px',
    ADD COLUMN IF NOT EXISTS shadow_style VARCHAR(160) DEFAULT '0 18px 60px rgba(15, 23, 42, 0.10)';

UPDATE website_settings
SET
    primary_contrast_color = COALESCE(primary_contrast_color, '#ffffff'),
    accent_contrast_color = COALESCE(accent_contrast_color, '#ffffff'),
    surface_contrast_color = COALESCE(surface_contrast_color, text_color, '#0f172a'),
    success_color = COALESCE(success_color, '#10b981'),
    warning_color = COALESCE(warning_color, '#f59e0b'),
    danger_color = COALESCE(danger_color, '#ef4444'),
    info_color = COALESCE(info_color, '#3b82f6'),
    label_color = COALESCE(label_color, text_color, '#334155'),
    label_background_color = COALESCE(label_background_color, secondary_color, '#f1f5f9'),
    font_family = COALESCE(font_family, 'Plus Jakarta Sans, ui-sans-serif, system-ui, sans-serif'),
    heading_font_family = COALESCE(heading_font_family, font_family, 'Plus Jakarta Sans, ui-sans-serif, system-ui, sans-serif'),
    border_radius = COALESCE(border_radius, '2rem'),
    button_radius = COALESCE(button_radius, '999px'),
    shadow_style = COALESCE(shadow_style, '0 18px 60px rgba(15, 23, 42, 0.10)');
