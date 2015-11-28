'use strict';

var sourceTextArea = document.querySelector('.editor .source textarea');
var sourceCodeMirror = CodeMirror.fromTextArea(sourceTextArea, {
  mode: 'clojure',
  title: 'Oden source code',
  matchBrackets: true
});

var outputTextArea = document.querySelector('.editor .output textarea');
var outputCodeMirror = CodeMirror.fromTextArea(outputTextArea, {
  mode: 'go',
  readOnly: true,
  title: 'Go output'
});

document.addEventListener('keyup', function (event) {
  if (event.ctrlKey && event.keyCode === 82) {
    document.querySelector('form').submit();
  }
});
