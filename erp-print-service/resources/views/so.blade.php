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
<section class="sheet" style="padding: 10px;">
    @if($order['order_type']['id'] == "393216")
    <div class="self-pickup">
        {{ $order['term_payment_sls']['name'] }} - {{ $order['order_type']['name'] }}
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
    <br/>
    <span style="font-size:12px;">{{ $setting['npwp'] }}</span>
    <br>
    <div style="float:right; margin-top:-120px;">
       <span style="float: right; margin-top:-40px; font-size:21px;">
           SALES ORDER
       </span>
        <table style="font-size: 12px; margin-top: 15px" cellspacing="0" border="0" width="100%">
            <tr>
                <td style="text-align: right"><strong>Sales Order# :</strong></td>
                <td style="text-align: right" width=30%><strong>{{ $order['code'] }}</strong></td>
            </tr>
            <tr>
                <td style="text-align: right">Order Date :</td>
                <td style="text-align: right" width=30%>{{ date_format($recognitionDate,"d/m/Y") }}</td>
            </tr>
            <tr>
                <td style="text-align: right">Delivery Date :</td>
                <td style="text-align: right"
                    width=30%>{{ date_format($deliveryDate,"d/m/Y") }} {{ $order['wrt']['name'] }}</td>
            </tr>
            <tr>
                <td style="text-align: right"></td>
                @if($order['order_type']['id'] == "393216")
                    <td style="font-weight:bold; text-align:right;">
                        {{ $order['archetype']['business_type']['abbreviation'] }}
                    </td>
                @else
                    <td style="text-align: right"><strong>{{ $order['archetype']['business_type']['abbreviation'] }} - {{ $order['order_type']['name'] }}</strong></td>
                @endif
            </tr>
        </table>
    </div>
    <table border="1" style="border-collapse: collapse; font-size: 12px; margin-top:20px;" width="100%">
        <tr>
            <td style="padding:5px;">Customer Name : {{ $order['branch']['name'] }}</td>
            <td style="padding:5px;">
                Contact Person : {{ $order['branch']['pic_name'] }}<br>
                Phone Number : {{ $order['branch']['phone_number'] }}
            </td>
        </tr>
        <tr>
            <td style="padding:5px;" colspan="2">
                Billing Address: {{ $order['billing_address'] }}
            </td>
        </tr>
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
            <th style="padding:5px; text-align: center;" width="65%">Product</th>
            <th style="padding:5px;" width="17%">Qty</th>
            <th style="padding:5px;" width="13%">Unit Price</th>
            <th style="padding:5px;" width="20%">Discount Amount</th>
            <th style="padding:5px;" width="20%">Subtotal</th>
        </tr>
        </thead>
        <tbody>
        @foreach($order['sales_order_items'] as $index => $items)
            <tr>
                <td style="padding:5px;">{{$index + 1}}</td>
                <td style="padding:5px;">{{ $items['item']['name']}}<br>
                    @if($items['note'])
                        <i class="second-color">Note: {{ $items['note'] }}</i>
                    @endif
                </td>
                <td style="padding:5px;">
                    {{ number_format($items['order_qty'], 2,',','.') }} {{ $items['item']['uom']['name'] }}
                </td>
                <td style="padding:5px; text-align:right;">{{ number_format($items['unit_price'], 2,",",".") }}</td>
                @if($items['discount_qty'] > 0)
                    <td style="padding:5px; text-align:right;">
                        {{ number_format(($items['sku_disc_amount']), 2,",",".") }}
                    </td>
                @else
                    <td style="padding:5px; text-align:right;"></td>
                @endif
                <td style="padding:5px; text-align:right;">{{ number_format($items['subtotal'], 2,",",".") }}</td>
            </tr>
        @endforeach
        </tbody>
    </table>
    <div class="row">
        <div class="column">
            <table>
                <tr>
                    <td style="word-wrap: break-word;">
                        @if($order['note'])
                            <strong>Note :</strong> <br> <i>{{ $order['note'] }}</i>
                        @endif
                    </td>
                </tr>
            </table>
        </div>
        <div class="column">
            <table style="border-collapse: collapse; font-size:12px;" width="100% float:right;" border="0">
                <tr>
                    <td>Subtotal (Rp)</td>
                    <td style="text-align: right;">
                        <span style="margin-right: 30px;">
                            {{ number_format($order['total_price'], 2,",",".") }}
                        </span>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">
                        <hr style="width: 305px; margin-left: 1px">
                    </td>
                </tr>
                @if($order['voucher'] && $order['voucher']['type'] == 1)
                    <tr>
                        <td>Total Discount (Rp)</td>
                        <td style="text-align: right;">
                            <span style="margin-right: 30px;">
                                {{ number_format($order['vou_disc_amount'], 2,",",".") }}
                            </span>
                        </td>
                    </tr>
                    <tr>
                        <td colspan="2">
                            <hr style="width: 305px; margin-left: 1px">
                        </td>
                    </tr>
                @endif
                @if($order['voucher'] && $order['voucher']['type'] == 3)
                    <tr>
                        <td>Delivery Discount (Rp)</td>
                        <td style="text-align: right;">
                            <span style="margin-right: 30px;">
                                {{ number_format($order['vou_disc_amount'], 2,",",".") }}
                            </span>
                        </td>
                    </tr>
                    <tr>
                        <td colspan="2">
                            <hr style="width: 305px; margin-left: 1px">
                        </td>
                    </tr>
                @endif
                @if($order['voucher'] && $order['voucher']['type'] == 2)
                    <tr>
                        <td>Grand Total Discount (Rp)</td>
                        <td style="text-align: right;">
                            <span style="margin-right: 30px;">
                                - {{ number_format($order['vou_disc_amount'], 2,",",".") }}
                            </span>
                        </td>
                    </tr>
                    <tr>
                        <td colspan="2">
                            <hr style="width: 305px; margin-left: 1px">
                        </td>
                    </tr>
                @endif
            <!-- POINT REDEEM -->
                @if($order['point_redeem_amount'] && $order['point_redeem_amount'] != 0.00)
                    <tr>
                        <td>EdenPoint</td>
                        <td style="text-align: right;">
                            <span style="margin-right: 30px;">
                                - {{ number_format($order['point_redeem_amount'], 2,",",".") }}
                            </span>
                        </td>
                    </tr>
                    <tr>
                        <td colspan="2">
                            <hr style="width: 305px; margin-left: 1px">
                        </td>
                    </tr>
                @endif
                <tr>
                    <td>Delivery Fee</td>
                    <td style="text-align: right;">
                        <span style="margin-right: 30px;">
                            {{ number_format($order['delivery_fee'], 2,",",".") }}
                        </span>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">
                        <hr style="width: 305px; margin-left: 1px">
                    </td>
                </tr>
                <tr>
                    <td><strong>Grand Total</strong></td>
                    <td style="text-align: right;">
                        <strong style="margin-right: 30px;">
                            {{ number_format($order['total_charge'], 2,",",".") }}
                        </strong>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">
                        <hr style="width: 305px; margin-left: 1px">
                    </td>
                </tr>
            </table>
        </div>
    </div>
</section>

<script type="text/javascript">
    window.onload = function () {
        window.print();
        window.close();
    };

</script>
</body>
</html>
