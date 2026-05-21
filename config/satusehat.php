<?php

return [
    'satusehat' => [
        'client_id' => env('SATUSEHAT_CLIENT_ID'),
        'client_secret' => env('SATUSEHAT_CLIENT_SECRET'),
        'oauth_url' => env('SATUSEHAT_OAUTH_URL', 'https://api-sandbox.dto.kemkes.go.id/oauth2/v1'),
        'api_url' => env('SATUSEHAT_API_URL', 'https://api-sandbox.dto.kemkes.go.id'),
        'environment' => env('SATUSEHAT_ENV', 'sandbox'),
    ],
];
