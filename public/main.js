//document.querySelector('textarea').select();

var sourceTextArea = document.querySelector('.editor .source textarea');
var sourceCodeMirror = CodeMirror.fromTextArea(sourceTextArea, {
  mode: 'clojure',
  title: 'Oden source code'
});

var outputTextArea = document.querySelector('.editor .output textarea');
var outputCodeMirror = CodeMirror.fromTextArea(outputTextArea, {
  mode: 'go',
  readOnly: true,
  title: 'Go output'
});
