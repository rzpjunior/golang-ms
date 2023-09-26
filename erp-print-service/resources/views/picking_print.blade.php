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
            <td style="text-align:center;" colspan="2">
                <h1>
                    {{ $items['picking']['sales_order']['code']}}-{{ $items['numKoli'] }}
                </h1> 
            </td>
        </tr>
        <tr>
            <td colspan="2" style="text-align:center; padding:10px; font-size:16px;">
                {{ $items['branch'] }}
            </td>
        </tr>
        <tr>
            <td rowspan="2" style="text-align:center; padding: 10px; font-size:16px;">
                <img src='data:image/png;base64," . {{ $items['qrCode'] }} . "' width="95" height="95">
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
