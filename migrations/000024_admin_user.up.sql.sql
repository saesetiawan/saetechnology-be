INSERT INTO users (
    full_name,
    email,
    phone,
    password,
    status,
    role,
    email_verified_at,
    created_at,
    updated_at
) VALUES (
    'Braditya',
    'braditya12@gmail.com',
    '080000000001',
    '$argon2id$v=19$m=65536,t=3,p=2$zq26JWV6s/hLnPPmfG1hBA$cM3yYeLQGEfIR10Xq+8TfjW6lXnhRcUu/pyGF0Libtk',
    'active',
    'admin',
    NOW(),
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE SET
    full_name = EXCLUDED.full_name,
    phone = EXCLUDED.phone,
    password = EXCLUDED.password,
    status = EXCLUDED.status,
    role = EXCLUDED.role,
    email_verified_at = COALESCE(users.email_verified_at, NOW()),
    updated_at = NOW();
