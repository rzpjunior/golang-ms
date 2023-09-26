<!DOCTYPE html>
<html lang="en">
<head>
    <title>{{ $item['item_code'] }}</title>
    <meta charset="utf-8">
    <link rel="stylesheet" href="/assets/css/print.css">
    <style>
        @page { margin: 2px; }
        body {
            font-family: Arial, Helvetica, sans-serif;
            margin: 1px;
            margin-top:5px;
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
    </style>
</head>
<body>
    <table style="height:30px;" width="100%">
        <tr>
            <td colspan="2" style="text-align:center; font-size:14px; padding-top:15px;">
                {{ $dateNow }} 
            </td>
        </tr>
        <tr>
            <td style="text-align:center;" colspan="2">
                <img src='data:image/png;base64," . {{ $qrCode }} . "' width="120" height="120">
            </td>
        </tr>
        <tr>
            <td colspan="2" style="text-align:center; padding-top:10px; font-size:18px;">
                <strong>{{ $item['item_name'] }}</strong> 
            </td>
        </tr>
        <tr>
            @if ($item['pack_size'] === null)
                <td colspan="2" style="text-align:center; font-size:16px;">
                    {{ $item['total_order'] }} {{ $item['item_uom'] }}
                </td>
            @else
                <td colspan="2" style="text-align:center; font-size:16px;">
                    {{ $item['pack_size'] }}
                </td>
            @endif
            
        </tr>
        
    </table>

<script type="text/javascript">
    window.onload = function () {
        window.print();
        window.close();
    };

</script>
</body>
</html>
