<?php

namespace App\Http\Controllers;

use App\Services\SatuSehatOAuth2Service;
use Illuminate\Http\JsonResponse;
use Exception;

class OAuth2Controller extends Controller
{
    protected SatuSehatOAuth2Service $oauthService;

    public function __construct(SatuSehatOAuth2Service $oauthService)
    {
        $this->oauthService = $oauthService;
    }

    /**
     * Get access token dari SATUSEHAT
     * 
     * @return JsonResponse
     */
    public function getToken(): JsonResponse
    {
        try {
            $result = $this->oauthService->getAccessToken();

            return response()->json([
                'success' => true,
                'data' => [
                    'access_token' => $result['access_token'],
                    'expires_in' => $result['expires_in'] ?? 3600,
                    'from_cache' => $result['from_cache'] ?? false,
                    'token_type' => 'Bearer'
                ],
                'message' => 'Access token berhasil diperoleh'
            ]);

        } catch (Exception $e) {
            \Log::error('OAuth2 Token Error: ' . $e->getMessage());

            return response()->json([
                'success' => false,
                'error' => $e->getMessage(),
                'message' => 'Gagal mendapatkan access token'
            ], 500);
        }
    }

    /**
     * Refresh access token (force new request)
     * 
     * @return JsonResponse
     */
    public function refreshToken(): JsonResponse
    {
        try {
            $result = $this->oauthService->refreshAccessToken();

            return response()->json([
                'success' => true,
                'data' => [
                    'access_token' => $result['access_token'],
                    'expires_in' => $result['expires_in'] ?? 3600,
                    'token_type' => 'Bearer'
                ],
                'message' => 'Access token berhasil di-refresh'
            ]);

        } catch (Exception $e) {
            \Log::error('OAuth2 Refresh Error: ' . $e->getMessage());

            return response()->json([
                'success' => false,
                'error' => $e->getMessage(),
                'message' => 'Gagal me-refresh access token'
            ], 500);
        }
    }

    /**
     * Validate access token
     * 
     * @return JsonResponse
     */
    public function validateToken(): JsonResponse
    {
        try {
            $token = request()->bearerToken();

            if (!$token) {
                return response()->json([
                    'success' => false,
                    'message' => 'Token tidak ditemukan'
                ], 401);
            }

            $isValid = $this->oauthService->validateToken($token);

            return response()->json([
                'success' => true,
                'data' => [
                    'is_valid' => $isValid
                ],
                'message' => $isValid ? 'Token valid' : 'Token tidak valid'
            ]);

        } catch (Exception $e) {
            \Log::error('Token Validation Error: ' . $e->getMessage());

            return response()->json([
                'success' => false,
                'error' => $e->getMessage(),
                'message' => 'Gagal validasi token'
            ], 500);
        }
    }
}
