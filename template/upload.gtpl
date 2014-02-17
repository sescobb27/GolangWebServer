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
  <nav class="navbar navbar-default" role="navigation">
    <ul class="nav nav-tabs nav-justified">
      <li><a href="/login" class="btn" title="">Login</a></li>
      <li> <a href="/signup" class="btn" title="">Join Us</a></li>
      <li> <a href="/upload" class="btn" title="">Upload Files</a></li>
      <li> <a href="/search" class="btn" title="">See Files</a> </li>
    </ul>
  </nav>
  <div class=" container form-group">
        <form enctype="multipart/form-data" action="/upload" method="post">
          <div class="row">
            <div class="col-lg-2">
                <div class="input-group">
                  <input type="file" name="uploadfile" class="input-control"/>
                </div>
            </div>
            <input type="hidden" name="token" value="{{.}}"/>
          </div>
          <div class="row">
            <div class="col-lg-2">
              <div class="input-group">
                    <input type="text" name="title" class="input-control" placeholder="Title"/>
              </div>
            </div>
            <div class="col-lg-2">
              <div class="input-group">
                    <input type="text" name="tags" class="input-control" placeholder="Tags" />
              </div>
            </div>
            <input type="submit" value="upload" class="btn btn-success" />
          </div>
        </form>
  </div>
</body>
</html>
