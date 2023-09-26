<!DOCTYPE html>
<html lang="en">
<head>
    <title>{{ $item['code'] }}</title>
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
    <div style="height:30px;">
        <img src='data:image/png;base64," . {{ $qrCode }} . "' width="155" height="155">
        <div style="float:right; font-size:19px;">
            <b>{{ $pkdata['code'] }}</b> <br>
            <div style="text-align: center">
                <h1>
                    {{ $pkdata['weight_scale'] }} KG 
                </h1>
               <h4>{{$dateNow}}</h4> 
            </div>
        </div>
        <div style="text-align: center;">
            <h3><span style="font-size:27px;">{{ $item['name'] }}</span></h3>
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
