'use strict';

var $sourceTextArea = $('.editor .source textarea');
var $goOutput = $('.editor .go-output');
var $goOutputTextArea = $goOutput.find('textarea');
var $console = $('.console');
var $consoleOutput = $console.find('.output');
var $spinner = $('.spinner');

var sourceCM = CodeMirror.fromTextArea($sourceTextArea.get(0), {
  mode: 'oden',
  title: 'Oden source code',
  theme: 'oden-playground',
  matchBrackets: true,
  indentUnit: 2,
  tabSize: 2,
  indentWithTabs: false
});

var outputCM = CodeMirror.fromTextArea($goOutputTextArea.get(0), {
  mode: 'go',
  readOnly: true,
  title: 'Go output'
});

function displayError(result) {
  $console.addClass('error');
  var md = result.error.split('\n').map(function (line) {
    return line
      .replace(/^.+?main.oden/i, 'main.oden')
      .replace(/\serror:(.?)/i, '<span class="heading heading-error"> error:</span>$1')
      .replace(/warning:(.?)/i, '<span class="heading heading-warning">warning:</span>$1')
      .replace (/^(For more information see.*)/i, '$1\n\n');
  }).join('\n');
  $consoleOutput.html(marked(md));
}

function displayEvent(event) {
  $consoleOutput.append($('<div>').text(event.Message));
}

function displayConsoleOutput(result) {
  if (result.consoleOutput.Errors) {
    $console.addClass('error');
    $consoleOutput.text(result.consoleOutput.Errors);
  } else if (result.consoleOutput.Events) {
    $console.removeClass('error');
    $consoleOutput.empty();
    result.consoleOutput.Events.forEach(displayEvent);
  } else {
    $console.removeClass('error');
    $consoleOutput.empty();
    $consoleOutput.append(
      $('<span>').addClass('no-output').text('No output.'));
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
  $(document).on('keydown', function (event) {
    if (event.shiftKey && event.keyCode === 13) {
      event.preventDefault();
      compileAndRun();
    }
  });
  $('button.run').click(function () {
    compileAndRun();
  });

  var $shareScreen = $('.share-screen');

  function hideShareScreen() {
    $shareScreen.css('display', 'none');
  }

  function showShareScreen(url) {
    var $url = $shareScreen.find('.url');
    $url.val(url);

    var $tweetButton = $shareScreen.find('.tweet-button');
    var intentUri = 'http://twitter.com/share?url=' +
      encodeURIComponent(url) +
      '&text=' +
      encodeURIComponent('Check this #odenlang program out:');
    $tweetButton.attr('href', intentUri);
    $tweetButton.click(function (e) {
      e.preventDefault();
      window.open(intentUri, '_blank', 'location=yes,height=300,width=600,scrollbars=no,status=no');
    });

    $shareScreen.css('display', 'flex');
    $url.select();
  }

  $shareScreen.find('.close').click(hideShareScreen);
  $shareScreen.click(function (e) {
    if ($shareScreen.is(e.target)) {
      hideShareScreen();
    }
  });

  function shareProgram() {
    var code = sourceCM.getValue();
    showSpinners();
    $.ajax({
      type: 'POST',
      url: '/p',
      contentType: 'application/json',
      data: JSON.stringify({ odenSource: code }),
      dataType: 'json'
    }).done(function (result) {
      var url = window.location.origin + result.path;
      window.history.replaceState({}, 'Shared Program', result.path);
      showShareScreen(url);
    }).fail(function () {
      console.error('Failed to share program.');
    }).always(function () {
      hideSpinners();
    });
  }

  $('button.share').click(shareProgram);

  $('#show-go-output').on('change', function () {
    var show = $(this).prop('checked');
    $goOutput.toggleClass('hidden', !show);
    outputCM.refresh();
  });
}

setupEditor();
