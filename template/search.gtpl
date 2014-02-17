<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>Search Uploads</title>
  <link rel="stylesheet" href="/css/bootstrap.min.css">
</head>
<body>
  <nav class="navbar navbar-default" role="navigation">
    <ul class="nav nav-tabs nav-justified">
      <li><a href="/login" class="btn" title="">Login</a></li>
      <li> <a href="/signup" class="btn" title="">Join Us</a></li>
      <li> <a href="/upload" class="btn" title="">Upload Files</a></li>
      <li> <a href="/search" class="btn" title="">Search Files</a> </li>
    </ul>
  </nav>
  <div class="container form-group">
    <form action="/show" method="get" accept-charset="utf-8">
      {{range $i,$tag := .}}
      <input type="submit"action="/show" method="get" name="tag" value="{{$tag}}" placeholder="" class="btn btn-info">
      {{end}}
      <input type="text" name="tag" placeholder="Tag Name">
    </form>
  </div>
</body>
</html>
