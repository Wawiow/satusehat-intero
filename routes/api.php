<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\OAuth2Controller;

// OAuth2 Routes
Route::prefix('oauth2')->group(function () {
    Route::post('/token', [OAuth2Controller::class, 'getToken'])->name('oauth2.token');
    Route::post('/refresh', [OAuth2Controller::class, 'refreshToken'])->name('oauth2.refresh');
    Route::post('/validate', [OAuth2Controller::class, 'validateToken'])->name('oauth2.validate');
});
