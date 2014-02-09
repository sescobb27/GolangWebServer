<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>Show Uploads</title>
</head>
<body>
{{range $i,$file := .}}
<div>
    <h3>{{$file.Title}}</h3>
    <img src="{{$file.Path}}" alt="{{$file.Title}}">
</div>
{{end}}

</body>
</html>
