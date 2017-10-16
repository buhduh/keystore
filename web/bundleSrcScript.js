var browserify = require('browserify');

var b = browserify(),
  path = require('path'),
  fs = require('fs'),
  bundlePath = path.join(__dirname, '../', 'build', '.js'),
  srcPath = './src/js/';

var entries = [
  'passwords.js',
  'edit_password.js'
];

entries.forEach(function(entry) {
  browserify().add(srcPath + entry)
    .bundle()
    .on('error', function (err) { console.error(err); })
    .pipe(fs.createWriteStream(path.join(__dirname, '../assets/js', entry)));
});
