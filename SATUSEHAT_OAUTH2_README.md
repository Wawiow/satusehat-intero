# SATUSEHAT OAuth2 Backend Documentation

## Setup

### 1. Update Environment Variables

Edit file `.env` dan update konfigurasi SATUSEHAT:

```env
SATUSEHAT_ENV=sandbox
SATUSEHAT_CLIENT_ID=your_client_id_here
SATUSEHAT_CLIENT_SECRET=your_client_secret_here
SATUSEHAT_OAUTH_URL=https://api-sandbox.dto.kemkes.go.id/oauth2/v1
SATUSEHAT_API_URL=https://api-sandbox.dto.kemkes.go.id
```

**Ganti nilai di atas dengan credentials SATUSEHAT Anda.**

### 2. Struktur File

- `app/Services/SatuSehatOAuth2Service.php` - Service untuk OAuth2 logic
- `app/Http/Controllers/OAuth2Controller.php` - Controller untuk endpoints
- `config/satusehat.php` - Konfigurasi SATUSEHAT
- `routes/api.php` - API routes

## API Endpoints

### 1. Get Access Token
**Endpoint:** `POST /api/oauth2/token`

Request:
```bash
curl -X POST http://localhost:8000/api/oauth2/token
```

Response:
```json
{
    "success": true,
    "data": {
        "access_token": "eyJhbGc...",
        "expires_in": 3600,
        "from_cache": false,
        "token_type": "Bearer"
    },
    "message": "Access token berhasil diperoleh"
}
```

### 2. Refresh Access Token
**Endpoint:** `POST /api/oauth2/refresh`

Merefresh token yang sudah ada (force new request dari server):

```bash
curl -X POST http://localhost:8000/api/oauth2/refresh
```

Response:
```json
{
    "success": true,
    "data": {
        "access_token": "eyJhbGc...",
        "expires_in": 3600,
        "token_type": "Bearer"
    },
    "message": "Access token berhasil di-refresh"
}
```

### 3. Validate Token
**Endpoint:** `POST /api/oauth2/validate`

Validasi apakah token masih valid:

```bash
curl -X POST http://localhost:8000/api/oauth2/validate \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

Response:
```json
{
    "success": true,
    "data": {
        "is_valid": true
    },
    "message": "Token valid"
}
```

## Features

### ✅ Token Caching
- Token secara otomatis di-cache untuk menghindari multiple requests
- Cache TTL otomatis berdasarkan `expires_in` dari server (dikurangi 5 menit)
- Lebih efficient dan mengurangi beban API

### ✅ Error Handling
- Exception handling yang robust
- Logging untuk debugging
- Response yang consistent

### ✅ Client Credentials Flow
- Menggunakan OAuth2 Client Credentials flow
- Cocok untuk backend-to-backend communication

## Usage dalam Code

### Mendapatkan Token di Controller/Service:

```php
use App\Services\SatuSehatOAuth2Service;

class MyController extends Controller
{
    public function __construct(private SatuSehatOAuth2Service $oauthService)
    {
    }

    public function myMethod()
    {
        try {
            // Get token (dari cache jika ada)
            $token = $this->oauthService->getToken();
            
            // Atau dapatkan lengkap dengan metadata
            $result = $this->oauthService->getAccessToken();
            
            // Gunakan token untuk API calls
            $response = Http::withToken($token)->get('https://api.satusehat.com/endpoint');
            
        } catch (Exception $e) {
            // Handle error
        }
    }
}
```

## Best Practices

1. **Selalu gunakan try-catch** saat memanggil OAuth service
2. **Monitor logs** untuk debugging issues
3. **Cache token** untuk mengurangi API calls (sudah otomatis dilakukan)
4. **Refresh token** jika mendapat error 401 (Unauthorized)
5. **Jangan hardcode** credentials - selalu gunakan environment variables

## Troubleshooting

### Error: "Failed to get access token"
- Cek credentials di `.env` apakah sudah benar
- Pastikan internet connection aktif
- Check logs: `storage/logs/laravel.log`

### Token tidak di-cache
- Cek apakah `CACHE_STORE` di `.env` sudah dikonfigurasi
- Default menggunakan database cache

### 401 Unauthorized
- Refresh token menggunakan endpoint `/api/oauth2/refresh`
- Atau restart aplikasi untuk clear cache

## Production Deployment

Untuk production, pastikan:
1. Update `.env` dengan credentials production SATUSEHAT
2. Set `SATUSEHAT_ENV=production`
3. Update URLs ke production endpoints
4. Monitor token expiration dan refresh mechanisms
5. Implement proper security untuk token storage
