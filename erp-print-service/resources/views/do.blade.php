<!DOCTYPE html>
<html lang="en">
<head>
    <title>{{ $order['code'] }}</title>
    <meta charset="utf-8">
    <link rel="stylesheet" href="/assets/css/print.css">
    <style>
        body {
            font-family: Arial, Helvetica, sans-serif;
            font-size: 14px;
        }

        th {
            text-align: center;
        }

        table.title {
            width: 100%;
            line-height: 10%;
        }

        table.company {
            width: 60%;
            font-size: 12px;
        }

        table.customer {
            font-size: 12px;
            width: 100%;
            margin-bottom: 20px;
        }

        .customer td {
            padding: 5px;
        }

        table.item {
            text-align: center;
        }

        .item td {
            padding-left: 5px;
            padding-right: 5px;
        }

        table.footer {
            width: 60%;
            /*margin-bottom:30px;*/
        }

        .footer td {
            padding-left: 5px;
            padding-top: 3px;
        }

        .company_name {
            font-size: 15px;
            font-weight: bold;
        }

        .total td {
            padding-right: 5px;
        }


        .second-color {
            color: #768F9C;
        }

        footer {
            position: fixed;
            left: 0;
            bottom: 0;
            width: 100%;
            text-align: left;
            font-size: 12px;
        }
        delta-print {
            position: fixed;
            left: 0;
            bottom: 0;
            width: 100%;
            text-align: center;
            font-size: 12px;
        }
        /* Create two equal columns that floats next to each other */
        .column {
            float: left;
            width: 50%;
            padding: 10px;
            height: 300px; /* Should be removed. Only for demonstration */
        }

        /* Clear floats after the columns */
        .row:after {
            content: "";
            display: table;
            clear: both;
        }

        .sign{
            width:100%;
            text-align:center;
            float:right;
            margin-top:50px;
        }
        
        .self-pickup {
            position: fixed;
            font-weight: 700;
            font-size: 22px;
            line-height: 26px;
            text-align: center;
            border: 1px solid #333333;
            padding: 8px;
            left: 50%;
            transform: translate(-50%, 0);
        }
    </style>
</head>
<body>
<footer>
    <span>#{{ $order['code'] }}</span>
</footer>
<delta-print>
    <span>X{{ $order['delta_print'] }}-{{date("dmYhis")}}-{{$session}}</span>
</delta-print>
<section class="sheet" style="padding: 10px;">
    @if($order['sales_order']['order_type']['id'] == "393216")
    <div class="self-pickup">
        {{ $order['sales_order']['term_payment_sls']['name'] }} - {{ $order['sales_order']['order_type']['name'] }}
    </div>
    @endif
    <table class="title" border="0">
        <tr>
            <td width=70%>
                <img src='data:image/png;base64," . {{ $logo }} . "' width="125" height="100" style="margin: -12px;">
            </td>
            <td width=30% style="text-align:right;">
                <img src='data:image/png;base64," . {{ $qrCode }} . "' width="60" height="60" style="margin: -12px;">
            </td>
        <tr>
            <td></td>
            <td>
            </td>
        </tr>
    </table>
    <table style="margin-top: 30px;" class="company" cellspacing="0" border="0">
        <tr>
            <td width=50% colspan=2>
                <span class="company_name">{{ strtoupper($setting['name']) }}</span>
                <br>
                <pre style="font-family: 'Helvetica';">{{ $setting['address'] }}</pre>
            </td>
        </tr>
        <tr>
            <td width=5% valign="top">Phone</td>
            <td width=95% valign="top">: {{ $setting['phone_number'] }}</td>
            <td></td>
        </tr>
        <tr>
            <td width=5% valign="top">Email</td>
            <td width=95% valign="top">: {{ $setting['email'] }}</td>
            <td><br/></td>
        </tr>
    </table>
    <table style="margin-top: 15px; margin-left:-10px; font-size:12px;">
        <tr>
            <td>
                <td width=5% valign="top">{{ $setting['npwp'] }}</td>
            </td>
        </tr>
    </table>
    <div style="float:right; margin-top:-120px;">
       <span style="float: right; margin-top:-40px; font-size:21px;">
           DELIVERY ORDER
       </span>
        <table style="font-size: 12px; margin-top: 15px" cellspacing="0" border="0" width="100%">
            <tr>
                <td style="text-align: right"><strong>Delivery Order# :</strong></td>
                <td style="text-align: right" width=30%><strong>{{ $order['code'] }}</strong></td>
            </tr>
            <tr>
                <td style="text-align: right">Sales Order# :</td>
                <td style="text-align: right" width=30%>{{ $order['sales_order']['code'] }}</td>
            </tr>
            <tr>
                <td style="text-align: right">Order Date :</td>
                <td style="text-align: right" width=30%>{{ date_format($orderDate,"d/m/Y") }}</td>
            </tr>
            <tr>
                <td style="text-align: right">Delivery Date :</td>
                <td style="text-align: right" width=30%>{{ date_format($deliveryDate,"d/m/Y") }} {{ $order['wrt']['name'] }}</td>
            </tr>
            <tr>
                <td style="text-align: right"></td>
                @if($order['sales_order']['order_type']['id'] != "393216")
                    <td style="text-align: right" width=30%><strong>{{ $order['sales_order']['term_payment_sls']['name'] }} - {{ $order['sales_order']['order_type']['name'] }}</strong></td>
                @endif
            </tr>
        </table>
    </div>
    <table border="1" style="border-collapse: collapse; margin-top:20px;" width="100%">
        <tr>
            <td style="border:1px solid blac; background-color: #5C5C5C; color: white; font-size:14px; padding: 10px;">
                <strong>
                    <img src='data:image/png;base64," . {{ $info }} . "' width="20" height="20" style="margin-top:5px;"> &nbsp;Customer memiliki <u>HAK</u> untuk :
                </strong>
            </td>
        </tr>
        <tr>
            <td style="padding:10px; font-size:14px; ">
                <p>
                    1. Memeriksa produk yang sudah dipesan bersama dengan kurir dalam waktu 15 menit. <br>
                    2. Mengembalikan produk yang tidak sesuai dengan pesanan <br> <br>
                    Dengan menandatangani, customer setuju untuk menerima barang dan tidak bisa mengembalikan pesanan
                    setelah kurir pergi
                </p>
            </td>
        </tr>
    </table>
    <table border="1" style="border-collapse: collapse; font-size: 12px; margin-top: 20px;" width="100%">
        <tr>
            <td style="padding:5px;">
                Customer Name : {{ $order['sales_order']['branch']['name'] }}
            </td>
            <td style="padding:5px;">
                Payment Term : {{ $order['sales_order']['term_invoice_sls']['name'] }}
            </td>
        </tr>
        <tr>
            <td style="padding:5px;">
                Contact Person : {{ $order['sales_order']['branch']['pic_name'] }}
            </td>
            <td style="padding:5px;">
                Phone Number : {{ $order['sales_order']['branch']['phone_number'] }}
            </td>
        </tr>
        <tr>
            <td style="padding:5px;">
                Total Koli: {{ $order['total_koli'] }} - 
                @if ($deliveryKoli)
                    @foreach ($deliveryKoli as $item)
                        {{ $item['koli']['name'] }} - ({{ $item['quantity'] }})
                    @endforeach
                @endif
            </td>
            @if($order['sales_order']['branch']['merchant']['tag_customer'] == 1) 
                <td style="padding:5px;">
                    Customer Tag: NC
                </td>
            @elseif($order['sales_order']['branch']['merchant']['tag_customer'] == 8) 
                <td style="padding:5px;">
                    Customer Tag: PC
                </td>
            @else
                <td style="padding:5px;">
                    Customer Tag: -
                </td>
            @endif 
        </tr>
        @if($order['sales_order']['archetype']['business_type']['abbreviation'] === 'LM' || $order['sales_order']['archetype']['business_type']['abbreviation'] === 'ED')
            <tr>
                <td style="padding:5px;" colspan="2">
                    Salesperson : {{ $order['sales_order']['salesperson']['name'] }} ( {{ $order['sales_order']['salesperson']['phone_number'] }} )
                </td>
            </tr>
        @endif
        <tr>
            <td style="padding:5px;" colspan="2">
                Shipping Address: {{ $order['shipping_address'] }}
            </td>
        </tr>
    </table>
    <br>
    <table border="1" style="border-collapse: collapse; font-size:12px;" width="100%">
        <thead>
        <tr>
            <th style="padding:5px;" width="5%">No</th>
            <th style="padding:5px; text-align: center;" width="72%">Product</th>
            <th style="padding:5px;" width="10%">Qty</th>
            <th style="padding:5px;" width="13%">Receive Note</th>
        </tr>
        </thead>
        <tbody>
        @foreach($order['delivery_order_items'] as $index => $items)
            <tr>
                <td style="padding:5px;">{{$index + 1}}</td>
                <td style="padding:5px;">{{ $items['item']['name']}}</td>
                <td style="padding:5px;">
                    {{ number_format($items['deliver_qty'], 2,',','.') }} {{ $items['item']['uom']['name'] }}
                </td>
                <td style="padding:5px;">
                    {{ $items['receipt_item_note'] }}
                </td>
            </tr>
        @endforeach
        </tbody>
        <tr>
            <td style="padding:10px;" colspan="4">
                <strong>Note</strong> : <br>
                {{ $order['note'] }}
            </td>
        </tr>
    </table>
    <br/>
    <table class="sign" cellspacing=0>
        <tr>
            <td style="padding-top:50px;margin-bottom: 10px;text-align: center">
                <img src='data:image/png;base64," . {{ $digstamp }} . "' width="240" height="40">
                <hr style="width: 90%">
            </td>
            <td style="padding-top:90px;margin-bottom: 10px;text-align: center">
                <hr style="width: 90%">
            </td>
            <td style="padding-top:90px;margin-bottom: 10px;text-align: center">
                <hr style="width: 90%">
            </td>
        </tr>
        <tr>
            <td width=33.3% style="text-align: center">Warehouse</td>
            <td width=33.3% style="text-align: center">Courier</td>
            <td width=33.3% style="text-align: center">Received By</td>
        </tr>
    </table>
</section>

</body>
</html>
