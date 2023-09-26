<!DOCTYPE html>
<html lang="en">
<head>
    <title>{{ $waste['code'] }}</title>
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

        .sign{
            width:70%;
            text-align:center;
            float:right;
            margin-top:20px;
        }
    </style>
</head>
<body>
<footer>
    <span>#{{ $waste['code'] }}</span>
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
           WASTE DISPOSAL
       </span>
        <table style="font-size: 12px; margin-top: 15px" cellspacing="0" border="0" width="100%">
            <tr>
                <td style="text-align: right"><strong>Waste Disposal# :</strong></td>
                <td style="text-align: right" width=30%><strong>{{ $waste['code'] }}</strong></td>
            </tr>
            <tr>
                <td style="text-align: right">Waste Disposal Date :</td>
                <td style="text-align: right" width=30%>{{ date_format($recognitionDate,"d/m/Y") }}</td>
            </tr>
        </table>
    </div>
    <table border="1" style="border-collapse: collapse; font-size: 12px; margin-top: 20px;" width="100%">
        <tr>
            <td style="padding:5px; border-right: solid 1px #FFF;">
                Area
            </td>
            <td style="padding:5px; border-left: solid 1px #FFF;">
                : {{ $waste['site']['area']['code'] }} - {{ $waste['site']['area']['name'] }}
            </td>
        </tr>
        <tr>
            <td style="padding:5px; border-right: solid 1px #FFF;">
                Warehouse
            </td>
            <td style="padding:5px; border-left: solid 1px #FFF;">
                : {{ $waste['site']['code'] }} - {{ $waste['site']['name'] }} - {{ $waste['site']['street_address'] }}, {{ $waste['site']['sub_district']['name'] }},
                {{ $waste['site']['sub_district']['district']['name'] }}, {{ $waste['city'] }}, {{ $waste['province'] }}, {{ $waste['site']['sub_district']['postal_code'] }}
            </td>
        </tr>
    </table>

    <br>
    <table border="1" style="border-collapse: collapse; font-size:12px;" width="100%">
        <thead>
        <tr>
            <th style="padding:5px;" width="5%">No</th>
            <th style="padding:5px; text-align: center;" width="65%">Product</th>
            <th style="padding:5px;" width="17%">UOM</th>
            <th style="padding:5px;" width="13%">Dispose Qty</th>
            <th style="padding:5px;" width="20%">Note</th>
        </tr>
        </thead>
        <tbody>
        @foreach($waste['waste_disposal_item'] as $index => $items)
            <tr>
                <td style="padding:5px;">{{$index + 1}}</td>
                <td style="padding:5px;">{{ $items['item']['name']}}</td>
                <td style="padding:5px;">
                    {{ $items['item']['uom']['name'] }}
                </td>
                <td style="padding:5px;">
                    {{ number_format($items['dispose_qty'], 2,',','.') }}
                </td>
                <td style="padding:5px; text-align:right;">{{ $items['note'] }}</td>
            </tr>
        @endforeach
            <tr>
                <td style="padding:5px;" colspan="5">
                    <strong>Note :</strong> <br> <i>{{ $waste['note'] }}</i>
                </td>
            </tr>
        </tbody>
    </table>
    <table class="sign" cellspacing=0>
        <tr>
            <td width=46.7% style="padding-top:90px;margin-bottom: 10px;text-align: center"></td>
            <td width=23.3% style="padding-top:90px;margin-bottom: 10px;text-align: center"><hr style="width: 90%"></td>
        </tr>
        <tr>
            <td width=46.7% style="text-align: center"></td>
            <td width=33.3% style="text-align: center">Manager on Duty</td>
        </tr>
    </table>
</section>

<script type="text/javascript">
    window.onload = function () {
        window.print();
        window.close();
    };

</script>
</body>
</html>
