<?php

namespace App\Http\Controllers;

use Dompdf\Dompdf;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;
use SimpleSoftwareIO\QrCode\Facades\QrCode;
use Aws\S3\S3Client;

class deliveryOrderController extends Controller
{
    public function getPrint(Request $request)
    {
        $hostname = env("URL_STORAGE", "http://storage.edenfarm.tech:8080/v1");
        $qrCode = QrCode::generate( $request['do']['code']);
        $pdf = new Dompdf();
        $options = $pdf->getOptions();
        $options->isPhpEnabled(true);
        $pdf->setOptions($options);
        $pdf->setPaper(
            'A4',
            'portrait'
        );
        $path = base_path();
        $info = file_get_contents($path."/public/img/info.png");
        $logo = file_get_contents($path."/public/img/LogoEden.png");
        $digstamp = file_get_contents($path."/public/img/digistamp.png");
        $data = array(
            'qrCode'       => base64_encode($qrCode),
            'logo' => base64_encode($logo),
            'digstamp' => base64_encode($digstamp),
            'order'     => $request['do'],
            'setting'    => $request['config'],
            'orderDate' => date_create($request['do']['sales_order']['recognition_date']),
            'deliveryDate' => date_create($request['do']['recognition_date']),
            'deliveryKoli' => $request['do']['delivery_koli'],
            'session' => $request['session'],
            'info' => base64_encode($info),
        );

        $html = view('do',$data)->render();

        $pdf->loadHtml($html);
        $pdf->render();
//       $pdf->setOptions(['defaultFont' => 'Open Sans', 'isRemoteEnabled'=> true, 'isPhpEnabled'=>true]);

        $x = 510;
        $y = 810;
        $text = "Page {PAGE_NUM} of {PAGE_COUNT}";
        $font = $pdf->getFontMetrics()->get_font("helvetica");
        $size = 9;
        $color = array(0,0,0);
        $word_space = 0.0;  //  default
        $char_space = 0.0;  //  default
        $angle = 0.0;   //  default
        $pdf->getCanvas()->page_text($x, $y, $text,$font, $size, $color, $word_space, $char_space, $angle);


        $content = $pdf->output();
        $chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz";
        $str = substr(str_shuffle($chars), 0, 5);
        $t=time();

        $filename= $data['order']['code'].'_'.$t.$str.'.pdf';
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
