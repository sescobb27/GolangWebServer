<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>Index</title>
  <link rel="stylesheet" href="/css/bootstrap.min.css">
</head>
<body>
<section>
  <ul class="nav nav-tabs nav-justified">
    <li><a href="/login" class="btn" title="">Login</a></li>
    <li> <a href="/signup" class="btn" title="">Join Us</a></li>
    <li> <a href="/upload" class="btn" title="">Upload Files</a></li>
    <li> <a href="/search" class="btn" title="">See Files</a> </li>
  </ul>
</section>
<section class=" container form-group">
    <form action="/signup" method="post" class="form-inline">
      <label for="username">Username</label>
      <input type="text" name="username" class="form-control" value="{{.}}">

      <label for="password">Password</label>
      <input type="password" name="password" class="form-control">
      <input type="submit" value="Login" class="btn btn-default">
    </form>
</section>
</body>
</html>