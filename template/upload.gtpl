<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>Upload file</title>
  <link rel="stylesheet" href="/css/bootstrap.min.css">
</head>
<body>
  <section>
    <ul class="nav nav-tabs nav-justified">
      <li><a href="/login" class="btn" title="">Login</a></li>
      <li> <a href="/signup" class="btn" title="">Join Us</a></li>
      <li> <a href="/upload" class="btn" title="">Upload Files</a></li>
      <li> <a href="/show" class="btn" title="">See Files</a> </li>
    </ul>
  </section>
  <section class=" container form-group">
    <form enctype="multipart/form-data" action="/upload" method="post" class="form-inline">
        <input type="file" name="uploadfile" />
        <input type="text" name="title" class="form-control" placeholder="Title"/>
        <input type="text" name="tags" class="form-control" placeholder="Tags" />
        <input type="hidden" name="token" value="{{.}}"/>
        <input type="submit" value="upload" class="btn btn-default" />
    </form>
  </section>
</body>
</html>
