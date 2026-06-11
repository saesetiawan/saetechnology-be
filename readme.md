# Saa Technology Solution Backend

Backend ini difokuskan untuk kebutuhan website promosi jasa, custom software, custom apps, dan langganan SaaS.

## Fokus utama

- Auth dan akun pengguna
- Website settings
- Website contents
- Upload file untuk kebutuhan konten

## Endpoint aktif

- `/api/login`
- `/api/register`
- `/api/refresh-token`
- `/api/logout`
- `/api/activate-account`
- `/api/customer/login`
- `/api/customer/register`
- `/api/profile`
- `/api/customer/profile`
- `/api/customer/password`
- `/api/storage/upload-file`
- `/api/public/website-contents`
- `/api/website-contents`
- `/api/public/website-settings`
- `/api/website-settings`

## Catatan

Modul marketplace lama seperti order, cart, payment, courier, product catalog, seller/store, admin dashboard, dan notifikasi sudah dihapus dari handler, usecase, repository, domain, migration, dan infrastructure.
