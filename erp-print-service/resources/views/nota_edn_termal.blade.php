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
                {{ $invoice['sales_order']['site']['name'] }}
            </b>
        </div>
        <div class="td12">
            {{ $invoice['sales_order']['site']['street_address'] }}
        </div>
        <div class="td12" style="margin-top:6px">
            No. Handphone: 0{{ $invoice['sales_order']['site']['phone_number'] }}
        </div>
        <div style="border-top: 1px solid black;margin-top:12px;margin-bottom:12px">

        </div>
        <table class="table1" border="0" cellspacing="0">
            <tr>
                <td class="td12" width=50%> <b>{{ $invoice['code'] }}</b></td>
                <td class="td22">
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
                    $source = $invoice['created_at'];
                    $date = new DateTime($source);
                    echo tgl_indo($date->format('Y-m-d'));
                ?>
                </td>
            </tr>
            <tr>
                <td class="td12">Mitra: {{ $invoice['sales_order']['branch']['pic_name'] }}</td>
                <td class="td22"></td>
            </tr>
        </table>
        
        <div style="border-top: 1px dotted black;margin-top:10px;">

        </div>
        <table border="0" cellspacing="0" width="100%">
            <tbody>
            @foreach($invoice['sales_invoice_items'] as $index => $items)
                <tr>
                    <td class="td12" style="">{{ $items['item']['name']}}</td>
                    <td class="td12" style="text-align:right;">
                        {{ number_format($items['invoice_qty'], 2,',','.') }} {{ $items['item']['uom']['name'] }}
                    </td>
                    <td class="td12" style=" text-align:right;">{{ number_format($items['unit_price'], 2,",",".") }}</td>
                    <td class="td12" style=" text-align:right;">{{ number_format($items['subtotal'], 2,",",".") }}</td>
                </tr>
            @endforeach
                <tr style="margin-top:10px;margin-bottom:10px">
                    <td class="td12" style="border-top: 1px dotted black;"></td>
                    <td class="td12" style="border-top: 1px dotted black;"></td>
                    <td class="td12" style="border-top: 1px dotted black;"></td>
                    <td class="td12" style="border-top: 1px dotted black;"></td>
                </tr>
                <tr>
                    <td class="td12" style="">Total Pembayaran</td>
                    <td class="td12" style="text-align:right;"></td>
                    <td class="td12" style=" text-align:right;"></td>
                    <td class="td12" style=" text-align:right;"> <b> {{ number_format($invoice['total_charge'], 2,",",".") }} </b></td>
                </tr>
            </tbody>
        </table>
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
