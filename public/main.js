'use strict';

var $sourceTextArea = $('.editor .source textarea');
var $goOutput = $('.editor .go-output');
var $goOutputTextArea = $goOutput.find('textarea');
var $console = $('.console');
var $consoleOutput = $console.find('.output');
var $spinner = $('.spinner');

var sourceCM = CodeMirror.fromTextArea($sourceTextArea.get(0), {
  mode: 'clojure',
  title: 'Oden source code',
  matchBrackets: true
});

var outputCM = CodeMirror.fromTextArea($goOutputTextArea.get(0), {
  mode: 'go',
  readOnly: true,
  title: 'Go output'
});

function displayError(result) {
  $goOutput.addClass('error');
  outputCM.setValue(result.error);
  $consoleOutput.empty();
}

function displayEvent(event) {
  $consoleOutput.append($('<div>').text(event.Message));
}

function displayConsoleOutput(result) {
  if (result.consoleOutput.Errors) {
    $console.addClass('error');
    $consoleOutput.text(result.consoleOutput.Errors);
  } else {
    $console.removeClass('error');
    $consoleOutput.empty();
    result.consoleOutput.Events.forEach(displayEvent);
  }
}

function showSpinners() {
  $spinner.css('display', 'flex');
  outputCM.setValue('');
  $consoleOutput.empty();
}

function hideSpinners() {
  $spinner.css('display', 'none');
}

function display(result) {
  $goOutput.removeClass('error');
  outputCM.setValue(result.goOutput);
  displayConsoleOutput(result);
}

function compileAndRun() {
  var code = sourceCM.getValue();
  showSpinners();

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
    hideSpinners();
  });
}

function setupEditor() {
  $(document).on('keyup', function (event) {
    if (event.ctrlKey && event.keyCode === 82) {
      compileAndRun();
    }
  });
  $('button.run').click(function () {
    compileAndRun();
  });
}

setupEditor();
