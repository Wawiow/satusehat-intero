<?php

namespace Tests\Feature;

use Tests\TestCase;
use App\Services\SatuSehatOAuth2Service;

class OAuth2Test extends TestCase
{
    /**
     * Test mendapatkan access token
     */
    public function test_can_get_access_token(): void
    {
        $response = $this->postJson('/api/oauth2/token');

        $response->assertStatus(200)
            ->assertJsonStructure([
                'success',
                'data' => [
                    'access_token',
                    'expires_in',
                    'from_cache',
                    'token_type'
                ],
                'message'
            ]);

        $this->assertTrue($response->json('success'));
        $this->assertNotEmpty($response->json('data.access_token'));
    }

    /**
     * Test refresh access token
     */
    public function test_can_refresh_access_token(): void
    {
        $response = $this->postJson('/api/oauth2/refresh');

        $response->assertStatus(200)
            ->assertJsonStructure([
                'success',
                'data' => [
                    'access_token',
                    'expires_in',
                    'token_type'
                ],
                'message'
            ]);

        $this->assertTrue($response->json('success'));
    }

    /**
     * Test validate token
     */
    public function test_can_validate_token(): void
    {
        // Pertama dapatkan token
        $getTokenResponse = $this->postJson('/api/oauth2/token');
        $token = $getTokenResponse->json('data.access_token');

        // Kemudian validate token
        $response = $this->postJson('/api/oauth2/validate', [], [
            'Authorization' => "Bearer {$token}"
        ]);

        $response->assertStatus(200)
            ->assertJsonStructure([
                'success',
                'data' => [
                    'is_valid'
                ],
                'message'
            ]);
    }

    /**
     * Test validate dengan token kosong
     */
    public function test_validate_without_token(): void
    {
        $response = $this->postJson('/api/oauth2/validate');

        $response->assertStatus(401);
        $this->assertFalse($response->json('success'));
    }

    /**
     * Test service getToken method
     */
    public function test_oauth_service_get_token(): void
    {
        $service = app(SatuSehatOAuth2Service::class);
        
        $result = $service->getAccessToken();

        $this->assertIsArray($result);
        $this->assertArrayHasKey('access_token', $result);
        $this->assertArrayHasKey('expires_in', $result);
        $this->assertArrayHasKey('from_cache', $result);
        $this->assertNotEmpty($result['access_token']);
    }
}
