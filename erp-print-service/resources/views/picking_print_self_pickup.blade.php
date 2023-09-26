<!DOCTYPE html>
<html lang="en">
<head>
    {{-- <title>{{ $picking['sales_order']['code'] }}</title> --}}
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
        .verticalLeft {
            text-transform: uppercase;
            text-align: center;
            white-space: nowrap;
            writing-mode: vertical-lr;
            -ms-writing-mode: tb-rl;
            transform: rotate(-270deg);
            width: 1px;
            padding: 1px;
            margin: 1px;
            left: 30px;
            vertical-align: middle;
            position: relative;
        }
        .verticalRight {
            text-transform: uppercase;
            text-align: center;
            white-space: nowrap;
            writing-mode: vertical-lr;
            -ms-writing-mode: tb-rl;
            transform: rotate(270deg);
            width: 1px;
            padding: 1px;
            margin: 1px;
            left: 30px;
            vertical-align: middle;
            position: relative;
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
    @foreach($temp as $index => $items)
    <table style="height:30px;" border="2" width="100%">
        <tr>
            <td style="text-align:center; padding: 3px; font: bold; font-size: 20px; white-space: break-spaces; padding: 10px;" colspan="2">
                {{ $items['picking']['sales_order']['code']}}-{{ $items['numKoli'] }}
            </td>
        </tr>
        <tr>
            <td colspan="2" style="text-align:center; padding: 10px; font-size:16px;">
                {{ $items['branch'] }}
            </td>
        </tr>
        <tr>
            <td rowspan="3" style="text-align:center; font-size:16px;">
                <img src='data:image/png;base64," . {{ $items['qrCode'] }} . "' width="120" height="120">
            </td>
            <td style="text-align:center; font-size:16px;">
                {{ $items['picking']['helper']['code'] }} <br>
                {{ $items['picking']['sales_order']['wrt']['name'] }}
            </td>
        </tr>
        <tr>
            <td style="text-align:center; font-size:16px;">
               {{ $items['numKoli'] }} / {{ $items['picking']['total_koli'] }} Koli
            </td>
        </tr>
        <tr>
            <td style="background-color: black; color: white; font: bold; font-size: 20px; padding: 10px;">
                <div style="text-align: center;">AMBIL SENDIRI</div>
            </td>
        </tr>
    </table>
    @endforeach
        

<script type="text/javascript">
    window.onload = function () {
        window.print();
        window.close();
    };

</script>
</body>
</html>
