<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <link rel="stylesheet" href="/assets/css/print.css">
    <style>
        @page { margin: 15px; }
        body {
            font-family: Arial, Helvetica, sans-serif;
        }

        td {
            border: 1px solid black;
            font-size: 10px;
        }
        .media-print{
            display: inline-block;
            width: 50%;
        }
        .column {
            margin-top:-30px;
            margin-left:-25px;
            width: 65mm;
            height:17mm
        }
        /* Clear floats after the columns */
        .row:after {
        content: "";
        display: table;
        clear: both;
        }
        .table1 {
            width:100%;
            margin-top:8px;
            border:0px
        }
        .td1 {
            font-size:12px;font-weight:bold;border:0px
        }
        .td2 {
            font-size:12px;font-weight:bold;text-align:right;border:0px
        }
        .td12 {
            font-size:10px;border:0px
        }
        .td22 {
            font-size:10px;text-align:right;border:0px
        }
        .td222 {
            font-size:10px;text-align:right;border:0px;font-weight:bold;
        }
        .border0 {
            border:0px
        }
    </style>
</head>
<body>
    <div>
        <div style="font-size:20px;font-weight:bold">
            {{ $company }}
        </div>
        <div class="td12">
            <b>
                {{ $site['name'] }}
            </b>
        </div>
        <div class="td12">
            {{ $site['street_address'] }}
        </div>
        <div class="td12" style="margin-top:6px">
            No. Handphone: 0{{ $site['phone_number'] }}
        </div>
        <div style="border-top: 1px solid black;margin-top:12px;margin-bottom:12px">
        </div>
        <div class="td12">
            <b>
                {{ $invoice['code'] }}
            </b>
        </div>
        <div style="border-top: 1px dotted black;margin-top:10px;margin-bottom:10px">
        </div>
        <?php
            function tgl_indo($tanggal){
                $bulan = array (
                    1 =>   'Januari',
                    'Februari',
                    'Maret',
                    'April',
                    'Mei',
                    'Juni',
                    'Juli',
                    'Agustus',
                    'September',
                    'Oktober',
                    'November',
                    'Desember'
                );
                $pecahkan = explode('-', $tanggal);
                
                // variabel pecahkan 0 = tanggal
                // variabel pecahkan 1 = bulan
                // variabel pecahkan 2 = tahun
            
                return $pecahkan[2] . ' ' . $bulan[ (int)$pecahkan[1] ] . ' ' . $pecahkan[0];
            }
        ?>
        <table border="0" cellspacing="0" width="100%">
            <tbody>
                <tr>
                    <td style="padding-top:4px;padding-bottom:4px" class="td12">Nota</td>
                    <td style="padding-top:4px;padding-bottom:4px" class="td222">{{ $invoice['sales_invoice']['code'] }}</td>
                </tr>
                <tr>
                    <td style="padding-top:4px;padding-bottom:4px" class="td12">Mitra</td>
                    <td style="padding-top:4px;padding-bottom:4px" class="td222">{{ $branch['pic_name'] }}</td>
                </tr>
                <tr>
                    <td style="padding-top:4px;padding-bottom:4px" class="td12">Metode Pembayaran</td>
                    <td style="padding-top:4px;padding-bottom:4px" class="td222">{{ $invoice['payment_method']['name'] }}</td>
                </tr>
                <tr>
                    <td style="padding-top:4px;padding-bottom:4px" class="td12">Tanggal Bayar</td>
                    <td style="padding-top:4px;padding-bottom:4px" class="td222">
                        <?php
                            $source = $invoice['received_date'];
                            $date = new DateTime($source);
                            echo tgl_indo($date->format('Y-m-d'));
                        ?>
                    </td>
                </tr>
                <tr>
                    <td style="padding-top:4px;padding-bottom:4px" class="td12">Sisa Hutang</td>
                    <td style="padding-top:4px;padding-bottom:4px" class="td222">
                        @if($sisa_hutang=="LUNAS")
                        LUNAS
                        @else
                        <?php
                        $int = (int)$sisa_hutang;
                        ?>
                        {{number_format($int, 2,",",".")}}
                        @endif
                    </td>
                </tr>
                <tr>
                    <td style="padding-top:4px;padding-bottom:4px;font-size:11px" class="td12">Total Pembayaran</td>
                    <td style="padding-top:4px;padding-bottom:4px" class="td222">
                        Rp{{ number_format($invoice['amount'], 2,",",".") }}
                    </td>
                </tr>
            </tbody>
        </table>
        <div style="font-size:10px">
            <b>
            <?php            
                function penyebut($nilai) {
                    $nilai = abs($nilai);
                    $huruf = array("", "Satu", "Dua", "Tiga", "Empat", "Lima", "Enam", "Tujuh", "Delapan", "Sembilan", "Sepuluh", "Sebelas");
                    $temp = "";
                    if ($nilai < 12) {
                        $temp = " ". $huruf[$nilai];
                    } else if ($nilai <20) {
                        $temp = penyebut($nilai - 10). " Belas";
                    } else if ($nilai < 100) {
                        $temp = penyebut($nilai/10)." Puluh". penyebut($nilai % 10);
                    } else if ($nilai < 200) {
                        $temp = " Seratus" . penyebut($nilai - 100);
                    } else if ($nilai < 1000) {
                        $temp = penyebut($nilai/100) . " Ratus" . penyebut($nilai % 100);
                    } else if ($nilai < 2000) {
                        $temp = " Seribu" . penyebut($nilai - 1000);
                    } else if ($nilai < 1000000) {
                        $temp = penyebut($nilai/1000) . " Ribu" . penyebut($nilai % 1000);
                    } else if ($nilai < 1000000000) {
                        $temp = penyebut($nilai/1000000) . " Juta" . penyebut($nilai % 1000000);
                    } else if ($nilai < 1000000000000) {
                        $temp = penyebut($nilai/1000000000) . " Milyar" . penyebut(fmod($nilai,1000000000));
                    } else if ($nilai < 1000000000000000) {
                        $temp = penyebut($nilai/1000000000000) . " Trilyun" . penyebut(fmod($nilai,1000000000000));
                    }     
                    return $temp;
                }
            
                function terbilang($nilai) {
                    if($nilai<0) {
                        $hasil = "minus ". trim(penyebut($nilai));
                    } else {
                        $hasil = trim(penyebut($nilai));
                    }     		
                    return $hasil;
                }
            
                $angka = $invoice['amount'];
                echo terbilang($angka);
            ?>
            Rupiah
            </b>
        </div>
        <div style="border-top: 1px dotted black;margin-top:10px;margin-bottom:10px">
        </div>
        <div class="td12" style="margin-top:10px;">
            @if($invoice['note'])
            <strong>Catatan :</strong> <br> <i>{{ $invoice['note'] }}</i>
            @endif
        </div>
        <div style="margin-top:30px;">
            <table width="100%">
                <tr>
                    <td style="border:0px;font-size:10px;">
                    <img src='data:image/png;base64," . {{ $logo }} . "'>
                    <?php
                    echo tgl_indo($dateNow);
                    ?>
                    ({{$timeNow}})
                    </td>
                    <td style="border:0px;text-align:right;font-size:10px;">
                        <div>
                            Hormat Kami,
                        </div>
                        <div style="font-weight:bold;font-size:14px;">
                            {{ $company }}
                        </div>
                    </td>
                </tr>
            </table>
            
        </div>
    </div>
<script type="text/javascript">
    window.onload = function () {
        window.print();
        window.close();
    };
    
</script>
</body>
</html>
