<?php

namespace App\Http\Controllers;

use Dompdf\Dompdf;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;
use SimpleSoftwareIO\QrCode\Facades\QrCode;
use Aws\S3\S3Client;
use Carbon\Carbon;

class qrCodeLabelController extends Controller
{
    public function getPrint(Request $request)
    {
        $qrCode = QrCode::generate($request['qr']['rfid_code']);
        $date = Carbon::now()->timezone('Asia/Phnom_Penh');

        $pdf = new Dompdf();
        $options = $pdf->getOptions();
        $options->isPhpEnabled(true);
        $pdf->setOptions($options);
        $pdf->setPaper(
            'A7',
            'landscape'
        );
        $path = base_path();
        $data = array(
            'qrCode'       => base64_encode($qrCode),
            'rfidCode'     => $request['qr']['rfid_code']
        );

        $html = view('qrCodeLabel_print',$data)->render();

        $pdf->loadHtml($html);
        $pdf->render();
//       $pdf->setOptions(['defaultFont' => 'Open Sans', 'isRemoteEnabled'=> true, 'isPhpEnabled'=>true]);

        $content = $pdf->output();
        $chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz";
            $str = substr(str_shuffle($chars), 0, 5);
            $t=time();
            $filename= $data['rfidCode'].'_'.$t.$str.'.pdf';
        file_put_contents($filename,$content);

        Storage::disk('s3Public')->put($filename, file_get_contents($path.'/public/'.$filename,'public'));

        $client = new S3Client([
            'version' => 'latest',
            'region'  => env('AWS_DEFAULT_REGION'),
            'endpoint' => env('AWS_ENDPOINT'),
            'credentials' => [
                'key'    => env('AWS_ACCESS_KEY_ID'),
                'secret' => env('AWS_SECRET_ACCESS_KEY'),
            ],
        ]);
        // $cmd = $client->getCommand('GetObject', [
        //     'Bucket' => env('AWS_BUCKET'),
        //     'Key'    => $filename
        // ]);
        $request = $client->getObjectUrl(env('AWS_BUCKET'),$filename);
        // $presignedUrl = (string) $request->getUri();
        File::delete($path.'/public/'.$filename);

        $datas = array('data' => $request);

        return  $datas;
    }


}