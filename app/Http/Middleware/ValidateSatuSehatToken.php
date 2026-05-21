<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use App\Services\SatuSehatOAuth2Service;

class ValidateSatuSehatToken
{
    public function __construct(private SatuSehatOAuth2Service $oauthService)
    {
    }

    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure(\Illuminate\Http\Request): (\Illuminate\Http\Response|\Illuminate\Http\RedirectResponse)  $next
     * @return \Illuminate\Http\Response|\Illuminate\Http\RedirectResponse
     */
    public function handle(Request $request, Closure $next)
    {
        $token = $request->bearerToken();

        if (!$token) {
            return response()->json([
                'success' => false,
                'message' => 'Token tidak ditemukan'
            ], 401);
        }

        // Jika token invalid, refresh dan retry
        if (!$this->oauthService->validateToken($token)) {
            return response()->json([
                'success' => false,
                'message' => 'Token tidak valid atau expired'
            ], 401);
        }

        $request->attributes->add(['satusehat_token' => $token]);

        return $next($request);
    }
}
