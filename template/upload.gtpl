<html>
<head>
    <title>Upload file</title>
</head>
<body>
  <section>
    <form enctype="multipart/form-data" action="/upload" method="post">
        <input type="file" name="uploadfile" />
        <input type="text" name="title"/>
        <select name="categories">
          <option value="animals">Animals</option>
          <option value="food">Food</option>
          <option value="people">People</option>
          <option value="cars">Cars</option>
          <option value="fraces">Fraces</option>
          <option value="memes">Memes</option>
        </select>
        <input type="hidden" name="token" value="{{.}}"/>
        <input type="submit" value="upload" />
    </form>
  </section>
</body>
</html>