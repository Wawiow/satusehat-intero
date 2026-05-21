<?php

namespace App\Services;

use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Cache;
use Exception;

class SatuSehatOAuth2Service
{
    protected string $clientId;
    protected string $clientSecret;
    protected string $oauthUrl;
    protected string $apiUrl;

    public function __construct()
    {
        $this->clientId = config('services.satusehat.client_id');
        $this->clientSecret = config('services.satusehat.client_secret');
        $this->oauthUrl = config('services.satusehat.oauth_url');
        $this->apiUrl = config('services.satusehat.api_url');
    }

    /**
     * Get access token from SATUSEHAT OAuth2 server
     * Token akan di-cache untuk menghindari multiple requests
     * 
     * @return array
     * @throws Exception
     */
    public function getAccessToken(): array
    {
        // Cek apakah token sudah ada di cache
        $cacheKey = 'satusehat_access_token';
        
        if (Cache::has($cacheKey)) {
            return [
                'access_token' => Cache::get($cacheKey),
                'from_cache' => true
            ];
        }

        try {
            $response = Http::post("{$this->oauthUrl}/accesstoken", [
                'client_id' => $this->clientId,
                'client_secret' => $this->clientSecret,
                'grant_type' => 'client_credentials'
            ]);

            if (!$response->successful()) {
                throw new Exception("Failed to get access token: " . $response->body());
            }

            $data = $response->json();
            $accessToken = $data['access_token'] ?? null;
            $expiresIn = $data['expires_in'] ?? 3600; // Default 1 jam jika tidak ada

            if (!$accessToken) {
                throw new Exception("No access token in response");
            }

            // Cache token dengan TTL yang sesuai
            // Kurangi 300 detik (5 menit) untuk refresh sebelum expire
            Cache::put($cacheKey, $accessToken, $expiresIn - 300);

            return [
                'access_token' => $accessToken,
                'expires_in' => $expiresIn,
                'from_cache' => false
            ];

        } catch (Exception $e) {
            \Log::error('SATUSEHAT OAuth2 Error: ' . $e->getMessage());
            throw $e;
        }
    }

    /**
     * Get access token (simple version, return string)
     * 
     * @return string
     * @throws Exception
     */
    public function getToken(): string
    {
        return $this->getAccessToken()['access_token'];
    }

    /**
     * Refresh access token (hapus dari cache untuk force new request)
     * 
     * @return array
     * @throws Exception
     */
    public function refreshAccessToken(): array
    {
        Cache::forget('satusehat_access_token');
        return $this->getAccessToken();
    }

    /**
     * Validate access token
     * 
     * @param string $token
     * @return bool
     */
    public function validateToken(string $token): bool
    {
        try {
            $response = Http::withToken($token)->get("{$this->apiUrl}/v1/user/profile");
            return $response->successful();
        } catch (Exception $e) {
            \Log::error('Token validation error: ' . $e->getMessage());
            return false;
        }
    }
}
