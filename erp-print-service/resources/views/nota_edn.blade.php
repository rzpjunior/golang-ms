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
                <td class="td1" width=50%>{{ $invoice['sales_order']['site']['name'] }}</td>
                <td class="td2">{{ $invoice['code'] }}</td>
            </tr>
            <tr>
                <td class="td12">{{ $invoice['sales_order']['site']['street_address'] }}</td>
                <td class="td22">Mitra : {{ $invoice['sales_order']['branch']['pic_name'] }}</td>
            </tr>
            <tr>
                <td class="td12">No. Handphone: 0{{ $invoice['sales_order']['site']['phone_number'] }}</td>
                <td class="td22">
                Nota :
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
        </table>
        <table border="1" style="border-collapse: collapse; font-size:12px;margin-top:20px" width="100%">
            <thead>
            <tr>
                <th style="padding:10px;border:1px solid black;text-align: center;" width="5%">No</th>
                <th style="padding:10px;border:1px solid black; text-align: center;" width="45%">Product</th>
                <th style="padding:10px;border:1px solid black;" width="15%">Qty</th>
                <th style="padding:10px;border:1px solid black;" width="15%">Harga Per Unit</th>
                <th style="padding:10px;border:1px solid black;text-align: right;" width="20%">Subtotal</th>
            </tr>
            </thead>
            <tbody>
            @foreach($invoice['sales_invoice_items'] as $index => $items)
                <tr>
                    <td style="padding:15px;text-align:center;">{{$index + 1}}</td>
                    <td style="padding:15px;">{{ $items['item']['name']}}</td>
                    <td style="padding:15px;text-align:right;">
                        {{ number_format($items['invoice_qty'], 2,',','.') }} {{ $items['item']['uom']['name'] }}
                    </td>
                    <td style="padding:15px; text-align:right;">{{ number_format($items['unit_price'], 2,",",".") }}</td>
                    <td style="padding:15px; text-align:right;">{{ number_format($items['subtotal'], 2,",",".") }}</td>
                </tr>
            @endforeach
            </tbody>
        </table>
        
        <table width=100% style="margin-top:15px">
            <tr>
                <td width=60% class="border0" style="vertical-align:top">
                    <table>
                        <tr>
                            <td class="border0">
                                @if($invoice['note'])
                                <strong>Catatan :</strong> <br> <i>{{ $invoice['note'] }}</i>
                                @endif
                            </td>
                        </tr>
                    </table>
                </td>
                <td class="border0" width=40%>
                    <table style="border-collapse: collapse; font-size:12px; font-weight:bold" width="100% float:right;" border="0">
                        <tr>
                            <td class="border0" style="margin-left:25px">Total Pembayaran</td>
                            <td style="text-align: right;margin-right:20px" class="border0">
                                <span>
                                    {{ number_format($invoice['total_charge'], 2,",",".") }}
                                </span>
                            </td>
                        </tr>
                    </table>
                    <div style="border-top: 1px solid black;margin-top:10px">
                    </div>
                </td>
            </tr>
        </table>
        <div style="margin-top:100px;">
            <table width="100%">
                <tr>
                    <td style="border:0px;font-size:10px;">
                    <img src='data:image/png;base64," . {{ $logo }} . "'>
                    <?php
                    echo tgl_indo($dateNow);
                    ?>
                    ({{$timeNow}})
                    </td>
                    <td style="border:0px;">
                        <div style="text-align:right;font-size:12px;">
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
