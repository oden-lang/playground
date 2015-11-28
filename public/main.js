'use strict';

var $sourceTextArea = $('.editor .source textarea');
var $output = $('.editor .output');
var $outputTextArea = $output.find('textarea');
var $console = $('.console');

var sourceCM = CodeMirror.fromTextArea($sourceTextArea.get(0), {
  mode: 'clojure',
  title: 'Oden source code',
  matchBrackets: true
});

var outputCM = CodeMirror.fromTextArea($outputTextArea.get(0), {
  mode: 'go',
  readOnly: true,
  title: 'Go output'
});


function displayError(result) {
  $output.addClass('error');
  outputCM.setValue(result.error);
  $console.empty();
}

function displayEvent(event) {
  $console.append($('<div>').text(event.Message));
}

function displayConsoleOutput(result) {
  if (result.consoleOutput.Errors) {
    $console.addClass('error').text(result.consoleOutput.Errors);
  } else {
    $console.removeClass('error').empty();
    result.consoleOutput.Events.forEach(displayEvent);
  }
}

function display(result) {
  $output.removeClass('error');
  outputCM.setValue(result.goOutput);
  displayConsoleOutput(result);
}

function compileAndRun(code) {
  console.log(code);
  $.ajax({
    type: 'POST',
    url: '/compile',
    data: JSON.stringify({ odenSource: code }),
    dataType: 'json'
  }).done(function (result) {
    if (result.error) {
      displayError(result);
    } else {
      display(result);
    }
  }).fail(function () {
    console.error('Failed to run code.');
  }).always(function () {
    console.log('Done');
  });
}

function setupEditor() {

  $(document).on('keyup', function (event) {
    if (event.ctrlKey && event.keyCode === 82) {
      var code = sourceCM.getValue();
      compileAndRun(code);
    }
  });
}

setupEditor();
