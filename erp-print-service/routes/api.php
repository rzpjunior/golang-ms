<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::middleware('auth:api')->get('/user', function (Request $request) {
    return $request->user();
});
Route::post('/read/so', 'salesOrderController@getPrint');
Route::post('/read/do', 'deliveryOrderController@getPrint');
Route::post('/read/si', 'salesInvoiceController@getPrint');
Route::post('/read/po', 'purchaseOrderController@getPrint');
Route::post('/read/wd', 'wasteDisposalController@getPrint');
Route::post('/read/gt', 'goodsTransferController@getPrint');
Route::post('/read/gr', 'goodsReceiptController@getPrint');
Route::post('/read/sr', 'supplierReturnController@getPrint');
Route::post('/read/picking_print', 'pickingPrintController@getPrint');
Route::post('/read/label_product', 'labelProductController@getPrint');
Route::post('/read/label_reprint', 'labelReprintController@getPrint');
Route::post('/read/dn', 'debitNoteController@getPrint');
Route::post('/read/label_packing', 'labelPackingPrintController@getPrint');
Route::post('/read/picking_list_print', 'pickingLabelController@getPrint');
Route::post('/read/nota_edn', 'notaEDNController@getPrint');
Route::post('/read/kwitansi_edn', 'kwitansiEDNController@getPrint');
Route::post('/read/nota_edn_termal', 'notaEDNTermalController@getPrint');
Route::post('/read/kwitansi_edn_termal', 'kwitansiEDNTermalController@getPrint');
Route::post('/read/qrcodeLabel', 'qrCodeLabelController@getPrint');


Route::get('/v1/health_check', 'healthCheckController@index');
