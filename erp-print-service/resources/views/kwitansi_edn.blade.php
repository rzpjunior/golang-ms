<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <link rel="stylesheet" href="/assets/css/print.css">
    <style>
        @page { margin: 2px; }
        body {
            font-family: Arial, Helvetica, sans-serif;
            margin: 1px;
            margin-top:5px;
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
        .border0 {
            border:0px
        }
    </style>
</head>
<body>
<div style="padding:24px">
        <div style="font-size:24px;font-weight:bold">
        {{ $company }}
        </div>
        <table class="table1" border="0" cellspacing="0">
            <tr> 
                <td class="td12" width=50% style="font-weight:700">{{ $site['name'] }}</td>
                <td class="td22">Kwitansi : <span style="font-weight:700">{{ $invoice['code'] }}</span></td>
            </tr>
            <tr>
                <td class="td12">{{ $site['street_address'] }}</td>
                <td class="td22">Nota : <span style="font-weight:700">{{ $invoice['sales_invoice']['code'] }}</span></td>
            </tr>
            <tr>
                <td class="td12">No. Handphone: 0{{ $site['phone_number'] }}</td>
                <td class="td22"></td>
            </tr>
        </table>
        <table width="100%" border="0" cellspacing="0" style="margin-top:20px">
            <tbody>
                <tr>
                    <td style="padding:4px;" width="50%">
                        <table width="100%">
                            <tr>
                                <td width=48% class="border0">Mitra</td>
                                <td width=3% class="border0">:</td>
                                <td width=49% class="border0">{{ $branch['pic_name'] }}</td>
                            </tr>
                        </table>
                    </td>
                    <td style="padding:4px;" width="50%">
                        <table width="100%">
                            <tr>
                                <td width=48% class="border0">Tanggal Bayar</td>
                                <td width=3% class="border0">:</td>
                                <td width=49% class="border0">
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
                                        
                                        return $pecahkan[2] . ' ' . $bulan[ (int)$pecahkan[1] ] . ' ' . $pecahkan[0];
                                    }
                                    $source = $invoice['received_date'];
                                    $date = new DateTime($source);
                                    echo tgl_indo($date->format('Y-m-d'));
                                ?>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
                <tr>
                    <td style="padding:4px;" width="50%">
                        <table width="100%">
                            <tr>
                                <td width=48% class="border0">Metode Pembayaran</td>
                                <td width=3% class="border0">:</td>
                                <td width=49% class="border0">{{ $invoice['payment_method']['name'] }}</td>
                            </tr>
                        </table>
                    </td>
                    <td style="padding:4px;" width="50%">
                        <table width="100%">
                            <tr>
                                <td width=48% class="border0">Sisa Hutang</td>
                                <td width=3% class="border0">:</td>
                                <td width=49% class="border0">
                                    @if($sisa_hutang=="LUNAS")
                                    LUNAS
                                    @else
                                    <?php
                                    $int = (int)$sisa_hutang;
                                    ?>
                                    Rp{{number_format($int, 2,",",".")}}
                                    @endif
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
                <tr>
                    <td colspan="2" style="padding:4px;">
                        <table width="100%">
                            <tr>
                                <td width=23.5% class="border0">Catatan</td>
                                <td width=1.5% class="border0">:</td>
                                <td class="border0">
                                @if($invoice['note'])
                                {{$invoice['note']}}
                                @else
                                -
                                @endif
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
                <tr>
                    <td colspan="2" style="padding:4px;">
                        <table width="100%">
                            <tr>
                                <td width=23.5% class="border0">Total Pembayaran</td>
                                <td width=1.5% class="border0">:</td>
                                <td class="border0">
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
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </tbody>
        </table>
        
        <table width=100% style="margin-top:20px">
            <tr>
                <td width=30% style="text-align:center;font-size:14px">
                    <b>
                        Rp{{ number_format($invoice['amount'], 2,",",".") }}
                    </b>
                </td>
                <td class="border0" >
                    <div style="margin-top:1px;text-align:right;font-size:12px;">
                        <div>
                            Hormat Kami,
                        </div>
                        <div style="font-weight:bold;font-size:14px;margin-top:8px">
                        {{ $company }}
                        </div>
                    </div>
                </td>
            </tr>
        </table>
        <div style="font-size:10px;margin-top:10px">
            <img src='data:image/png;base64," . {{ $logo }} . "'>
            <?php
            echo tgl_indo($dateNow);
            ?>
            ({{$timeNow}})
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
