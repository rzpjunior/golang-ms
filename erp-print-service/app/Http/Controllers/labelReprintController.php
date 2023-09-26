<?php

namespace App\Http\Controllers;

use Dompdf\Dompdf;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;
use SimpleSoftwareIO\QrCode\Facades\QrCode;
use Aws\S3\S3Client;
use Illuminate\Support\Str;

class labelReprintController extends Controller
{


    public function getPrint(Request $request)
    {
        $pdf = new Dompdf();
        $options = $pdf->getOptions();
        $options->isPhpEnabled(true);
        $pdf->setOptions($options);
        $pdf->setPaper(
            'A7',
            'landscape'
        );
        $plis = $request['plis'];
        $temp = [];
        for($i = 0;$i < count($plis); $i++)
            {
                $branchElipsis = Str::limit($request['plis'][$i]['sales_order']['branch']['name'], 70);
                $qr = QrCode::generate($request['plis'][$i]['sales_order']['code'] . '-' . $request['plis'][$i]['increment']);
                $data = array(
                        'qrCode'       => base64_encode($qr),
                        'picking'     => $request['plis'][$i],
                        'branch'      => $branchElipsis
                    );
                array_push($temp, $data);
            }
        $path = base_path();
        $temp = array(
            'temp' =>  $temp
        );
        $html = view('label_reprint',$temp)->render();

        $pdf->loadHtml($html);
        $pdf->render();
//       $pdf->setOptions(['defaultFont' => 'Open Sans', 'isRemoteEnabled'=> true, 'isPhpEnabled'=>true]);

        $content = $pdf->output();
        $chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz";
            $str = substr(str_shuffle($chars), 0, 5);
            $t=time();
            $filename= $data['picking']['sales_order']['code'].'_'.$t.$str.'.pdf';
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
