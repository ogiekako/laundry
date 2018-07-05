var xhr = new XMLHttpRequest();
xhr.onreadystatechange = function() {
  if (this.readyState == 4 && this.status == 200) {
    var data = JSON.parse(this.responseText);
    drawData(data)
  }
};
xhr.open("GET", "/data", true);
xhr.send();

var drawData = function(data) {
      // https://developers.google.com/chart/interactive/docs/gallery/annotatedtimeline
      google.charts.load('current', {'packages':['annotatedtimeline']});
      google.charts.setOnLoadCallback(drawChart);
      function drawChart() {
        var table = new google.visualization.DataTable();
        table.addColumn('date', 'DateTime');
        table.addColumn('number', 'Shake');
        table.addColumn('number', 'Tap');

        var len = data.shakes.length
        for (var i = 0; i < len; i++) {
          table.addRow([
              new Date(data.shakes[i].timestamp * 1000),
              data.shakes[i].value,
              undefined
          ])
        }
        var len = data.taps.length
        for (var i = 0; i < len; i++) {
          table.addRow([
              new Date(data.taps[i].timestamp * 1000),
              undefined,
              data.taps[i].value
          ])
        }

        var chart = new google.visualization.AnnotatedTimeLine(document.getElementById('chart_div'));
        chart.draw(table, {displayAnnotations: false});
      }
}
