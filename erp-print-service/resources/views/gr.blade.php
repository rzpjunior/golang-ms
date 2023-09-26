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
            width:70%;
            text-align:center;
            float:right;
            margin-top:50px;
        }
    </style>
</head>
<body>
<footer>
    <span>#{{ $order['code'] }}</span>
</footer>
<section class="sheet" style="padding: 10px;">
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
    <br/>
    <span style="font-size:12px;">{{ $setting['npwp'] }}</span>
    <br>
    <div style="float:right; margin-top:-120px;">
       <span style="float: right; margin-top:-40px; font-size:21px;">
           GOODS RECEIPT
       </span>
        <table style="font-size: 12px; margin-top: 15px" cellspacing="0" border="0" width="100%">
            <tr>
                <td style="text-align: right"><strong>Goods Receipt# :</strong></td>
                <td style="text-align: right" width=30%><strong>{{ $order['code'] }}</strong></td>
            </tr>
             <tr>
                <td style="text-align: right">Goods Receipt Date :</td>
                <td style="text-align: right" width=30%>{{ date_format($createdAtGR,"d/m/Y") }}</td>
            </tr>
            <tr>
                <td style="text-align: right">Actual Arrival Date :</td>
                <td style="text-align: right" width=30%>{{ date_format($ataGR,"d/m/Y") }}</td>
            </tr>
            <tr>
                <td style="text-align: right">Purchase Order# :</td>
                <td style="text-align: right" width=30%><strong>{{ $order['purchase_order']['code'] }}</strong></td>
            </tr>
            <tr>
                <td style="text-align: right">Purchase Order Date :</td>
                <td style="text-align: right" width=30%>{{ date_format($recognition_datePO,"d/m/Y") }}</td>
            </tr>
            <tr>
                <td style="text-align: right">Estimated Arrival Date :</td>
                <td style="text-align: right" width=30%>{{ date_format($etaPO,"d/m/Y") }}</td>
            </tr>
            <tr>
                <td style="text-align: right">Status :</td>
                <td style="text-align: right" width=30%>
                    @if ($order['status'] == 1)
                        Active
                    @elseif ($order['status'] == 2)
                        Finished
                    @else
                        Cancelled
                    @endif
                </td>
            </tr>
        </table>
    </div>
    <table border="1" style="border-collapse: collapse; font-size: 12px; margin-top:20px;" width="100%">
        <tr>
            <td width="50%" style="padding:8px;">Supplier : {{ $order['purchase_order']['supplier']['name'] }}</td>
            <td width="50%" style="padding:8px;">Payment Term : {{ $order['purchase_order']['supplier']['purchase_term']['name'] }}</td>
        </tr>
        <tr>
            <td width="50%" style="padding:8px;">PIC Supplier : {{ $order['purchase_order']['supplier']['pic_name'] }}</td>
            <td width="50%" style="padding:8px;">Payment Term : {{ $order['purchase_order']['eta_time'] }}</td>
        </tr>
        <tr>
            <td colspan="2" style="padding:8px;">
                <table border="0" style="width:100%">
                    <tr>
                        <td width="20%">Supplier Address</td>
                        <td width="1%">:</td>
                        <td width="80%">
                            {{ $order['purchase_order']['supplier']['address'] }}
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td colspan="2" style="padding:8px;">
                <table border="0" style="width:100%">
                    <tr>
                        <td width="20%">Ship to</td>
                        <td width="1%">:</td>
                        <td width="80%">
                            {{ $order['site']['street_address'] }}
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
    <br>
    <table border="1" style="border-collapse: collapse; font-size:12px;" width="100%">
        <thead>
        <tr>
            <th style="padding-top:8px;padding-bottom:8px;" width="5%">No</th>
            <th style="padding:8px; text-align: left;" width="40%">Product </th>
            <th style="padding-top:8px;padding-bottom:8px;">UOM </th>
            <th style="padding-top:8px;padding-bottom:8px;">Order Qty </th>
            <th style="padding-top:8px;padding-bottom:8px;">Delivery Qty </th>
            <th style="padding-top:8px;padding-bottom:8px;">Reject Qty </th>
            <th style="padding-top:8px;padding-bottom:8px;">Received Qty </th>
        </tr>
        </thead>
        <tbody>
        @foreach($order['goods_receipt_items'] as $index => $items)
            <tr>
                <td style="padding:5px;">{{$index + 1}}</td>
                <td style="padding:5px;">{{ $items['item']['code']  }} - {{ $items['item']['name']}}<br>
                    @if($items['note'])
                        <i class="second-color">Note: {{ $items['note'] }}</i>
                    @endif
                </td>
                <td  style="padding-top:8px;padding-bottom:8px;text-align: center">
                    {{ $items['item']['uom']['name'] }}
                </td>
                <td style="padding-top:8px;padding-bottom:8px;text-align: center">
                    {{ $items['purchase_order_item']['order_qty'] }}
                </td>
                <td style="padding-top:8px;padding-bottom:8px;text-align: center">
                    {{ $items['delivery_qty'] }}
                </td>
                <td style="padding-top:8px;padding-bottom:8px;text-align: center">
                    {{ $items['reject_qty'] }}
                </td>
                <td style="padding-top:8px;padding-bottom:8px;text-align: center">
                    {{ $items['receive_qty'] }}
                </td>
            </tr>
        @endforeach
        </tbody>
    </table>
    <table>
        <tr>
            <td>
                @if($order['note'])
                    <strong>Note :</strong> <br> <i>{{ $order['note'] }}</i>
                @endif
            </td>
        </tr>
    </table>
    <br/>
    <table class="sign" cellspacing=0>
        <tr>
            <td width=23.3% style="padding-top:90px;margin-bottom: 10px;text-align: center"><hr style="width: 90%"></td>
            <td width=23.3% style="padding-top:90px;margin-bottom: 10px;text-align: center"><hr style="width: 90%"></td>
            <td width=23.3% style="padding-top:90px;margin-bottom: 10px;text-align: center"><hr style="width: 90%"></td>
        </tr>
        <tr>
            <td width=33.3% style="text-align: center">Warehouse</td>
            <td width=33.3% style="text-align: center">Supplier</td>
            <td width=33.3% style="text-align: center">Security</td>
        </tr>
    </table>
</section>

</body>
</html>
