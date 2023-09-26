<?php

namespace App\Http\Controllers;

use Dompdf\Dompdf;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;
use SimpleSoftwareIO\QrCode\Facades\QrCode;
use Aws\S3\S3Client;

class healthCheckController extends Controller
{
    public function index()
    { 
        $data = array(
            'code'=> 200,
            'status'=> 'success',
            'message'=> 'erp-print-service',
        );
 
        return  $data;
    }


}
