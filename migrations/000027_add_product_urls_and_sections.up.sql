ALTER TABLE products
    ADD COLUMN IF NOT EXISTS price_url TEXT,
    ADD COLUMN IF NOT EXISTS cta_url TEXT;

UPDATE products
SET
    price_url = '/#contact',
    cta_url = '/#contact',
    metadata = '{
        "features": [
            "Product and catalog management",
            "Order and payment workflow",
            "Customer and admin dashboards",
            "Custom checkout and integration flow",
            "Reporting-ready backend APIs"
        ],
        "use_cases": [
            "B2B commerce portal",
            "Internal ordering system",
            "Multi-branch product sales",
            "Custom marketplace foundation"
        ],
        "tech_stack": [
            "Next.js",
            "NestJS",
            "Golang",
            "PostgreSQL",
            "Redis",
            "S3"
        ],
        "sections": [
            {
                "key": "workflow",
                "type": "content",
                "title": "Operational Workflow Built Around Your Business",
                "subtitle": "Process-first commerce",
                "body": "SAECommerce can be shaped around your real sales flow, from catalog setup and customer segmentation to order approvals, payment handling, fulfillment tracking, reporting, and internal administration.",
                "link_label": "Discuss the workflow",
                "link_url": "/#contact",
                "sort_order": 10
            },
            {
                "key": "customization",
                "type": "highlight",
                "title": "Everything Important Can Be Customized",
                "subtitle": "Flexible implementation",
                "body": "Use SAECommerce as a ready foundation for a new system, custom ecommerce operation, marketplace flow, internal ordering platform, or SaaS commerce product. The interface, roles, integrations, reports, and business rules can be adjusted to match the way your team works.",
                "link_label": "Request customization",
                "link_url": "/#contact",
                "sort_order": 20
            },
            {
                "key": "delivery",
                "type": "content",
                "title": "Designed for Professional Delivery",
                "subtitle": "From idea to production",
                "body": "We can support discovery, architecture, backend API development, frontend implementation, admin dashboards, deployment preparation, and post-launch improvement so the product is ready for real operational use.",
                "link_label": "Plan my product",
                "link_url": "/#contact",
                "sort_order": 30
            },
            {
                "key": "contact-us",
                "type": "contact",
                "title": "Ready to Build With SAECommerce?",
                "subtitle": "Contact Us",
                "body": "Tell us about your commerce workflow, integration needs, and launch target. We will help map the implementation path and recommend the right development approach for your business.",
                "link_label": "Contact Us",
                "link_url": "/#contact",
                "sort_order": 40
            }
        ],
        "faqs": [
            {
                "question": "Can SAECommerce be customized?",
                "answer": "Yes. It is intended as a flexible foundation that can be adapted to business rules, checkout flows, integrations, dashboards, and operations."
            },
            {
                "question": "Is this suitable for SaaS?",
                "answer": "Yes. The architecture can support subscription logic, role-based dashboards, tenant-aware data, and backend integrations."
            },
            {
                "question": "Can I request a completely new system instead of ecommerce?",
                "answer": "Yes. SAE Technology Solution can build new custom systems, SaaS products, and operational platforms based on your requirements."
            }
        ]
    }'::jsonb,
    updated_at = NOW()
WHERE slug = 'saecommerce';
