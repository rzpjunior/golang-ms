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
        .page-break {
            page-break-after: always;
        }
    </style>
</head>
<body>
    @foreach($temp as $index => $items)
    <table style="height:50px;" border="2" width="100%">
        <tr>
            <td style="text-align:center;font-size:16px;font-weight:bold;" colspan="2">
                {{ $items['branch'] }}
            </td>
        </tr>
        <tr>
            <td colspan="2" style="text-align:center;font-size:18px;">
                <h1>
                    {{ $items['picking']['code'] }}
                </h1> 
            </td>
        </tr>
        <tr>
            <td style="text-align:center; padding: 10px; font-size:16px;">
                {{ $items['deliveryDate'] }}
            </td>
            <td style="text-align:center; font-size:16px;">
                {{ $items['picking']['wrt']['name'] }}
            </td>
        </tr>
        <tr>
            <td colspan="2" style="text-align:center;font-size:16px;">
                {{ $dateNow }}
            </td>
        </tr>
    </table>
    @if(!$loop->last)
        <div class="page-break"></div>
    @endif
    @endforeach
<script type="text/javascript">
    window.onload = function () {
        window.print();
        window.close();
    };
</script>
</body>
</html>
