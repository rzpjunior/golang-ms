<?php

namespace App\Http\Controllers;

use Dompdf\Dompdf;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;
use SimpleSoftwareIO\QrCode\Facades\QrCode;
use Aws\S3\S3Client;
use Illuminate\Support\Str;
use Carbon\Carbon;

class kwitansiEDNController extends Controller
{


    public function getPrint(Request $request)
    {
        $date = Carbon::now()->timezone('Asia/Phnom_Penh');

        $pdf = new Dompdf();
        $options = $pdf->getOptions();
        $options->isPhpEnabled(true);
        $pdf->setOptions($options);
        $pdf->setPaper(
            'A6',
            'landscape'
        );
        $path = base_path();
        $logo = file_get_contents($path."/public/img/print.png");
        $temp = array(
            'company'     => $request['company'],
            'sisa_hutang'     => $request['sisa_hutang'],
            'invoice'     => $request['sp'],
            'branch'     => $request['branch'],
            'site'     => $request['site'],
            'logo'      => base64_encode($logo),
            'dateNow'      => $date->format("Y-m-d"),
            'timeNow'      => $date->format("H:i")
        );
        $html = view('kwitansi_edn',$temp)->render();

        $pdf->loadHtml($html);
        $pdf->render();
        $content = $pdf->output();
        $chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz";
            $str = substr(str_shuffle($chars), 0, 5);
            $t=time();
            $filename= $request['sp']['code'].'_'.$t.$str.'.pdf';
        file_put_contents($filename,$content);

        Storage::disk('s3EDN')->put($filename, file_get_contents($path.'/public/'.$filename,'public'));

        $client = new S3Client([
            'version' => 'latest',
            'region'  => env('AWS_DEFAULT_REGION'),
            'endpoint' => env('AWS_ENDPOINT'),
            'credentials' => [
                'key'    => env('AWS_ACCESS_KEY_ID'),
                'secret' => env('AWS_SECRET_ACCESS_KEY'),
            ],
        ]);
        $request = $client->getObjectUrl(env('AWS_BUCKET'),$filename);
        File::delete($path.'/public/'.$filename);

        $datas = array('data' => $request);

        return  $datas;
    }


}
