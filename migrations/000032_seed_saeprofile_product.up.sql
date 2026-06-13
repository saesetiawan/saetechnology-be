INSERT INTO products (
    slug,
    name,
    tagline,
    summary,
    description,
    category,
    status,
    price_label,
    price_url,
    image_url,
    demo_url,
    cta_label,
    cta_url,
    sort_order,
    is_featured,
    metadata,
    created_at,
    updated_at
) VALUES (
    'saeprofile',
    'SAEProfile',
    'A highly customizable company profile website with content management built for professional business presence.',
    'SAEProfile helps businesses manage company pages, services, portfolios, testimonials, product highlights, brand assets, and contact funnels from an easy admin panel without repeatedly editing code.',
    'SAEProfile is designed for companies, consultants, agencies, and growing businesses that need more than a static company profile website. It provides a flexible content management foundation for home pages, about pages, product pages, profile sections, carousels, banners, contact forms, SEO metadata, theme colors, and custom business storytelling. The system can be adjusted to match industry-specific messaging, campaign pages, landing pages, and future content requests.',
    'Company Profile CMS',
    'published',
    'Custom website implementation',
    '/#contact',
    '/images/saas-dashboard-preview.png',
    '/products/saeprofile',
    'Build My Company Profile',
    '/#contact',
    20,
    TRUE,
    '{
        "features": [
            "Admin-managed website content by placement",
            "Custom home, about, profile, product, and landing sections",
            "Carousel, banner, promo, announcement, contact, and profile content blocks",
            "Theme setting for colors, contrast, labels, fonts, radius, and shadows",
            "Image upload support for logo, hero visuals, galleries, and content cards",
            "SEO-ready metadata and professional copy structure",
            "Contact form with captcha protection",
            "Flexible product and service showcase pages"
        ],
        "use_cases": [
            "Company profile website",
            "Professional service website",
            "Consultant or agency portfolio",
            "Product and SaaS landing pages",
            "Business profile with managed content",
            "Campaign-specific landing pages"
        ],
        "tech_stack": [
            "Next.js",
            "Golang",
            "PostgreSQL",
            "Redis",
            "S3-compatible storage",
            "Admin CMS"
        ],
        "sections": [
            {
                "key": "content-management",
                "type": "highlight",
                "title": "Manage Website Content Without Rebuilding the Site",
                "subtitle": "Content-first profile website",
                "body": "SAEProfile lets your team update company messaging, hero sections, service cards, product highlights, carousels, banners, contact sections, and profile pages directly from admin content placement.",
                "link_label": "Request SAEProfile",
                "link_url": "/#contact",
                "sort_order": 10
            },
            {
                "key": "custom-profile",
                "type": "content",
                "title": "More Flexible Than a Standard Company Profile",
                "subtitle": "Built to customize",
                "body": "Instead of a fixed brochure website, SAEProfile can be adapted for your brand structure, service categories, industry-specific copy, portfolio sections, team profiles, product pages, SEO content, and future marketing campaigns.",
                "link_label": "Discuss customization",
                "link_url": "/#contact",
                "sort_order": 20
            },
            {
                "key": "growth-ready",
                "type": "content",
                "title": "Ready for New Requests as Your Business Grows",
                "subtitle": "Expandable foundation",
                "body": "Start with a professional company profile, then expand into product landing pages, SaaS previews, inquiry flows, portfolio galleries, multilingual content, integrations, and custom modules when needed.",
                "link_label": "Plan the roadmap",
                "link_url": "/#contact",
                "sort_order": 30
            },
            {
                "key": "contact-us",
                "type": "contact",
                "title": "Ready to Launch a Better Company Profile?",
                "subtitle": "Contact Us",
                "body": "Tell us about your company, services, brand direction, content needs, and future website requests. We will help define the right SAEProfile implementation for your business.",
                "link_label": "Contact Us",
                "link_url": "/#contact",
                "sort_order": 40
            }
        ],
        "faqs": [
            {
                "question": "Can SAEProfile be used for a normal company profile website?",
                "answer": "Yes. It is suitable for company profile websites, service businesses, agencies, consultants, portfolios, and professional brand websites."
            },
            {
                "question": "What makes it different from a static website?",
                "answer": "SAEProfile is content-managed. Sections, placements, images, colors, product pages, profile pages, and contact content can be updated from admin without changing frontend code."
            },
            {
                "question": "Can the design and sections be customized?",
                "answer": "Yes. The layout, copy, colors, sections, carousels, pages, content types, and future modules can be customized around your business needs."
            },
            {
                "question": "Can this grow into a larger system?",
                "answer": "Yes. SAEProfile can become the foundation for landing pages, product catalogs, inquiry management, customer portals, dashboards, or SaaS modules."
            }
        ]
    }'::jsonb,
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
    price_url = EXCLUDED.price_url,
    image_url = EXCLUDED.image_url,
    demo_url = EXCLUDED.demo_url,
    cta_label = EXCLUDED.cta_label,
    cta_url = EXCLUDED.cta_url,
    sort_order = EXCLUDED.sort_order,
    is_featured = EXCLUDED.is_featured,
    metadata = EXCLUDED.metadata,
    deleted_at = NULL,
    updated_at = NOW();

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
    'products.saeprofile.solution',
    'section',
    'products/saeprofile',
    'A company profile website your team can actually manage.',
    'Managed content system',
    'SAEProfile turns your website into a flexible content platform. Admin users can manage homepage sections, about content, profile pages, product pages, carousels, service cards, banners, and contact messaging without touching code.',
    '',
    '/#contact',
    'Request SAEProfile',
    10,
    TRUE,
    '{"items":[{"title":"Placement-based content","description":"Show different content on home, about-us, profile, products/saeprofile, and future landing pages using placement keys.","sort_order":10},{"title":"Flexible section types","description":"Use hero, section, carousel, banner, promo, announcement, profile, and contact blocks for professional website storytelling.","sort_order":20},{"title":"Theme-controlled branding","description":"Manage colors, contrast, labels, typography, border radius, shadows, logos, and images from website settings.","sort_order":30},{"title":"Contact-ready conversion","description":"Use captcha-protected contact forms and clear call-to-action sections to receive business inquiries.","sort_order":40}]}'::jsonb
),
(
    'products.saeprofile.customization',
    'banner',
    'products/saeprofile',
    'Not just a company profile. A customizable website foundation.',
    'Flexible implementation',
    'SAEProfile can be shaped for corporate websites, service companies, agencies, consultants, professional portfolios, SaaS previews, product landing pages, and custom business campaigns.',
    '/images/saas-dashboard-preview.png',
    '/#contact',
    'Plan My Website',
    20,
    TRUE,
    '{"tone":"profile"}'::jsonb
),
(
    'products.saeprofile.carousel',
    'carousel',
    'products/saeprofile',
    'What you can manage with SAEProfile.',
    'CMS capabilities',
    'Highlight the content modules and business sections that make your website easier to customize and maintain.',
    '',
    '/#contact',
    'Customize SAEProfile',
    30,
    TRUE,
    '{"items":[{"title":"Company pages","subtitle":"Profile content","background":"#ffffff","description":"Manage home, about, company story, team, values, service descriptions, and contact sections.","link_label":"Discuss pages","link_url":"/#contact","image_url":"/images/saas-dashboard-preview.png","sort_order":10},{"title":"Service and product showcases","subtitle":"Business offering","background":"#ffffff","description":"Create service cards, product detail pages, feature lists, use cases, pricing labels, CTAs, and FAQ sections.","link_label":"Plan showcase","link_url":"/#contact","image_url":"/images/saas-dashboard-preview.png","sort_order":20},{"title":"Brand and theme control","subtitle":"Website setting","background":"#ffffff","description":"Adjust logo, images, colors, contrast, labels, fonts, rounded corners, shadows, and SEO metadata from admin.","link_label":"Set branding","link_url":"/#contact","image_url":"/images/saas-dashboard-preview.png","sort_order":30},{"title":"Campaign landing pages","subtitle":"Growth-ready","background":"#ffffff","description":"Add new placements for promotional pages, product launches, industry pages, and custom request funnels.","link_label":"Launch campaign","link_url":"/#contact","image_url":"/images/saas-dashboard-preview.png","sort_order":40}]}'::jsonb
),
(
    'products.saeprofile.contact',
    'section',
    'products/saeprofile',
    'Contact us to build your company profile system.',
    'Contact Us',
    'Share your business profile, services, preferred website structure, content needs, and future customization plans. We will help design the right SAEProfile setup.',
    '',
    '',
    '',
    40,
    TRUE,
    '{"tone":"contact"}'::jsonb
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
    deleted_at = NULL,
    updated_at = NOW();
